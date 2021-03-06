package ipfscluster

import (
	"context"

	peer "gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"

	"github.com/ipfs/ipfs-cluster/api"
)

// RPCAPI is a go-libp2p-gorpc service which provides the internal ipfs-cluster
// API, which enables components and cluster peers to communicate and
// request actions from each other.
//
// The RPC API methods are usually redirects to the actual methods in
// the different components of ipfs-cluster, with very little added logic.
// Refer to documentation on those methods for details on their behaviour.
type RPCAPI struct {
	c *Cluster
}

/*
   Cluster components methods
*/

// ID runs Cluster.ID()
func (rpcapi *RPCAPI) ID(ctx context.Context, in struct{}, out *api.IDSerial) error {
	id := rpcapi.c.ID().ToSerial()
	*out = id
	return nil
}

// Pin runs Cluster.Pin().
func (rpcapi *RPCAPI) Pin(ctx context.Context, in api.PinSerial, out *struct{}) error {
	return rpcapi.c.Pin(in.ToPin())
}

// Unpin runs Cluster.Unpin().
func (rpcapi *RPCAPI) Unpin(ctx context.Context, in api.PinSerial, out *struct{}) error {
	c := in.DecodeCid()
	return rpcapi.c.Unpin(c)
}

// Pins runs Cluster.Pins().
func (rpcapi *RPCAPI) Pins(ctx context.Context, in struct{}, out *[]api.PinSerial) error {
	cidList := rpcapi.c.Pins()
	cidSerialList := make([]api.PinSerial, 0, len(cidList))
	for _, c := range cidList {
		cidSerialList = append(cidSerialList, c.ToSerial())
	}
	*out = cidSerialList
	return nil
}

// PinGet runs Cluster.PinGet().
func (rpcapi *RPCAPI) PinGet(ctx context.Context, in api.PinSerial, out *api.PinSerial) error {
	cidarg := in.ToPin()
	pin, err := rpcapi.c.PinGet(cidarg.Cid)
	if err == nil {
		*out = pin.ToSerial()
	}
	return err
}

// Version runs Cluster.Version().
func (rpcapi *RPCAPI) Version(ctx context.Context, in struct{}, out *api.Version) error {
	*out = api.Version{
		Version: rpcapi.c.Version(),
	}
	return nil
}

// Peers runs Cluster.Peers().
func (rpcapi *RPCAPI) Peers(ctx context.Context, in struct{}, out *[]api.IDSerial) error {
	peers := rpcapi.c.Peers()
	var sPeers []api.IDSerial
	for _, p := range peers {
		sPeers = append(sPeers, p.ToSerial())
	}
	*out = sPeers
	return nil
}

// PeerAdd runs Cluster.PeerAdd().
func (rpcapi *RPCAPI) PeerAdd(ctx context.Context, in string, out *api.IDSerial) error {
	pid, _ := peer.IDB58Decode(in)
	id, err := rpcapi.c.PeerAdd(pid)
	*out = id.ToSerial()
	return err
}

// ConnectGraph runs Cluster.GetConnectGraph().
func (rpcapi *RPCAPI) ConnectGraph(ctx context.Context, in struct{}, out *api.ConnectGraphSerial) error {
	graph, err := rpcapi.c.ConnectGraph()
	*out = graph.ToSerial()
	return err
}

// PeerRemove runs Cluster.PeerRm().
func (rpcapi *RPCAPI) PeerRemove(ctx context.Context, in peer.ID, out *struct{}) error {
	return rpcapi.c.PeerRemove(in)
}

// Join runs Cluster.Join().
func (rpcapi *RPCAPI) Join(ctx context.Context, in api.MultiaddrSerial, out *struct{}) error {
	addr := in.ToMultiaddr()
	err := rpcapi.c.Join(addr)
	return err
}

// StatusAll runs Cluster.StatusAll().
func (rpcapi *RPCAPI) StatusAll(ctx context.Context, in struct{}, out *[]api.GlobalPinInfoSerial) error {
	pinfos, err := rpcapi.c.StatusAll()
	*out = GlobalPinInfoSliceToSerial(pinfos)
	return err
}

// StatusAllLocal runs Cluster.StatusAllLocal().
func (rpcapi *RPCAPI) StatusAllLocal(ctx context.Context, in struct{}, out *[]api.PinInfoSerial) error {
	pinfos := rpcapi.c.StatusAllLocal()
	*out = pinInfoSliceToSerial(pinfos)
	return nil
}

// Status runs Cluster.Status().
func (rpcapi *RPCAPI) Status(ctx context.Context, in api.PinSerial, out *api.GlobalPinInfoSerial) error {
	c := in.DecodeCid()
	pinfo, err := rpcapi.c.Status(c)
	*out = pinfo.ToSerial()
	return err
}

// StatusLocal runs Cluster.StatusLocal().
func (rpcapi *RPCAPI) StatusLocal(ctx context.Context, in api.PinSerial, out *api.PinInfoSerial) error {
	c := in.DecodeCid()
	pinfo := rpcapi.c.StatusLocal(c)
	*out = pinfo.ToSerial()
	return nil
}

// SyncAll runs Cluster.SyncAll().
func (rpcapi *RPCAPI) SyncAll(ctx context.Context, in struct{}, out *[]api.GlobalPinInfoSerial) error {
	pinfos, err := rpcapi.c.SyncAll()
	*out = GlobalPinInfoSliceToSerial(pinfos)
	return err
}

// SyncAllLocal runs Cluster.SyncAllLocal().
func (rpcapi *RPCAPI) SyncAllLocal(ctx context.Context, in struct{}, out *[]api.PinInfoSerial) error {
	pinfos, err := rpcapi.c.SyncAllLocal()
	*out = pinInfoSliceToSerial(pinfos)
	return err
}

// Sync runs Cluster.Sync().
func (rpcapi *RPCAPI) Sync(ctx context.Context, in api.PinSerial, out *api.GlobalPinInfoSerial) error {
	c := in.DecodeCid()
	pinfo, err := rpcapi.c.Sync(c)
	*out = pinfo.ToSerial()
	return err
}

// SyncLocal runs Cluster.SyncLocal().
func (rpcapi *RPCAPI) SyncLocal(ctx context.Context, in api.PinSerial, out *api.PinInfoSerial) error {
	c := in.DecodeCid()
	pinfo, err := rpcapi.c.SyncLocal(c)
	*out = pinfo.ToSerial()
	return err
}

// RecoverAllLocal runs Cluster.RecoverAllLocal().
func (rpcapi *RPCAPI) RecoverAllLocal(ctx context.Context, in struct{}, out *[]api.PinInfoSerial) error {
	pinfos, err := rpcapi.c.RecoverAllLocal()
	*out = pinInfoSliceToSerial(pinfos)
	return err
}

// Recover runs Cluster.Recover().
func (rpcapi *RPCAPI) Recover(ctx context.Context, in api.PinSerial, out *api.GlobalPinInfoSerial) error {
	c := in.DecodeCid()
	pinfo, err := rpcapi.c.Recover(c)
	*out = pinfo.ToSerial()
	return err
}

// RecoverLocal runs Cluster.RecoverLocal().
func (rpcapi *RPCAPI) RecoverLocal(ctx context.Context, in api.PinSerial, out *api.PinInfoSerial) error {
	c := in.DecodeCid()
	pinfo, err := rpcapi.c.RecoverLocal(c)
	*out = pinfo.ToSerial()
	return err
}

// BlockAllocate returns allocations for blocks. This is used in the adders.
// It's different from pin allocations when ReplicationFactor < 0.
func (rpcapi *RPCAPI) BlockAllocate(ctx context.Context, in api.PinSerial, out *[]string) error {
	pin := in.ToPin()
	err := rpcapi.c.setupPin(&pin)
	if err != nil {
		return err
	}

	// Return the current peer list.
	if pin.ReplicationFactorMin < 0 {
		// Returned metrics are Valid and belong to current
		// Cluster peers.
		metrics := rpcapi.c.monitor.LatestMetrics(pingMetricName)
		peers := make([]string, len(metrics), len(metrics))
		for i, m := range metrics {
			peers[i] = peer.IDB58Encode(m.Peer)
		}

		*out = peers
		return nil
	}

	allocs, err := rpcapi.c.allocate(
		pin.Cid,
		pin.ReplicationFactorMin,
		pin.ReplicationFactorMax,
		[]peer.ID{}, // blacklist
		[]peer.ID{}, // prio list
	)

	if err != nil {
		return err
	}

	*out = api.PeersToStrings(allocs)
	return nil
}

// SendInformerMetric runs Cluster.sendInformerMetric().
func (rpcapi *RPCAPI) SendInformerMetric(ctx context.Context, in struct{}, out *api.Metric) error {
	m, err := rpcapi.c.sendInformerMetric()
	*out = m
	return err
}

/*
   Tracker component methods
*/

// Track runs PinTracker.Track().
func (rpcapi *RPCAPI) Track(ctx context.Context, in api.PinSerial, out *struct{}) error {
	return rpcapi.c.tracker.Track(in.ToPin())
}

// Untrack runs PinTracker.Untrack().
func (rpcapi *RPCAPI) Untrack(ctx context.Context, in api.PinSerial, out *struct{}) error {
	c := in.DecodeCid()
	return rpcapi.c.tracker.Untrack(c)
}

// TrackerStatusAll runs PinTracker.StatusAll().
func (rpcapi *RPCAPI) TrackerStatusAll(ctx context.Context, in struct{}, out *[]api.PinInfoSerial) error {
	*out = pinInfoSliceToSerial(rpcapi.c.tracker.StatusAll())
	return nil
}

// TrackerStatus runs PinTracker.Status().
func (rpcapi *RPCAPI) TrackerStatus(ctx context.Context, in api.PinSerial, out *api.PinInfoSerial) error {
	c := in.DecodeCid()
	pinfo := rpcapi.c.tracker.Status(c)
	*out = pinfo.ToSerial()
	return nil
}

// TrackerRecoverAll runs PinTracker.RecoverAll().f
func (rpcapi *RPCAPI) TrackerRecoverAll(ctx context.Context, in struct{}, out *[]api.PinInfoSerial) error {
	pinfos, err := rpcapi.c.tracker.RecoverAll()
	*out = pinInfoSliceToSerial(pinfos)
	return err
}

// TrackerRecover runs PinTracker.Recover().
func (rpcapi *RPCAPI) TrackerRecover(ctx context.Context, in api.PinSerial, out *api.PinInfoSerial) error {
	c := in.DecodeCid()
	pinfo, err := rpcapi.c.tracker.Recover(c)
	*out = pinfo.ToSerial()
	return err
}

/*
   IPFS Connector component methods
*/

// IPFSPin runs IPFSConnector.Pin().
func (rpcapi *RPCAPI) IPFSPin(ctx context.Context, in api.PinSerial, out *struct{}) error {
	c := in.DecodeCid()
	depth := in.ToPin().MaxDepth
	return rpcapi.c.ipfs.Pin(ctx, c, depth)
}

// IPFSUnpin runs IPFSConnector.Unpin().
func (rpcapi *RPCAPI) IPFSUnpin(ctx context.Context, in api.PinSerial, out *struct{}) error {
	c := in.DecodeCid()
	return rpcapi.c.ipfs.Unpin(ctx, c)
}

// IPFSPinLsCid runs IPFSConnector.PinLsCid().
func (rpcapi *RPCAPI) IPFSPinLsCid(ctx context.Context, in api.PinSerial, out *api.IPFSPinStatus) error {
	c := in.DecodeCid()
	b, err := rpcapi.c.ipfs.PinLsCid(ctx, c)
	*out = b
	return err
}

// IPFSPinLs runs IPFSConnector.PinLs().
func (rpcapi *RPCAPI) IPFSPinLs(ctx context.Context, in string, out *map[string]api.IPFSPinStatus) error {
	m, err := rpcapi.c.ipfs.PinLs(ctx, in)
	*out = m
	return err
}

// IPFSConnectSwarms runs IPFSConnector.ConnectSwarms().
func (rpcapi *RPCAPI) IPFSConnectSwarms(ctx context.Context, in struct{}, out *struct{}) error {
	err := rpcapi.c.ipfs.ConnectSwarms()
	return err
}

// IPFSConfigKey runs IPFSConnector.ConfigKey().
func (rpcapi *RPCAPI) IPFSConfigKey(ctx context.Context, in string, out *interface{}) error {
	res, err := rpcapi.c.ipfs.ConfigKey(in)
	*out = res
	return err
}

// IPFSRepoStat runs IPFSConnector.RepoStat().
func (rpcapi *RPCAPI) IPFSRepoStat(ctx context.Context, in struct{}, out *api.IPFSRepoStat) error {
	res, err := rpcapi.c.ipfs.RepoStat()
	*out = res
	return err
}

// IPFSSwarmPeers runs IPFSConnector.SwarmPeers().
func (rpcapi *RPCAPI) IPFSSwarmPeers(ctx context.Context, in struct{}, out *api.SwarmPeersSerial) error {
	res, err := rpcapi.c.ipfs.SwarmPeers()
	*out = res.ToSerial()
	return err
}

// IPFSBlockPut runs IPFSConnector.BlockPut().
func (rpcapi *RPCAPI) IPFSBlockPut(ctx context.Context, in api.NodeWithMeta, out *struct{}) error {
	return rpcapi.c.ipfs.BlockPut(in)
}

// IPFSBlockGet runs IPFSConnector.BlockGet().
func (rpcapi *RPCAPI) IPFSBlockGet(ctx context.Context, in api.PinSerial, out *[]byte) error {
	c := in.DecodeCid()
	res, err := rpcapi.c.ipfs.BlockGet(c)
	*out = res
	return err
}

/*
   Consensus component methods
*/

// ConsensusLogPin runs Consensus.LogPin().
func (rpcapi *RPCAPI) ConsensusLogPin(ctx context.Context, in api.PinSerial, out *struct{}) error {
	c := in.ToPin()
	return rpcapi.c.consensus.LogPin(c)
}

// ConsensusLogUnpin runs Consensus.LogUnpin().
func (rpcapi *RPCAPI) ConsensusLogUnpin(ctx context.Context, in api.PinSerial, out *struct{}) error {
	c := in.ToPin()
	return rpcapi.c.consensus.LogUnpin(c)
}

// ConsensusAddPeer runs Consensus.AddPeer().
func (rpcapi *RPCAPI) ConsensusAddPeer(ctx context.Context, in peer.ID, out *struct{}) error {
	return rpcapi.c.consensus.AddPeer(in)
}

// ConsensusRmPeer runs Consensus.RmPeer().
func (rpcapi *RPCAPI) ConsensusRmPeer(ctx context.Context, in peer.ID, out *struct{}) error {
	return rpcapi.c.consensus.RmPeer(in)
}

// ConsensusPeers runs Consensus.Peers().
func (rpcapi *RPCAPI) ConsensusPeers(ctx context.Context, in struct{}, out *[]peer.ID) error {
	peers, err := rpcapi.c.consensus.Peers()
	*out = peers
	return err
}

/*
   PeerMonitor
*/

// PeerMonitorLogMetric runs PeerMonitor.LogMetric().
func (rpcapi *RPCAPI) PeerMonitorLogMetric(ctx context.Context, in api.Metric, out *struct{}) error {
	rpcapi.c.monitor.LogMetric(in)
	return nil
}

// PeerMonitorLatestMetrics runs PeerMonitor.LatestMetrics().
func (rpcapi *RPCAPI) PeerMonitorLatestMetrics(ctx context.Context, in string, out *[]api.Metric) error {
	*out = rpcapi.c.monitor.LatestMetrics(in)
	return nil
}
