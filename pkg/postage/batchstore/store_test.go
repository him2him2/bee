// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package batchstore_test

import (
	"testing"

	"github.com/ethersphere/bee/pkg/postage"
	"github.com/ethersphere/bee/pkg/postage/batchstore"
	postagetest "github.com/ethersphere/bee/pkg/postage/testing"
	"github.com/ethersphere/bee/pkg/statestore/mock"
)

func TestBatchStoreGet(t *testing.T) {
	stateStore := mock.NewStateStore()

	testBatch := postagetest.MustNewBatch()
	key := batchstore.BatchKey(testBatch.ID)

	if err := stateStore.Put(key, testBatch); err != nil {
		t.Fatalf("store batch: %v", err)
	}

	bStore := batchstore.New(stateStore)

	got, err := bStore.Get(testBatch.ID)
	if err != nil {
		t.Fatalf("get batch: %v", err)
	}

	postagetest.CompareBatches(t, testBatch, got)
}

func TestBatchStorePut(t *testing.T) {
	stateStore := mock.NewStateStore()
	bStore := batchstore.New(stateStore)

	testBatch := postagetest.MustNewBatch()
	key := batchstore.BatchKey(testBatch.ID)

	if err := bStore.Put(testBatch); err != nil {
		t.Fatalf("put batch: %v", err)
	}

	var got postage.Batch
	if err := stateStore.Get(key, &got); err != nil {
		t.Fatalf("store get batch: %v", err)
	}

	postagetest.CompareBatches(t, testBatch, &got)
}

func TestBatchStoreGetChainState(t *testing.T) {
	stateStore := mock.NewStateStore()
	bStore := batchstore.New(stateStore)

	testState := postagetest.MustNewChainState()

	err := stateStore.Put(batchstore.StateKey, testState)
	if err != nil {
		t.Fatalf("stateStore put: %v", err)
	}

	got, err := bStore.GetChainState()
	if err != nil {
		t.Fatalf("get chain state: %v", err)
	}

	postagetest.CompareChainState(t, testState, got)
}

func TestBatchStorePutChainState(t *testing.T) {
	stateStore := mock.NewStateStore()
	bStore := batchstore.New(stateStore)

	testState := postagetest.MustNewChainState()

	err := bStore.PutChainState(testState)
	if err != nil {
		t.Fatalf("put chain state: %v", err)
	}

	var got postage.ChainState
	err = stateStore.Get(batchstore.StateKey, &got)
	if err != nil {
		t.Fatalf("statestore get: %v", err)
	}

	postagetest.CompareChainState(t, testState, &got)
}
