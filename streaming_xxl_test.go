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
var result chan string

func startReceiver(t *testing.T, r Reference, wg *sync.WaitGroup, n int) {
	s, err := r.Child("live").Subscribe(false, true)
	allSubscriptions[n] = s
	assert.NoError(t, err)
	wg.Done()
	for e := range s.Events() {
		var x string
		p, err := e.Value(&x)
		assert.NoError(t, err)
		assert.Equal(t, "/", p)
		result <- x
	}
}

func startSender(t *testing.T, r Reference, wg *sync.WaitGroup, n int, nobjs int) {
	for i := 0; i < nobjs; i++ {
		objectId := fmt.Sprintf("XXL-%06d-%06d", n, i)
		err := r.Child("live").Set(&objectId, nil)
		assert.NoError(t, err)
		data := map[string]string{"seen": "yes"}
		err = r.Child("historical").Child(objectId).Set(&data, nil)
		assert.NoError(t, err)
	}
	wg.Done()
}

func TestStreamXXL(t *testing.T) {
	const numberOfReceivers = 10
	const numberOfSenders = 5
	const numberOfObjects = 3

	result = make(chan string)

	allSubscriptions = make([]Subscription, numberOfReceivers)
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
		go startSender(t, root, ready, i, numberOfObjects)
	}
	ready.Wait()

	checker := map[string]int{}
	counter := 0
	var finished <-chan time.Time

	outer:
	for {
		select {
		case x := <-result:
			counter++
			i, ok := checker[x]
			if ok {
				checker[x] = i+1
			} else {
				checker[x] = 1
			}
			if counter == (numberOfSenders * numberOfObjects + 1) * numberOfReceivers {
				// all received... wait 1 second more
				finished = time.After(1 * time.Second)
			}
		case <- finished:
			break outer
		case <- time.After(5 * time.Second):
			assert.Fail(t, "timeout!")
			break outer
		}
	}

	assert.Equal(t, (numberOfSenders * numberOfObjects + 1) * numberOfReceivers, counter)
	assert.Contains(t, checker, "")
	assert.EqualValues(t, checker[""], numberOfReceivers)
	assert.Len(t, checker, numberOfSenders * numberOfObjects + 1)

	for i := 0; i < numberOfReceivers; i++ {
		allSubscriptions[i].Close()
	}

	err = root.Remove()
	assert.NoError(t, err)
}
