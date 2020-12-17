package batchservice

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/bee/pkg/postage"
)

type chainState struct {
	Block uint64   `json:"block"` // The block number of the last postage event
	Total *big.Int `json:"total"` // Cumulative amount paid per stamp
	Price *big.Int `json:"price"` // Bzz/chunk/block normalised price
}

// BatchService implements EventUpdater.
type BatchService struct {
	cs     chainState
	storer postage.BatchStorer
	logger logging.Logger
}

// BatchServiceOption is an option passed to NewBatchService
type BatchServiceOption func(*BatchService)

// NewBatchService will create a new BatchService
func NewBatchService(storer postage.BatchStorer, logger logging.Logger, opts ...BatchServiceOption) postage.EventUpdater {
	b := BatchService{
		storer: storer,
		logger: logger,
	}

	for _, opt := range opts {
		opt(&b)
	}

	return &b
}

// Create will create a new batch and store it in the BatchStore.
func (svc *BatchService) Create(id, owner []byte, value *big.Int, depth uint8) error {
	b := postage.Batch{
		ID:    id,
		Owner: owner,
		Value: value,
		Start: svc.cs.Block,
		Depth: depth,
	}

	err := svc.storer.Put(&b)
	if err != nil {
		return fmt.Errorf("CreateBatch: %w", err)
	}

	svc.logger.Debugf("created batch id %x", hex.EncodeToString(b.ID))
	return nil
}

// TopUp implements the EventUpdater interface. It tops ups a batch with the
// given ID with the given amount of BZZ.
func (svc *BatchService) TopUp(id []byte, amount *big.Int) error {
	b, err := svc.storer.Get(id)
	if err != nil {
		return fmt.Errorf("TopUp: %w", err)
	}

	b.Value.Add(b.Value, amount)

	err = svc.storer.Put(b)
	if err != nil {
		return fmt.Errorf("TopUp: %w", err)
	}

	svc.logger.Debugf("topped up batch id %x with %v", hex.EncodeToString(b.ID), b.Value)
	return nil
}

// UpdateDepth implements the EventUpdater inteface. It sets the new depth of a
// batch with the given ID.
func (svc *BatchService) UpdateDepth(id []byte, depth uint8) error {
	b, err := svc.storer.Get(id)
	if err != nil {
		return err
	}

	b.Depth = depth

	err = svc.storer.Put(b)
	if err != nil {
		return fmt.Errorf("update depth: %w", err)
	}

	svc.logger.Debugf("updated depth of batch id %x to %d", hex.EncodeToString(b.ID), b.Depth)
	return nil
}

// UpdatePrice implements the EventUpdater interface. It sets the current
// price from the chain in the service chain state.
func (svc *BatchService) UpdatePrice(price *big.Int) error {
	svc.cs.Price = price

	svc.logger.Debugf("updated chain price to %s", svc.cs.Price)
	return nil
}
