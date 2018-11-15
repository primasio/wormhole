/* * Copyright 2018 Primas Lab Foundation
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
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/models"
)

func TestIncrementCommentVote(t *testing.T) {
	user := models.User{CommentUpVotes: 30, CommentDownVotes: 20}

	like := true
	user.IncrementCommentVote(like)
	assert.Equal(t, user.CommentUpVotes, uint(31))
	assert.Equal(t, user.CommentDownVotes, uint(20))

	like = false
	user.IncrementCommentVote(like)
	assert.Equal(t, user.CommentUpVotes, uint(31))
	assert.Equal(t, user.CommentDownVotes, uint(21))
}

func TestSwitchCommentVote(t *testing.T) {
	user := models.User{CommentUpVotes: 30, CommentDownVotes: 20}

	like := true
	user.SwitchCommentVote(like)
	assert.Equal(t, user.CommentUpVotes, uint(31))
	assert.Equal(t, user.CommentDownVotes, uint(19))

	like = false
	user.SwitchCommentVote(like)
	assert.Equal(t, user.CommentUpVotes, uint(30))
	assert.Equal(t, user.CommentDownVotes, uint(20))
}

func TestCancelCommentVote(t *testing.T) {
	user := models.User{CommentUpVotes: 30, CommentDownVotes: 20}

	like := true
	user.CancelCommentVote(like)
	assert.Equal(t, user.CommentUpVotes, uint(29))
	assert.Equal(t, user.CommentDownVotes, uint(20))

	like = false
	user.CancelCommentVote(like)
	assert.Equal(t, user.CommentUpVotes, uint(29))
	assert.Equal(t, user.CommentDownVotes, uint(19))
}

func TestIncrementIntegration(t *testing.T) {
	user := models.User{CommentUpVotes: 30, CommentDownVotes: 20, Integration: 55}

	user.IncrementIntegration(5)
	assert.Equal(t, user.Integration, int64(60))

	user.IncrementIntegration(-5)
	assert.Equal(t, user.Integration, int64(55))
}
