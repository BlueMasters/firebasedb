package firebasedb

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func TestSecret(t *testing.T) {
	s := Secret{Token: "password"}
	assert.Equal(t, "password", s.String())
	assert.Error(t, s.Renew())
}

func TestJwt(t *testing.T) {
	uid := uuid()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"v":   0,
		"iat": time.Now().Unix(),
		"d":   map[string]interface{}{"uid": uid},
	})
	tokenString, err := token.SignedString([]byte(testingDbSecret))
	assert.NoError(t, err)

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
	root := db.Auth(Secret{Token: tokenString}).Ref(uuid())
	err = root.Child("pikachu").Set(&pika)
	assert.NoError(t, err)

	p2 := pokemon{}
	err = root.Child("pikachu").Value(&p2)
	assert.NoError(t, err)
	assert.Equal(t, p2.CP, pika.CP)

	err = root.Remove()
	assert.NoError(t, err)
}

func TestBadJwt(t *testing.T) {
	uid := uuid()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"v":     0,
		"iat":   time.Now().Unix(),
		"d":     map[string]interface{}{"uid": uid},
		"debug": true,
	})
	tokenString, err := token.SignedString([]byte("BAD-SECRET"))
	assert.NoError(t, err)

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
	root := db.Auth(Secret{Token: tokenString}).Debug(os.Stderr).Ref(uuid())
	err = root.Child("pikachu").Set(&pika)
	assert.Error(t, err)
}

type jwtToken struct {
	uid        string
	key        string
	str        string
	renewCount int
}

func (t *jwtToken) String() string {
	return t.str
}

func (t *jwtToken) ParamName() string {
	return "auth"
}

func (t *jwtToken) Renew() error {
	t.renewCount += 1
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"v":   0,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(4 * time.Second).Unix(),
		"d":   map[string]interface{}{"uid": t.uid},
	})
	tokenString, err := token.SignedString([]byte(t.key))
	if err != nil {
		return err
	}
	t.str = tokenString
	return nil
}

func TestAuthRevoked(t *testing.T) {
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
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		time.Sleep(2 * time.Second)
		change := map[string]interface{}{"combat_point": 370}
		err := root.Child("pikachu").Update(&change)
		assert.NoError(t, err)
		time.Sleep(5 * time.Second)
		change = map[string]interface{}{"combat_point": 375}
		err = root.Child("pikachu").Update(&change)
		assert.NoError(t, err)
		wg.Done()
	}()

	rt := jwtToken{
		uid: uuid(),
		key: testingDbSecret,
	}
	err = rt.Renew()
	assert.NoError(t, err)

	s, err := root.Auth(&rt).Child("pikachu").Subscribe()
	assert.NoError(t, err)
outer:
	for {
		select {
		case e := <-s.Events():
			if e.Type == "put" {
				var d map[string]interface{}
				_, err := e.Value(&d)
				assert.NoError(t, err)
				if d["combat_point"].(float64) > 372.0 {
					break outer
				}
			}

		case <-time.After(15 * time.Second):
			assert.Fail(t, "Timeout")
			break outer
		}
	}

	assert.True(t, rt.renewCount > 1)
	s.Close()
	wg.Wait()
	err = root.Remove()
	assert.NoError(t, err)

}
