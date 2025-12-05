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

package gears

import (
	"sync"
	"testing"
)

func TestUID(t *testing.T) {
	t.Parallel()

	goroutines := 10
	generations := 100

	var wg1, wg2 sync.WaitGroup
	idChan := make(chan string, goroutines*generations)
	idMap := make(map[string]struct{})

	wg1.Add(goroutines)
	wg2.Add(1)

	go func() {
		defer wg2.Done()

		for id := range idChan {
			idMap[id] = struct{}{}
		}
	}()

	for range goroutines {
		go func() {
			defer wg1.Done()

			for range generations {
				idChan <- UID()
			}
		}()
	}

	wg1.Wait()
	close(idChan)
	wg2.Wait()

	if length := len(idMap); length != goroutines*generations {
		t.Errorf("there are %v unique IDs, want: %v", length, goroutines*generations)
	}
}
