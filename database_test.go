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
	"net/url"
	"testing"
	"fmt"
	"crypto/rand"
)

const dinoFactsUrl = "https://dinosaur-facts.firebaseio.com/"

type dinosaurFact struct {
	Appeared int64   `json:"appeared"`
	Height   float32 `json:"height"`
	Length   float32 `json:"length"`
	Order    string  `json:"order"`
	Vanished int64   `json:"vanished"`
	Weight   int32   `json:"weight"`
}

type dinosaurs map[string]dinosaurFact
type dinoScores map[string]int

func uuid() string {
	u := [16]byte{}
	_, err := rand.Read(u[:16])
	if err != nil {
		panic(err)
	}
	u[8] = (u[8] | 0x80) & 0xBf
	u[6] = (u[6] | 0x40) & 0x4f
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func TestRefOperators(t *testing.T) {
	db, err := NewFirebaseDB(dinoFactsUrl, "")
	assert.NoError(t, err)
	dino := db.Ref("/dinosaurs")
	assert.Equal(t, db.Key(), "")
	pterodactyl := dino.Child("pterodactyl")
	assert.Equal(t, dino.Key(), "dinosaurs")
	assert.Equal(t, pterodactyl.Key(), "pterodactyl")
	assert.Equal(t, pterodactyl.jsonUrl(), "https://dinosaur-facts.firebaseio.com/dinosaurs/pterodactyl.json")
	assert.Equal(t, pterodactyl.Root().jsonUrl(), "https://dinosaur-facts.firebaseio.com/.json")
	assert.Equal(t, pterodactyl.Parent().jsonUrl(), "https://dinosaur-facts.firebaseio.com/dinosaurs.json")
}

func TestArguments(t *testing.T) {
	db, err := NewFirebaseDB(dinoFactsUrl, "")
	assert.NoError(t, err)

	pretty := db.Pretty()
	v := pretty.url.Query()
	assert.Contains(t, v, "print")
	assert.Equal(t, v["print"], []string{"pretty"})

	silent := db.Silent()
	v = silent.url.Query()
	assert.Contains(t, v, "print")
	assert.Equal(t, v["print"], []string{"silent"})

	export := pretty.Export()
	v = export.url.Query()
	assert.Contains(t, v, "print")
	assert.Equal(t, v["print"], []string{"pretty"})
	assert.Contains(t, v, "format")
	assert.Equal(t, v["format"], []string{"export"})
}

func TestDinoShallow(t *testing.T) {
	db, err := NewFirebaseDB(dinoFactsUrl, "")
	assert.NoError(t, err)
	generic := make(map[string]interface{})
	err = db.Ref("/").Shallow().Value(&generic)
	assert.NoError(t, err)
	assert.Contains(t, generic, "dinosaurs")
	assert.Contains(t, generic, "scores")
	assert.NotContains(t, generic, "fat")
	assert.True(t, generic["dinosaurs"].(bool))
	assert.True(t, generic["scores"].(bool))
}

func TestRefFromUrl(t *testing.T) {
	db, err := NewFirebaseDB(dinoFactsUrl, "")
	assert.NoError(t, err)
	generic := make(map[string]interface{})
	u, err := url.Parse("https://dinosaur-facts.firebaseio.com/dinosaurs")
	assert.NoError(t, err)
	r, err := db.RefFromUrl(*u)
	assert.NoError(t, err)
	err = r.Shallow().Value(&generic)
	assert.NoError(t, err)
	assert.Contains(t, generic, "pterodactyl")
	assert.True(t, generic["pterodactyl"].(bool))
	u, err = url.Parse("https://not-the-same-host.firebaseio.com/dinosaurs")
	assert.NoError(t, err)
	_, err = db.RefFromUrl(*u)
	assert.Error(t, err)

}

func TestDotUrl(t *testing.T) {
	db, err := NewFirebaseDB(dinoFactsUrl, "")
	assert.NoError(t, err)
	generic := make(map[string]interface{})
	r := db.Ref("/")
	r.url.Path = "." // force dot path
	err = r.Shallow().Value(&generic)
	assert.NoError(t, err)
	assert.Contains(t, generic, "dinosaurs")
	assert.NotContains(t, generic, "fat")
	assert.True(t, generic["dinosaurs"].(bool))
}

func TestDino(t *testing.T) {
	db, err := NewFirebaseDB(dinoFactsUrl, "")
	assert.NoError(t, err)
	var dinos = dinosaurs{}
	err = db.Ref("/dinosaurs").Value(&dinos)
	assert.NoError(t, err)
	assert.Contains(t, dinos, "pterodactyl")
	assert.NotContains(t, dinos, "pikachu")
	assert.EqualValues(t, dinos["pterodactyl"].Appeared, -150000000)
	assert.EqualValues(t, dinos["pterodactyl"].Order, "pterosauria")
	var scores = dinoScores{}
	err = db.Ref("/scores").Value(&scores)
	assert.NoError(t, err)
	assert.Contains(t, scores, "pterodactyl")
	assert.EqualValues(t, scores["pterodactyl"], 93)
}

func TestBadUrl(t *testing.T) {
	_, err := NewFirebaseDB(":", "")
	assert.Error(t, err)
}

func TestBasic(t *testing.T) {
	db, err := NewFirebaseDB(testingDbUrl, testingDbSecret)
	if err != nil {
		t.Fatal(err)
	}
	pika := struct {
		Name         string `json:"name"`
		CombatPoints int    `json:"combat_point"`
	}{
		Name:         "Pikachu",
		CombatPoints: 450,
	}
	err = db.Ref("/pikachu").Shallow().Set(&pika, nil)
	assert.NoError(t, err)
}
