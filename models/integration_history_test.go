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

package models_test

import (
	"crypto/sha1"
	"fmt"
	"io"
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/primasio/wormhole/models"
)

func TestSetUniqueID(t *testing.T) {
	history := models.IntegrationHistory{}

	err := history.SetUniqueID()
	if err == nil {
		t.Errorf("err should be nil, but got error: %s", err)
	}

	history.Data = "Hello, there"
	history.SetUniqueID()

	h := sha1.New()
	io.WriteString(h, history.Data)
	uid := fmt.Sprintf("%x", h.Sum(nil))
	assert.Equal(t, uid, history.UniqueID)
}
