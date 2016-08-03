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
	"strconv"
)

func (r Reference) OrderBy(childKey string) Reference {
	return r.withQuotedParam("orderBy", childKey)
}

func (r Reference) OrderByKey() Reference {
	return r.OrderBy("$key")
}

func (r Reference) OrderByValue() Reference {
	return r.OrderBy("$value")
}

func (r Reference) LimitToFirst(n uint64) Reference {
	return r.withParam("limitToFirst", strconv.FormatUint(n, 10))
}

func (r Reference) LimitToLast(n uint64) Reference {
	return r.withParam("limitToLast", strconv.FormatUint(n, 10))
}

func (r Reference) StartAt(n interface{}) Reference {
	return r.withQuotedParam("startAt", n)
}

func (r Reference) EndAt(n interface{}) Reference {
	return r.withQuotedParam("endAt", n)
}

func (r Reference) EqualTo(n interface{}) Reference {
	return r.withQuotedParam("equalTo", n)
}
