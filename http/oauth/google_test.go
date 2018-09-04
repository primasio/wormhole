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

package oauth_test

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/http/oauth"
	"testing"
)

func TestResponseUnmarshal(t *testing.T) {
	jsonStr := `{
    "id": "12345",
        "email": "test@gmail.com",
        "verified_email": true,
        "name": "chen zhao",
        "given_name": "chen",
        "family_name": "zhao",
        "link": "https://plus.google.com/105681419613076020853",
        "picture": "https://lh5.googleusercontent.com/-YdfbEoumJVE/AAAAAAAAAAI/AAAAAAAAAAo/Iiq2f7lNVFY/photo.jpg",
        "locale": "zh-CN"
}`
	response := oauth.GoogleUserInfoResponse{}

	err := json.Unmarshal([]byte(jsonStr), &response)

	assert.Equal(t, err, nil)

	assert.Equal(t, response.Id, "12345")
	assert.Equal(t, response.Email, "test@gmail.com")
}
