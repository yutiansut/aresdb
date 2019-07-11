// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bootstrap "github.com/uber/aresdb/datanode/bootstrap"
import client "github.com/uber/aresdb/datanode/client"
import common "github.com/uber/aresdb/memstore/common"
import memstore "github.com/uber/aresdb/memstore"
import mock "github.com/stretchr/testify/mock"
import topology "github.com/uber/aresdb/cluster/topology"

// MemStore is an autogenerated mock type for the MemStore type
type MemStore struct {
	mock.Mock
}

// AddTableShard provides a mock function with given fields: table, shardID, needPeerCopy
func (_m *MemStore) AddTableShard(table string, shardID int, needPeerCopy bool) {
	_m.Called(table, shardID, needPeerCopy)
}

// Archive provides a mock function with given fields: table, shardID, cutoff, reporter
func (_m *MemStore) Archive(table string, shardID int, cutoff uint32, reporter memstore.ArchiveJobDetailReporter) error {
	ret := _m.Called(table, shardID, cutoff, reporter)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int, uint32, memstore.ArchiveJobDetailReporter) error); ok {
		r0 = rf(table, shardID, cutoff, reporter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Backfill provides a mock function with given fields: table, shardID, reporter
func (_m *MemStore) Backfill(table string, shardID int, reporter memstore.BackfillJobDetailReporter) error {
	ret := _m.Called(table, shardID, reporter)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int, memstore.BackfillJobDetailReporter) error); ok {
		r0 = rf(table, shardID, reporter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Bootstrap provides a mock function with given fields: peerSource, origin, topo, topoState, options
func (_m *MemStore) Bootstrap(peerSource client.PeerSource, origin string, topo topology.Topology, topoState *topology.StateSnapshot, options bootstrap.Options) error {
	ret := _m.Called(peerSource, origin, topo, topoState, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(client.PeerSource, string, topology.Topology, *topology.StateSnapshot, bootstrap.Options) error); ok {
		r0 = rf(peerSource, origin, topo, topoState, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchSchema provides a mock function with given fields:
func (_m *MemStore) FetchSchema() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHostMemoryManager provides a mock function with given fields:
func (_m *MemStore) GetHostMemoryManager() common.HostMemoryManager {
	ret := _m.Called()

	var r0 common.HostMemoryManager
	if rf, ok := ret.Get(0).(func() common.HostMemoryManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.HostMemoryManager)
		}
	}

	return r0
}

// GetMemoryUsageDetails provides a mock function with given fields:
func (_m *MemStore) GetMemoryUsageDetails() (map[string]memstore.TableShardMemoryUsage, error) {
	ret := _m.Called()

	var r0 map[string]memstore.TableShardMemoryUsage
	if rf, ok := ret.Get(0).(func() map[string]memstore.TableShardMemoryUsage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]memstore.TableShardMemoryUsage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetScheduler provides a mock function with given fields:
func (_m *MemStore) GetScheduler() memstore.Scheduler {
	ret := _m.Called()

	var r0 memstore.Scheduler
	if rf, ok := ret.Get(0).(func() memstore.Scheduler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(memstore.Scheduler)
		}
	}

	return r0
}

// GetSchema provides a mock function with given fields: table
func (_m *MemStore) GetSchema(table string) (*common.TableSchema, error) {
	ret := _m.Called(table)

	var r0 *common.TableSchema
	if rf, ok := ret.Get(0).(func(string) *common.TableSchema); ok {
		r0 = rf(table)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.TableSchema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(table)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSchemas provides a mock function with given fields:
func (_m *MemStore) GetSchemas() map[string]*common.TableSchema {
	ret := _m.Called()

	var r0 map[string]*common.TableSchema
	if rf, ok := ret.Get(0).(func() map[string]*common.TableSchema); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*common.TableSchema)
		}
	}

	return r0
}

// GetTableShard provides a mock function with given fields: table, shardID
func (_m *MemStore) GetTableShard(table string, shardID int) (*memstore.TableShard, error) {
	ret := _m.Called(table, shardID)

	var r0 *memstore.TableShard
	if rf, ok := ret.Get(0).(func(string, int) *memstore.TableShard); ok {
		r0 = rf(table, shardID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*memstore.TableShard)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(table, shardID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleIngestion provides a mock function with given fields: table, shardID, upsertBatch
func (_m *MemStore) HandleIngestion(table string, shardID int, upsertBatch *common.UpsertBatch) error {
	ret := _m.Called(table, shardID, upsertBatch)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int, *common.UpsertBatch) error); ok {
		r0 = rf(table, shardID, upsertBatch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitShards provides a mock function with given fields: schedulerOff, shardOwner
func (_m *MemStore) InitShards(schedulerOff bool, shardOwner topology.ShardOwner) {
	_m.Called(schedulerOff, shardOwner)
}

// Lock provides a mock function with given fields:
func (_m *MemStore) Lock() {
	_m.Called()
}

// Purge provides a mock function with given fields: table, shardID, batchIDStart, batchIDEnd, reporter
func (_m *MemStore) Purge(table string, shardID int, batchIDStart int, batchIDEnd int, reporter memstore.PurgeJobDetailReporter) error {
	ret := _m.Called(table, shardID, batchIDStart, batchIDEnd, reporter)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int, int, int, memstore.PurgeJobDetailReporter) error); ok {
		r0 = rf(table, shardID, batchIDStart, batchIDEnd, reporter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RLock provides a mock function with given fields:
func (_m *MemStore) RLock() {
	_m.Called()
}

// RUnlock provides a mock function with given fields:
func (_m *MemStore) RUnlock() {
	_m.Called()
}

// RemoveTableShard provides a mock function with given fields: table, shardID
func (_m *MemStore) RemoveTableShard(table string, shardID int) {
	_m.Called(table, shardID)
}

// Snapshot provides a mock function with given fields: table, shardID, reporter
func (_m *MemStore) Snapshot(table string, shardID int, reporter memstore.SnapshotJobDetailReporter) error {
	ret := _m.Called(table, shardID, reporter)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int, memstore.SnapshotJobDetailReporter) error); ok {
		r0 = rf(table, shardID, reporter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unlock provides a mock function with given fields:
func (_m *MemStore) Unlock() {
	_m.Called()
}
