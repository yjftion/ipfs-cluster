// Package local implements a ClusterDAGService that chunks and adds content
// to a local peer, before pinning it.
package local

import (
	"context"
	"errors"

	adder "github.com/ipfs/ipfs-cluster/adder"
	"github.com/ipfs/ipfs-cluster/api"

	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	rpc "gx/ipfs/QmTfA73jjmEphGCYGYyZksqy4vRKdv9sKJLKb6WzbCBqJB/go-libp2p-gorpc"
	peer "gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
	ipld "gx/ipfs/QmcKKBwfz6FyQdHR2jsXrrF6XeSBXYL86anmWNewpFpoF5/go-ipld-format"
	logging "gx/ipfs/QmcuXC5cxs79ro2cUuHs4HQ2bkDLJUYokwL8aivcX6HW3C/go-log"
)

var errNotFound = errors.New("dagservice: block not found")

var logger = logging.Logger("localdags")

// DAGService is an implementation of an adder.ClusterDAGService which
// puts the added blocks directly in the peers allocated to them (without
// sharding).
type DAGService struct {
	adder.BaseDAGService

	rpcClient *rpc.Client

	dests   []peer.ID
	pinOpts api.PinOptions
}

// New returns a new Adder with the given rpc Client. The client is used
// to perform calls to IPFSBlockPut and Pin content on Cluster.
func New(rpc *rpc.Client, opts api.PinOptions) *DAGService {
	return &DAGService{
		rpcClient: rpc,
		dests:     nil,
		pinOpts:   opts,
	}
}

// Add puts the given node in the destination peers.
func (dgs *DAGService) Add(ctx context.Context, node ipld.Node) error {
	if dgs.dests == nil {
		dests, err := adder.BlockAllocate(ctx, dgs.rpcClient, dgs.pinOpts)
		if err != nil {
			return err
		}
		dgs.dests = dests
	}

	size, err := node.Size()
	if err != nil {
		return err
	}
	nodeSerial := &api.NodeWithMeta{
		Cid:     node.Cid().String(),
		Data:    node.RawData(),
		CumSize: size,
	}

	return adder.PutBlock(ctx, dgs.rpcClient, nodeSerial, dgs.dests)
}

// Finalize pins the last Cid added to this DAGService.
func (dgs *DAGService) Finalize(ctx context.Context, root cid.Cid) (cid.Cid, error) {
	// Cluster pin the result
	rootPin := api.PinWithOpts(root, dgs.pinOpts)
	rootPin.Allocations = dgs.dests

	dgs.dests = nil

	return root, dgs.rpcClient.CallContext(
		ctx,
		"",
		"Cluster",
		"Pin",
		rootPin.ToSerial(),
		&struct{}{},
	)
}

// AddMany calls Add for every given node.
func (dgs *DAGService) AddMany(ctx context.Context, nodes []ipld.Node) error {
	for _, node := range nodes {
		err := dgs.Add(ctx, node)
		if err != nil {
			return err
		}
	}
	return nil
}
