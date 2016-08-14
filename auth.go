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

// Note: with jwt (e.g. using github.com/dgrijalva/jwt-go), you get
// the following info:
// auth={
//     "uid":"<UID>",
//     "token":{
//         "aud":"<DATABASE_NAME>",
//         "sub":"<UID>",
//         "iss":"https://securetoken.google.com/<DATABASE_NAME>",
//         "iat":<UNIX_TIME>,
//         "auth_time":<UNIX_TIME>,
//         "exp":<UNIX_TIME>
//     }
// }

package firebasedb

import "errors"

// Authenticator is the interface used to add authentication data to the requests. The String() method
// returns the current token and Renew() is called if the current token has expired.
type Authenticator interface {
	String() string
	ParamName() string // usually "auth" ou "access_token"
	Renew() error
}

// Secret implements the Authenticator interface and is used with static Database secret.
type Secret struct {
	Token string
}

// String returns the static Database secret
func (s Secret) String() string {
	return s.Token
}

// ParamName returns "auth"
func (s Secret) ParamName() string {
	return "auth"
}

// Renew is not allowed for static secret and thus always returns an error.
func (s Secret) Renew() error {
	return errors.New("Can't renew a static token")
}
