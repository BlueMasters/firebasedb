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
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"net/http"
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
	db := NewReference(dinoFactsUrl)
	assert.NoError(t,db.Error)
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
	db := NewReference(dinoFactsUrl)
	assert.NoError(t, db.Error)

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
	db := NewReference(dinoFactsUrl)
	assert.NoError(t, db.Error)
	generic := make(map[string]interface{})
	err := db.Ref("/").Shallow().Value(&generic)
	assert.NoError(t, err)
	assert.Contains(t, generic, "dinosaurs")
	assert.Contains(t, generic, "scores")
	assert.NotContains(t, generic, "fat")
	assert.True(t, generic["dinosaurs"].(bool))
	assert.True(t, generic["scores"].(bool))
}

func TestRefFromUrl(t *testing.T) {
	db := NewReference(dinoFactsUrl)
	assert.NoError(t, db.Error)
	generic := make(map[string]interface{})
	u, err := url.Parse("https://dinosaur-facts.firebaseio.com/dinosaurs")
	assert.NoError(t, err)
	r := db.RefFromUrl(*u)
	assert.NoError(t, r.Error)
	err = r.Shallow().Value(&generic)
	assert.NoError(t, err)
	assert.Contains(t, generic, "pterodactyl")
	assert.True(t, generic["pterodactyl"].(bool))
	u, err = url.Parse("https://not-the-same-host.firebaseio.com/dinosaurs")
	assert.NoError(t, err)
	r = db.RefFromUrl(*u)
	assert.Error(t, r.Error)
}

func TestDotUrl(t *testing.T) {
	db := NewReference(dinoFactsUrl)
	assert.NoError(t, db.Error)
	generic := make(map[string]interface{})
	r := db.Ref("/")
	r.url.Path = "." // force dot path
	err := r.Shallow().Value(&generic)
	assert.NoError(t, err)
	assert.Contains(t, generic, "dinosaurs")
	assert.NotContains(t, generic, "fat")
	assert.True(t, generic["dinosaurs"].(bool))
}

func TestDino(t *testing.T) {
	db := NewReference(dinoFactsUrl)
	assert.NoError(t, db.Error)
	var dinos = dinosaurs{}
	err := db.Ref("/dinosaurs").Value(&dinos)
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

func TestQueries(t *testing.T) {
	db := NewReference(dinoFactsUrl)
	assert.NoError(t, db.Error)
	var dinos = dinosaurs{}
	db.Ref("/dinosaurs").OrderByChild("height").StartAt(3).EndAt(5).Value(&dinos)
	assert.Contains(t, dinos, "triceratops")
	assert.Contains(t, dinos, "stegosaurus")
	assert.NotContains(t, dinos, "lambeosaurus")
	assert.NotContains(t, dinos, "bruhathkayosaurus")

	dinos = dinosaurs{}
	db.Ref("/dinosaurs").OrderByChild("height").EqualTo(uint(4)).Value(&dinos)
	assert.Contains(t, dinos, "stegosaurus")
	assert.NotContains(t, dinos, "triceratops")
	assert.NotContains(t, dinos, "lambeosaurus")
	assert.NotContains(t, dinos, "bruhathkayosaurus")

	dinos = dinosaurs{}
	db.Ref("/dinosaurs").OrderByChild("height").StartAt(2.5).EndAt(4.5).Value(&dinos)
	assert.Contains(t, dinos, "triceratops")
	assert.Contains(t, dinos, "stegosaurus")
	assert.NotContains(t, dinos, "lambeosaurus")
	assert.NotContains(t, dinos, "bruhathkayosaurus")

	dinos = dinosaurs{}
	db.Ref("/dinosaurs").OrderByKey().LimitToFirst(2).Value(&dinos)
	assert.Contains(t, dinos, "bruhathkayosaurus")
	assert.Contains(t, dinos, "lambeosaurus")
	assert.NotContains(t, dinos, "linhenykus")
	assert.NotContains(t, dinos, "stegosaurus")

	scores := dinoScores{}
	db.Ref("/scores").OrderByValue().LimitToLast(2).Value(&scores)
	assert.Contains(t, scores, "linhenykus")
	assert.Contains(t, scores, "pterodactyl")
	assert.NotContains(t, scores, "bruhathkayosaurus")
	assert.NotContains(t, scores, "lambeosaurus")
}

func TestBadUrl(t *testing.T) {
	db := NewReference(":")
	assert.Error(t, db.Error)
}

func TestSet(t *testing.T) {
	db := NewReference(testingDbUrl)
	assert.NoError(t, db.Error)
	type pokemon struct {
		Name string `json:"name"`
		CP   int    `json:"combat_point"`
	}
	pika := pokemon{
		Name: "Pikachu",
		CP:   365,
	}
	root := db.Auth(Secret{Token: testingDbSecret}).Ref(uuid())
	err := root.Child("pikachu").Set(&pika)
	assert.NoError(t, err)

	p2 := pokemon{}
	err = root.Child("pikachu").Value(&p2)
	assert.NoError(t, err)
	assert.Equal(t, p2.CP, pika.CP)

	err = root.Remove()
	assert.NoError(t, err)
}

func TestPatch(t *testing.T) {
	db := NewReference(testingDbUrl)
	assert.NoError(t, db.Error)
	type pokemon struct {
		Name string `json:"name"`
		CP   int    `json:"combat_point"`
	}
	pika := pokemon{
		Name: "Pikachu",
		CP:   365,
	}
	root := db.Auth(Secret{Token: testingDbSecret}).Ref(uuid())
	err := root.Child("pikachu").Set(&pika)
	assert.NoError(t, err)

	res := map[string]interface{}{}
	err = root.Child("pikachu-2").SetWithResult(&pika, &res)
	assert.NoError(t, err)
	assert.Contains(t, res, "combat_point")
	assert.EqualValues(t, 365, res["combat_point"])
	assert.Contains(t, res, "name")


	p2 := pokemon{}
	err = root.Child("pikachu").Value(&p2)
	assert.NoError(t, err)
	assert.Equal(t, pika.CP, p2.CP)

	change := map[string]interface{}{"combat_point": 370}
	err = root.Child("pikachu").Update(&change)
	assert.NoError(t, err)

	res = map[string]interface{}{}
	err = root.Child("pikachu-2").UpdateWithResult(&change, &res)
	assert.NoError(t, err)
	assert.Contains(t, res, "combat_point")
	assert.EqualValues(t, 370, res["combat_point"])
	assert.NotContains(t, res, "name")

	p2 = pokemon{}
	err = root.Child("pikachu").Value(&p2)
	assert.NoError(t, err)
	assert.Equal(t, 370, p2.CP)

	err = root.Remove()
	assert.NoError(t, err)
}

func TestPush(t *testing.T) {
	db := NewReference(testingDbUrl)
	assert.NoError(t, db.Error)
	type pokemon struct {
		Name string `json:"name"`
		CP   int    `json:"combat_point"`
	}
	pika := pokemon{
		Name: "Pikachu",
		CP:   365,
	}
	root := db.Auth(Secret{Token: testingDbSecret}).Ref(uuid())

	_, err := root.Child("pokemons").Push(&pika)
	assert.NoError(t, err)

	bulb := pokemon{
		Name: "Bulbasaur",
		CP:   412,
	}

	_, err = root.Child("pokemons").Push(&bulb)
	assert.NoError(t, err)

	var p map[string]pokemon
	err = root.Child("pokemons").Value(&p)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(p))

	err = root.Remove()
	assert.NoError(t, err)
}

func TestCustomClient(t *testing.T) {
	client := http.Client{}
	db := NewReference(dinoFactsUrl).WithHttpClient(&client)
	assert.NoError(t, db.Error)
	var dinos = dinosaurs{}
	err := db.Ref("/dinosaurs").Value(&dinos)
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