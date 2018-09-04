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

package oauth

import (
	"context"
	"encoding/json"
	"github.com/primasio/wormhole/config"
	"github.com/primasio/wormhole/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

type GoogleUserInfoResponse struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func getGoogleOAuthConfig() *oauth2.Config {
	if googleOAuthConfig != nil {
		return googleOAuthConfig
	}

	c := config.GetConfig()

	clientId := c.GetString("oauth.google.client_id")
	clientSecret := c.GetString("oauth.google.client_secret")

	scheme := c.GetString("application.scheme")
	domain := c.GetString("application.domain")

	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		RedirectURL:  scheme + "://" + domain + "/v1/oauth/callback/google",
	}

	return googleOAuthConfig
}

func HandleGoogleAuthCallback(code string) (err error, userId uint) {

	googleConfig := getGoogleOAuthConfig()

	// 1. Use code to get Google access token

	token, e := googleConfig.Exchange(context.Background(), code)

	if e != nil {
		return e, 0
	}

	// 2. Use access token to get user info

	url := "https://www.googleapis.com/oauth2/v2/userinfo"
	client := googleConfig.Client(context.Background(), token)
	resp, e := client.Get(url)

	if e != nil {
		return e, 0
	}

	defer resp.Body.Close()

	userInfo := &GoogleUserInfoResponse{}

	if e := json.NewDecoder(resp.Body).Decode(userInfo); e != nil {
		return e, 0
	}

	// 3. Process user info

	result := &OAuthResult{
		Type:      models.OAuthGoogle,
		Id:        userInfo.Id,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
	}

	if err, userId := result.Process(); err != nil {
		return err, 0
	} else {
		return nil, userId
	}
}

func HandleGoogleAuth(state string) (redirectUrl string) {

	url := getGoogleOAuthConfig().AuthCodeURL(state)

	return url
}
