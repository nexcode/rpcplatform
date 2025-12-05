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
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"sync"

	"github.com/nexcode/rpcplatform/internal/config"
	"github.com/nexcode/rpcplatform/internal/grpcattrs"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/balancer/endpointsharding"
	"google.golang.org/grpc/connectivity"
)

func New(childStates []endpointsharding.ChildState, config *config.Client) balancer.Picker {
	if len(childStates) == 0 {
		return base.NewErrPicker(errNoServerAvailableForPick)
	}

	var connecting bool
	var totalWeight int

	pickerStates := make([]*state, 0, len(childStates))

	for _, childState := range childStates {
		attributes := grpcattrs.GetAttributes(childState.Endpoint.Attributes)
		if attributes.BalancerWeight <= 0 {
			continue
		}

		if childState.State.ConnectivityState == connectivity.Connecting {
			connecting = true
		}

		if childState.State.ConnectivityState != connectivity.Ready {
			continue
		}

		pickerStates = append(pickerStates, &state{
			picker:   childState.State.Picker,
			priority: attributes.BalancerPriority,
			weight:   attributes.BalancerWeight,
		})
	}

	if len(pickerStates) == 0 {
		if connecting {
			return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
		}

		return base.NewErrPicker(errNoServerAvailableForPick)
	}

	slices.SortFunc(pickerStates, func(a, b *state) int {
		return cmp.Compare(b.priority, a.priority)
	})

	if config.MaxActiveServers > 0 && config.MaxActiveServers < len(pickerStates) {
		pickerStates = pickerStates[:config.MaxActiveServers]
	}

	for _, pickerState := range pickerStates {
		pickerState.factor = int(math.Ceil(float64(pickerState.weight) / float64(len(pickerStates))))
		totalWeight += pickerState.weight
	}

	picker := &picker{
		pickers: make([]balancer.Picker, 0, totalWeight),
	}

	for {
		prevLen := len(picker.pickers)

		for _, pickerState := range pickerStates {
			if pickerState.count < pickerState.factor && pickerState.weight > 0 {
				picker.pickers = append(picker.pickers, pickerState.picker)
				pickerState.weight--
				pickerState.count++
			}
		}

		if totalWeight == len(picker.pickers) {
			break
		}

		if prevLen == len(picker.pickers) {
			for _, pickerState := range pickerStates {
				pickerState.count = 0
			}
		}
	}

	picker.next = rand.IntN(totalWeight)
	return picker
}

type picker struct {
	pickers []balancer.Picker
	mu      sync.Mutex
	next    int
}
