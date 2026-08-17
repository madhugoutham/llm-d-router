/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package predictedlatency

import (
	"testing"

	attrtopology "github.com/llm-d/llm-d-router/pkg/epp/framework/plugins/datalayer/attribute/topology"
	latencypredictor "github.com/llm-d/llm-d-router/pkg/epp/framework/plugins/requestcontrol/dataproducer/predictedlatency/latencypredictorclient"
	"github.com/stretchr/testify/require"
)

func TestTopologyDistanceAndScore(t *testing.T) {
	peer := &attrtopology.Topology{Hostname: "host-a", Rack: "rack-a", Zone: "zone-a", Region: "region-a"}
	candidate := &attrtopology.Topology{Hostname: "host-a", Rack: "rack-a", Zone: "zone-a", Region: "region-a"}

	distance, score := topologyDistanceAndScore(peer, candidate, true, true)
	require.Equal(t, "host", distance)
	require.Equal(t, 1.0, score)

	distance, score = topologyDistanceAndScore(nil, candidate, false, true)
	require.Equal(t, "unknown", distance)
	require.Zero(t, score)

	other := &attrtopology.Topology{Hostname: "host-b", Rack: "rack-b", Zone: "zone-b", Region: "region-b"}
	distance, score = topologyDistanceAndScore(peer, other, true, true)
	require.Equal(t, "none", distance)
	require.Zero(t, score)
}

func TestStampTopologyOnEntry(t *testing.T) {
	peerKnown, candidateKnown := true, true
	ctx := &predictedLatencyCtx{
		topologyDistance:       "rack",
		topologyAffinityScore:  0.2,
		peerTopologyKnown:      peerKnown,
		candidateTopologyKnown: candidateKnown,
		requestID:              "request-1",
	}
	entry := latencypredictor.TrainingEntry{}

	stampTopologyOnEntry(&entry, ctx)

	require.Equal(t, "request-1", entry.RequestID)
	require.Equal(t, "rack", entry.TopologyDistance)
	require.Equal(t, 0.2, entry.TopologyAffinityScore)
	require.NotNil(t, entry.PeerTopologyKnown)
	require.NotNil(t, entry.CandidateTopologyKnown)
	require.True(t, *entry.PeerTopologyKnown)
	require.True(t, *entry.CandidateTopologyKnown)
}

func TestStampTopologyOnEntryMonolithic(t *testing.T) {
	entry := latencypredictor.TrainingEntry{}
	stampTopologyOnEntry(&entry, &predictedLatencyCtx{requestID: "request-2"})

	require.Equal(t, "request-2", entry.RequestID)
	require.Empty(t, entry.TopologyDistance)
	require.Nil(t, entry.PeerTopologyKnown)
	require.Nil(t, entry.CandidateTopologyKnown)
}
