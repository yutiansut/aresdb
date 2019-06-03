//  Copyright (c) 2017-2018 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memstore

import (
	m3Shard "github.com/m3db/m3/src/cluster/shard"
	xerrors "github.com/m3db/m3/src/x/errors"
	xsync "github.com/m3db/m3/src/x/sync"
	"github.com/uber/aresdb/cluster/topology"
	"github.com/uber/aresdb/datanode/bootstrap"
	"github.com/uber/aresdb/diskstore"
	"github.com/uber/aresdb/memstore/common"
	"github.com/uber/aresdb/metastore"
	"github.com/uber/aresdb/redolog"
	"github.com/uber/aresdb/utils"
	"math"
	"math/rand"
	"runtime"
	"sync"
)

// TableShard stores the data for one table shard in memory.
type TableShard struct {
	bootstrapLock sync.RWMutex

	// Wait group used to prevent the stores from being prematurely deleted.
	Users sync.WaitGroup `json:"-"`

	ShardID int `json:"-"`

	// For convenience, reference to the table schema struct.
	Schema *common.TableSchema `json:"schema"`

	// For convenience.
	metaStore            metastore.MetaStore
	diskStore            diskstore.DiskStore
	redoLogManagerMaster *redolog.RedoLogManagerMaster

	// Live store. Its locks also cover the primary key.
	LiveStore *LiveStore `json:"liveStore"`

	// Archive store.
	ArchiveStore *ArchiveStore `json:"archiveStore"`

	// The special column deletion lock,
	// see https://docs.google.com/spreadsheets/d/1QI3s1_4wgP3Cy-IGoKFCx9BcN23FzIfZGRSNC8I-1Sk/edit#gid=0
	columnDeletion sync.Mutex

	// For convenience.
	HostMemoryManager common.HostMemoryManager `json:"-"`

	bootstrapState bootstrap.BootstrapState
}

// NewTableShard creates and initiates a table shard based on the schema.
func NewTableShard(schema *common.TableSchema, metaStore metastore.MetaStore,
	diskStore diskstore.DiskStore, hostMemoryManager common.HostMemoryManager, shard int, redoLogManagerMaster *redolog.RedoLogManagerMaster) *TableShard {
	tableShard := &TableShard{
		ShardID:              shard,
		Schema:               schema,
		diskStore:            diskStore,
		metaStore:            metaStore,
		HostMemoryManager:    hostMemoryManager,
		redoLogManagerMaster: redoLogManagerMaster,
	}
	archiveStore := NewArchiveStore(tableShard)
	tableShard.ArchiveStore = archiveStore
	tableShard.LiveStore = NewLiveStore(schema.Schema.Config.BatchSize, tableShard)
	return tableShard
}

// Destruct destructs the table shard.
// Caller must detach the shard from memstore first.
func (shard *TableShard) Destruct() {
	// TODO: if this blocks on archiving for too long, figure out a way to cancel it.
	shard.Users.Wait()

	shard.redoLogManagerMaster.Close(shard.Schema.Schema.Name, shard.ShardID)

	shard.LiveStore.Destruct()

	if shard.Schema.Schema.IsFactTable {
		shard.ArchiveStore.Destruct()
	}
}

// DeleteColumn deletes the data for the specified column.
func (shard *TableShard) DeleteColumn(columnID int) error {
	shard.columnDeletion.Lock()
	defer shard.columnDeletion.Unlock()

	// Delete from live store
	shard.LiveStore.WriterLock.Lock()
	batchIDs, _ := shard.LiveStore.GetBatchIDs()
	for _, batchID := range batchIDs {
		batch := shard.LiveStore.GetBatchForWrite(batchID)
		if batch == nil {
			continue
		}
		if columnID < len(batch.Columns) {
			vp := batch.Columns[columnID]
			if vp != nil {
				bytes := vp.GetBytes()
				batch.Columns[columnID] = nil
				vp.SafeDestruct()
				shard.HostMemoryManager.ReportUnmanagedSpaceUsageChange(int64(-bytes))
			}
		}
		batch.Unlock()
	}
	shard.LiveStore.WriterLock.Unlock()

	if !shard.Schema.Schema.IsFactTable {
		return nil
	}

	// Delete from disk store
	// Schema cannot be changed while this function is called.
	// Only delete unsorted columns from disk.
	if utils.IndexOfInt(shard.Schema.Schema.ArchivingSortColumns, columnID) < 0 {
		err := shard.diskStore.DeleteColumn(shard.Schema.Schema.Name, columnID, shard.ShardID)
		if err != nil {
			return err
		}
	}

	// Delete from archive store
	currentVersion := shard.ArchiveStore.GetCurrentVersion()
	defer currentVersion.Users.Done()

	var batches []*ArchiveBatch
	currentVersion.RLock()
	for _, batch := range currentVersion.Batches {
		batches = append(batches, batch)
	}
	currentVersion.RUnlock()

	for _, batch := range batches {
		batch.BlockingDelete(columnID)
	}
	return nil
}

// PreloadColumn loads the column into memory and wait for completion of loading
// within (startDay, endDay]. Note endDay is inclusive but startDay is exclusive.
func (shard *TableShard) PreloadColumn(columnID int, startDay int, endDay int) {
	archiveStoreVersion := shard.ArchiveStore.GetCurrentVersion()
	for batchID := endDay; batchID > startDay; batchID-- {
		batch := archiveStoreVersion.RequestBatch(int32(batchID))
		// Only do loading if this batch does not have any data yet.
		if batch.Size > 0 {
			vp := batch.RequestVectorParty(columnID)
			vp.WaitForDiskLoad()
			vp.Release()
		}
	}
	archiveStoreVersion.Users.Done()
}

// IsBootstrapped returns whether this table shard is already bootstrapped.
func (shard *TableShard) IsBootstrapped() bool {
	return shard.BootstrapState() == bootstrap.Bootstrapped
}

// BootstrapState returns this table shards' bootstrap state.
func (shard *TableShard) BootstrapState() bootstrap.BootstrapState {
	shard.bootstrapLock.RLock()
	bs := shard.bootstrapState
	shard.bootstrapLock.RUnlock()
	return bs
}

// Bootstrap returns this table shards' bootstrap state.
func (shard *TableShard) Bootstrap(
	origin topology.Host, topo topology.Topology, topoState *topology.StateSnapshot) error {
	shard.bootstrapLock.Lock()
	if shard.bootstrapState == bootstrap.Bootstrapping {
		shard.bootstrapLock.Unlock()
		return bootstrap.ErrTableShardIsBootstrapping
	}
	shard.bootstrapState = bootstrap.Bootstrapping
	shard.bootstrapLock.Unlock()

	success := false
	defer func() {
		shard.bootstrapLock.Lock()
		if success {
			shard.bootstrapState = bootstrap.Bootstrapped
		} else {
			shard.bootstrapState = bootstrap.BootstrapNotStarted
		}
		shard.bootstrapLock.Unlock()
	}()

	shard.bootstrapLock.Lock()

	//peerNode := shard.findBootstrapSource(origin, topo, topoState)
	// TODO: Get the table shard metadata for the peerNode via GRPC
	// 1. distributed log offset
	// 2. archiving cutoff
	// 3. archive version
	// 4. a list of archive log files or snapshot files

	workers := xsync.NewWorkerPool(int(math.Ceil(float64(runtime.NumCPU()) / 2)))
	workers.Init()

	var (
		multiErr = xerrors.NewMultiError()
		//mutex    sync.Mutex
		//wg       sync.WaitGroup
	)
	/*
		for _, file := range files {
			wg.Add(1)
			workers.Go(func() {
				// copy archive log or snapshot
				mutex.Lock()
				multiErr = multiErr.Add(err)
				mutex.Unlock()

				wg.Done()
			})
		}
		wg.Wait()
	*/
	err := multiErr.FinalError()
	success = err == nil

	return err
}

func (shard *TableShard) findBootstrapSource(
	origin topology.Host, topo topology.Topology, topoState *topology.StateSnapshot) topology.Host {

	peers := make([]topology.Host, 0, topo.Get().HostsLen())
	hostShardStates, ok := topoState.ShardStates[topology.ShardID(shard.ShardID)]
	if !ok {
		// This shard was not part of the topology when the bootstrapping
		// process began.
		return nil
	}

	for _, hostShardState := range hostShardStates {
		if hostShardState.Host.ID() == origin.ID() {
			// Don't take self into account
			continue
		}
		shardState := hostShardState.ShardState
		switch shardState {
		// Don't want to peer bootstrap from a node that has not yet completely
		// taken ownership of the shard.
		case m3Shard.Initializing:
			// Success cases - We can bootstrap from this host, which is enough to
			// mark this shard as bootstrappable.
		case m3Shard.Leaving:
			fallthrough
		case m3Shard.Available:
			peers = append(peers, hostShardState.Host)
		case m3Shard.Unknown:
			fallthrough
		default:
		}
	}

	if len(peers) == 0 {
		utils.GetLogger().
			With("table", shard.Schema.Schema.Name).
			With("shardID", shard.ShardID).
			With("origin", origin.ID()).
			With("source", "").
			Info("no available bootstrap sorce")
	} else {

		utils.GetLogger().
			With("table", shard.Schema.Schema.Name).
			With("shardID", shard.ShardID).
			Info("bootstrap peers")
	}
	idx := rand.Intn(len(peers))
	return peers[idx]
}
