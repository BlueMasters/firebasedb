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
	"time"
)

func TestStream(t *testing.T) {
	db, err := NewFirebaseDB(testingDbUrl, testingDbSecret)
	if err != nil {
		t.Fatal(err)
	}
	type pokemon struct {
		Name string `json:"name"`
		CP   int    `json:"combat_point"`
	}

	root := db.Ref(uuid())
	pika := pokemon{
		Name: "Pikachu",
		CP:   365,
	}
	err = root.Child("pikachu").Set(&pika, nil)
	assert.NoError(t, err)

	s, err := root.Subscribe()
	assert.NoError(t, err)

	select {
	case e := <-s.Events():
		assert.Equal(t, e.Type, "put")
		p := map[string]pokemon{}
		path, err := e.Value(&p)
		assert.NoError(t, err)
		assert.Equal(t, "/", path)
		assert.Contains(t, p, "pikachu")
		assert.Equal(t, "Pikachu", p["pikachu"].Name)
	case <-time.After(5 * time.Second):
		assert.Fail(t, "Got Timeout instead of first event")
	}

    select {
    case  <-s.Events():
        assert.Fail(t, "Got a second event!")
    case <-time.After(1 * time.Second):
        // pass
    }

	p2 := pokemon{}
	err = root.Child("pikachu").Value(&p2)
	assert.NoError(t, err)
	assert.Equal(t, p2.CP, pika.CP)

	err = s.Close()
	assert.NoError(t, err)

	err = root.Remove()
	assert.NoError(t, err)

	generic := map[string]interface{}{}
	err = root.Value(&generic)


}
