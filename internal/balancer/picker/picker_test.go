/*
 * Copyright 2025 RPCPlatform Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package picker

import (
	"slices"
	"testing"

	"github.com/nexcode/rpcplatform/internal/attributes"
	"github.com/nexcode/rpcplatform/internal/config"
	"github.com/nexcode/rpcplatform/internal/grpcattrs"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/endpointsharding"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

func TestPicker(t *testing.T) {
	t.Parallel()

	childStates := []endpointsharding.ChildState{{
		State: balancer.State{
			ConnectivityState: connectivity.Ready,
			Picker:            &namedPicker{name: 1},
		},
		Endpoint: resolver.Endpoint{
			Attributes: grpcattrs.SetAttributes(nil, &attributes.Attributes{
				BalancerWeight:   1,
				BalancerPriority: 1,
			}),
		},
	}, {
		State: balancer.State{
			ConnectivityState: connectivity.Ready,
			Picker:            &namedPicker{name: 2},
		},
		Endpoint: resolver.Endpoint{
			Attributes: grpcattrs.SetAttributes(nil, &attributes.Attributes{
				BalancerWeight:   7,
				BalancerPriority: 2,
			}),
		},
	}, {
		State: balancer.State{
			ConnectivityState: connectivity.Ready,
			Picker:            &namedPicker{name: 3},
		},
		Endpoint: resolver.Endpoint{
			Attributes: grpcattrs.SetAttributes(nil, &attributes.Attributes{
				BalancerWeight:   5,
				BalancerPriority: 2,
			}),
		},
	}, {
		State: balancer.State{
			ConnectivityState: connectivity.Ready,
			Picker:            &namedPicker{name: 4},
		},
		Endpoint: resolver.Endpoint{
			Attributes: grpcattrs.SetAttributes(nil, &attributes.Attributes{
				BalancerWeight:   3,
				BalancerPriority: 2,
			}),
		},
	}, {
		State: balancer.State{
			ConnectivityState: connectivity.Ready,
			Picker:            &namedPicker{name: 5},
		},
		Endpoint: resolver.Endpoint{
			Attributes: grpcattrs.SetAttributes(nil, &attributes.Attributes{
				BalancerWeight:   0,
				BalancerPriority: 3,
			}),
		},
	}, {
		State: balancer.State{
			ConnectivityState: connectivity.TransientFailure,
			Picker:            &namedPicker{name: 6},
		},
		Endpoint: resolver.Endpoint{
			Attributes: grpcattrs.SetAttributes(nil, &attributes.Attributes{
				BalancerWeight:   1,
				BalancerPriority: 3,
			}),
		},
	}}

	config := &config.Client{
		MaxActiveServers: 3,
	}

	picker := New(childStates, config).(*picker)
	actualSequence := make([]int, len(picker.pickers))
	pickerNext := picker.next

	for i, childPicker := range picker.pickers {
		actualSequence[i] = childPicker.(*namedPicker).name

		if _, err := picker.Pick(balancer.PickInfo{}); err != nil {
			t.Fatalf("Pick() failed: %v", err)
		}

		pickerNext++
		if pickerNext == len(picker.pickers) {
			pickerNext = 0
		}

		if picker.next != pickerNext {
			t.Errorf("picker next = %v, want: %v", picker.next, pickerNext)
		}
	}

	expectedSequence := []int{2, 3, 4, 2, 3, 2, 2, 3, 4, 2, 3, 2, 2, 3, 4}

	if !slices.Equal(actualSequence, expectedSequence) {
		t.Errorf("picker sequence = %v, want: %v", actualSequence, expectedSequence)
	}
}

type namedPicker struct {
	name int
}

func (namedPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{}, nil
}
