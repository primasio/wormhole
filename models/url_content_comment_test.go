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

func TestIncrementVote(t *testing.T) {
	comment := models.URLContentComment{CommentUpVotes: 30, CommentDownVotes: 20}

	like := true
	comment.IncrementVote(like)
	assert.Equal(t, comment.CommentUpVotes, uint(31))
	assert.Equal(t, comment.CommentDownVotes, uint(20))

	like = false
	comment.IncrementVote(like)
	assert.Equal(t, comment.CommentUpVotes, uint(31))
	assert.Equal(t, comment.CommentDownVotes, uint(21))
}

func TestSwitchVote(t *testing.T) {
	comment := models.URLContentComment{CommentUpVotes: 30, CommentDownVotes: 20}

	like := true
	comment.SwitchVote(like)
	assert.Equal(t, comment.CommentUpVotes, uint(31))
	assert.Equal(t, comment.CommentDownVotes, uint(19))

	like = false
	comment.SwitchVote(like)
	assert.Equal(t, comment.CommentUpVotes, uint(30))
	assert.Equal(t, comment.CommentDownVotes, uint(20))
}

func TestCancelVote(t *testing.T) {
	comment := models.URLContentComment{CommentUpVotes: 30, CommentDownVotes: 20}

	like := true
	comment.CancelVote(like)
	assert.Equal(t, comment.CommentUpVotes, uint(29))
	assert.Equal(t, comment.CommentDownVotes, uint(20))

	like = false
	comment.CancelVote(like)
	assert.Equal(t, comment.CommentUpVotes, uint(29))
	assert.Equal(t, comment.CommentDownVotes, uint(19))
}
