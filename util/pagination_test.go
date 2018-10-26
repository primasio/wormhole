/*
 * Copyright 2018 Primas Lab Foundation
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

package util_test

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/util"
)

func TestPaginate(t *testing.T) {
	type Temp struct {
		Hello string `json:"hello"`
	}

	temp := &Temp{Hello: "world"}
	p := util.Paginate(1, 10, 200, temp)
	assert.Equal(t, p.CurrentPage, uint(1), "Current page should be 1")
	assert.Equal(t, p.From, uint(1), "From should be 1")
	assert.Equal(t, p.To, uint(10), "To Sholud be 10")
	assert.Equal(t, p.Data, temp)

	p = util.Paginate(1, 10, 9, temp)
	assert.Equal(t, p.CurrentPage, uint(1))
	assert.Equal(t, p.From, uint(1))
	assert.Equal(t, p.To, uint(9))
	assert.Equal(t, p.Data, temp)

	p = util.Paginate(2, 10, 9, temp)
	assert.Equal(t, p.CurrentPage, uint(2))
	assert.Equal(t, p.From, uint(0))
	assert.Equal(t, p.To, uint(0))
	assert.Equal(t, p.Data, temp)
}

func TestCanPaginate(t *testing.T) {
	var page, pageSize, count uint
	page = 0
	pageSize = 20
	count = 100

	pass := util.CanPaginate(page, pageSize, count)
	assert.Equal(t, pass, true)

	page = 2
	pageSize = 20
	count = 10
	pass = util.CanPaginate(page, pageSize, count)
	assert.Equal(t, pass, false)
}

func TestPurePageArgs(t *testing.T) {
	var page, pageSize uint
	page, pageSize = util.PurePageArgs(page, pageSize)
	assert.Equal(t, page, uint(1))
	assert.Equal(t, pageSize, uint(20))

	page = 99
	pageSize = 20000

	page, pageSize = util.PurePageArgs(page, pageSize)
	assert.Equal(t, page, uint(99))
	assert.Equal(t, pageSize, uint(100))
}

func TestEmptyPagination(t *testing.T) {
	p := util.EmptyPagination(0, 20)
	d, err := json.Marshal(p.Data)

	assert.Equal(t, err, nil)
	assert.Equal(t, string(d), "[]")
}
