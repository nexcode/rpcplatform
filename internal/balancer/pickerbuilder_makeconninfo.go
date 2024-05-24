/*
 * Copyright 2022 RPCPlatform Authors
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

package balancer

import (
	"cmp"
	"math"
	"slices"
	"strconv"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type connInfo struct {
	subConn  balancer.SubConn
	priority int
	weight   int
	factor   int
	count    int
}

func (pb *pickerBuilder) makeConnInfo(pickerInfo base.PickerBuildInfo) ([]*connInfo, int) {
	connInfoArr := make([]*connInfo, 0, len(pickerInfo.ReadySCs))
	var totalWeight int

	for subConn, subConnInfo := range pickerInfo.ReadySCs {
		connInfo := connInfo{
			subConn: subConn,
		}

		priority, _ := subConnInfo.Address.Attributes.Value("balancerPriority").(string)
		connInfo.priority, _ = strconv.Atoi(priority)

		weight, _ := subConnInfo.Address.Attributes.Value("balancerWeight").(string)
		connInfo.weight, _ = strconv.Atoi(weight)

		if connInfo.weight <= 0 {
			connInfo.weight = 1
		}

		connInfoArr = append(connInfoArr, &connInfo)
	}

	slices.SortFunc(connInfoArr, func(a, b *connInfo) int {
		return cmp.Compare(b.priority, a.priority)
	})

	if pb.maxActiveServers > 0 && pb.maxActiveServers < len(connInfoArr) {
		connInfoArr = connInfoArr[:pb.maxActiveServers]
	}

	for _, connInfo := range connInfoArr {
		connInfo.factor = int(math.Ceil(float64(connInfo.weight) / float64(len(connInfoArr))))
		totalWeight += connInfo.weight
	}

	return connInfoArr, totalWeight
}
