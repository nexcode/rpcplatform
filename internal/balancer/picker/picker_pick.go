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
	"google.golang.org/grpc/balancer"
)

func (p *picker) Pick(pickInfo balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	picker := p.pickers[p.next]
	p.next = (p.next + 1) % len(p.pickers)
	p.mu.Unlock()

	return picker.Pick(pickInfo)
}
