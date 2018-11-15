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

package service

import (
	"sync"

	"github.com/primasio/wormhole/config"
)

var integration *Integration
var integrationOnce sync.Once

type Integration struct{}

func GetIntegration() *Integration {
	integrationOnce.Do(func() {
		integration = &Integration{}
	})

	return integration
}

func (s *Integration) GetURLContentCommentVoteScore(like bool) int64 {
	c := config.GetConfig()

	if like {
		return c.GetInt64("integration.url_content_comment_vote.like")
	}

	return c.GetInt64("integration.url_content_comment_vote.hate")
}

func (s *Integration) GetRegisterScore() int64 {
	return config.GetConfig().GetInt64("integration.register")
}
