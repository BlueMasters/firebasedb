// Copyright 2016 Jacques Supcik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firebasedb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"sync"
	"fmt"
	"time"
)

var allSubscriptions []Subscription

func startReceiver(t *testing.T, r Reference, wg *sync.WaitGroup, n int) {
	s, err := r.Child("live").Subscribe(false, true)
	allSubscriptions[n] = s
	assert.NoError(t, err)

	wg.Done()
}

func startSender(t *testing.T, r Reference, wg *sync.WaitGroup, n int) {
	const numberOfObjects = 10
	for i := 0; i < numberOfObjects; i++ {
		objectId := fmt.Sprintf("XXL-%06d-%06d", n, i)
		err := r.Child("live").Set(&objectId, nil)
		assert.NoError(t, err)
		data := map[string]string{"seen": "yes"}
		err = r.Child("historical").Child(objectId).Set(&data, nil)
		fmt.Println(r.Child("historical").Child(objectId).jsonUrl())
		assert.NoError(t, err)
	}
	wg.Done()
}

func TestStreamXXL(t *testing.T) {
	const numberOfReceivers = 20
	const numberOfSenders = 10
	db, err := NewFirebaseDB(testingDbUrl, testingDbSecret)
	assert.NoError(t, err)
	root := db.Ref(uuid())

	ready := &sync.WaitGroup{}
	ready.Add(numberOfReceivers)
	for i := 0; i < numberOfReceivers; i++ {
		go startReceiver(t, root, ready, i)
	}
	ready.Wait()

	ready = &sync.WaitGroup{}
	ready.Add(numberOfSenders)
	for i := 0; i < numberOfSenders; i++ {
		go startSender(t, root, ready, i)
	}
	ready.Wait()
	time.Sleep(1 * time.Second) // give some time to Firebase

	for i := 0; i < numberOfReceivers; i++ {
		allSubscriptions[i].Close()
	}
}
