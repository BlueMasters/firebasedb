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
    "testing"
    "os"
    "log"
)

var (
    testingDbUrl string
    testingDbSecret string
)

func TestMain(m *testing.M) {
    testingDbUrl = os.Getenv("FIREBASE_DB_TESTING_URL")
    if (testingDbUrl == "") {
        log.Fatal("Please set the 'FIREBASE_DB_TESTING_URL' environment variable with the URL of your database")
    }
    testingDbSecret = os.Getenv("FIREBASE_DB_TESTING_SECRET")
    if (testingDbSecret == "") {
        log.Fatal("Please set the 'FIREBASE_DB_TESTING_SECRET' environment variable with the secret token of your database")
    }
    agree := os.Getenv("FIREBASE_DB_TESTING_I_UNDERSTAND_THAT_THIS_WILL_DELETE_EXISTING_DATA")
    if (agree != "I AGREE") {
        log.Fatal("Please set the 'FIREBASE_DB_TESTING_I_UNDERSTAND_THAT_THIS_WILL_DELETE_EXISTING_DATA' to 'I AGREE'")
    }
    os.Exit(m.Run())
}