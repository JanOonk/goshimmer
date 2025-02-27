package common

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/goshimmer/tools/integration-tests/tester/framework"
	"github.com/iotaledger/goshimmer/tools/integration-tests/tester/framework/config"
	"github.com/iotaledger/goshimmer/tools/integration-tests/tester/tests"
)

// TestCommonSynchronization checks whether messages are relayed through the network,
// a node that joins later solidifies, stop and start this node again, and whether all messages
// are available on all nodes at the end (persistence).
func TestCommonSynchronization(t *testing.T) {
	const (
		initialPeers    = 2
		numMessages     = 100
		numSyncMessages = 5 * initialPeers
	)

	ctx, cancel := tests.Context(context.Background(), t)
	defer cancel()
	n, err := f.CreateNetwork(ctx, t.Name(), initialPeers, framework.CreateNetworkConfig{
		StartSynced: true,
	})
	require.NoError(t, err)
	defer tests.ShutdownNetwork(ctx, t, n)

	// 1. issue data messages
	log.Printf("Issuing %d messages to sync...", numMessages)
	ids := tests.SendDataMessages(t, n.Peers(), numMessages)
	log.Println("Issuing messages... done")

	// 2. spawn peer without knowledge of previous messages
	log.Println("Spawning new node to sync...")
	newPeer, err := n.CreatePeer(ctx, createNewPeerConfig())
	require.NoError(t, err)
	err = n.DoManualPeering(ctx)
	require.NoError(t, err)
	log.Println("Spawning new node... done")

	// 3. issue some messages on old peers so that new peer can solidify
	log.Printf("Issuing %d messages on the %d initial peers...", numSyncMessages, initialPeers)
	ids = tests.SendDataMessages(t, n.Peers()[:initialPeers], numSyncMessages, ids)
	log.Println("Issuing messages... done")

	// 4. check whether all issued messages are available on to the new peer
	tests.RequireMessagesAvailable(t, []*framework.Node{newPeer}, ids, time.Minute, tests.Tick)
	tests.RequireMessagesEqual(t, []*framework.Node{newPeer}, ids)
	require.True(t, tests.Synced(t, newPeer))

	// 5. shut down newly added peer
	log.Println("Stopping new node...")
	require.NoError(t, newPeer.Stop(ctx))
	log.Println("Stopping new node... done")

	log.Printf("Issuing %d messages and waiting until they have old tangle time...", numMessages)
	ids = tests.SendDataMessages(t, n.Peers()[:initialPeers], numMessages, ids)
	// wait to assure that the new peer is actually out of sync when starting
	time.Sleep(newPeer.Config().MessageLayer.TangleTimeWindow)
	log.Println("Issuing messages... done")

	// 6. let it startup again
	log.Println("Restarting new node to sync again...")
	err = newPeer.Start(ctx)
	require.NoError(t, err)
	err = n.DoManualPeering(ctx)
	require.NoError(t, err)
	log.Println("Restarting node... done")

	// the node should not be in sync as all the message are outside its sync time window
	require.False(t, tests.Synced(t, newPeer))

	// 7. issue some messages on old peers so that new peer can sync again
	log.Printf("Issuing %d messages on the %d initial peers...", numSyncMessages, initialPeers)
	ids = tests.SendDataMessages(t, n.Peers()[:initialPeers], numSyncMessages, ids)
	log.Println("Issuing messages... done")

	// 9. check whether all issued messages are available on to the new peer
	tests.RequireMessagesAvailable(t, []*framework.Node{newPeer}, ids, time.Minute, tests.Tick)
	tests.RequireMessagesEqual(t, []*framework.Node{newPeer}, ids)
	// check that the new node is synced
	require.Eventuallyf(t,
		func() bool { return tests.Synced(t, newPeer) },
		tests.Timeout, tests.Tick,
		"the peer %s did not sync again after restart", newPeer)
}

func createNewPeerConfig() config.GoShimmer {
	conf := framework.PeerConfig
	// the new peer should use a shorter TangleTimeWindow than regular peers to go out of sync before them
	conf.MessageLayer.TangleTimeWindow = 30 * time.Second
	return conf
}
