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

package captcha

import (
	"context"
	"encoding/json"
	"github.com/primasio/wormhole/config"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const endpoint = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaVerifyResponse struct {
	Success     bool
	ChallengeTs time.Time
	Hostname    string
}

func VerifyRecaptchaToken(token string) (error, bool) {

	secret := config.GetConfig().GetString("recaptcha.secret")

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", token)

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFn()

	resp, err := ctxhttp.Do(ctx, nil, req)

	if err != nil {
		return err, false
	}

	defer resp.Body.Close()

	recaptchaResponse := &RecaptchaVerifyResponse{}

	if err := json.NewDecoder(resp.Body).Decode(recaptchaResponse); err != nil {
		return err, false
	}

	if !recaptchaResponse.Success {
		return nil, false
	}

	duration := time.Now().Unix() - recaptchaResponse.ChallengeTs.Unix()

	if duration <= 0 || duration >= int64(time.Minute*10) {
		return nil, false
	}

	return nil, true
}
