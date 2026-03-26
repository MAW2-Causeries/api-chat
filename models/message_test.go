package models

import (
	"testing"
	"time"

	"github.com/bouk/monkey"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestMessageToMapWithNonNilDeletedAt(t *testing.T) {
	deletedAt := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	message := &Message{
		ID:        "1",
		Content:   "Hello, World!",
		AuthorID:  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
		ChannelID: "27731CCA-ADB5-42DB-AA8C-500994FC4098",
		CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		DeletedAt: &deletedAt,
	}
	expectedMap := map[string]any{
		"id":         "1",
		"content":    "Hello, World!",
		"author_id":  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
		"channel_id": "27731CCA-ADB5-42DB-AA8C-500994FC4098",
		"created_at": "2024-06-01T12:00:00Z",
		"updated_at": "2024-06-01T12:00:00Z",
		"deleted_at": "2024-06-01T12:00:00Z",
	}

	assert.Equal(t, expectedMap, message.ToMap())
}

func TestMessageToMap(t *testing.T) {
	message := &Message{
		ID:        "1",
		Content:   "Hello, World!",
		AuthorID:  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
		ChannelID: "27731CCA-ADB5-42DB-AA8C-500994FC4098",
		CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		DeletedAt: nil,
	}
	expectedMap := map[string]any{
		"id":         "1",
		"content":    "Hello, World!",
		"author_id":  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
		"channel_id": "27731CCA-ADB5-42DB-AA8C-500994FC4098",
		"created_at": "2024-06-01T12:00:00Z",
		"updated_at": "2024-06-01T12:00:00Z",
		"deleted_at": nil,
	}

	assert.Equal(t, expectedMap, message.ToMap())
}

func TestNewMessage(t *testing.T) {
	var channelID = "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	var content = "Hello, World!"
	var authorID = "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF"

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Exec, func(_ *gocql.Query) error {
		return nil
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		now := time.Now().UTC()
		if len(dest) >= 1 {
			if p, ok := dest[0].(*string); ok {
				*p = "1"
			}
		}
		if len(dest) >= 2 {
			if p, ok := dest[1].(*time.Time); ok {
				*p = now
			}
		}
		if len(dest) >= 3 {
			if p, ok := dest[2].(*time.Time); ok {
				*p = now
			}
		}
		if len(dest) >= 4 {
			if p, ok := dest[3].(**time.Time); ok {
				*p = nil
			}
		}
		return nil
	})
	
	var message = NewMessage(authorID, channelID, content)

	assert.Equal(t, message.Content, content)
	assert.Equal(t, message.AuthorID, authorID)
	assert.Equal(t, message.ChannelID, channelID)
	assert.Equal(t, "1", message.ID)
	assert.Nil(t, message.DeletedAt)
	assert.NotEmpty(t, message.CreatedAt)
	assert.NotEmpty(t, message.UpdatedAt)
}

func TestNewMessageWithDBError(t *testing.T) {
	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Exec, func(_ *gocql.Query) error {
		return assert.AnError
	})
	
	var message = NewMessage("authorID", "channelID", "content")
	assert.Nil(t, message)
}

func TestGetMessagesByChannelId(t *testing.T) {
	var channelID = "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	var limit = 2
	var page = 2
	var fakeMessages = []*Message{
		{
			ID:        "1",
			Content:   "Hello, World!",
			AuthorID:  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
			ChannelID: channelID,
			CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			DeletedAt: nil,
		},
		{
			ID:        "2",
			Content:   "Hello, again!",
			AuthorID:  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
			ChannelID: channelID,
			CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			DeletedAt: nil,
		},
		{
			ID:        "3",
			Content:   "Hello, once more!",
			AuthorID:  "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
			ChannelID: channelID,
			CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			DeletedAt: nil,
		},
		{
			ID:        "4",
			Content:   "Hello, once more!",
			AuthorID: "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF",
			ChannelID: channelID,
			CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			DeletedAt: nil,
		},
	}

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...any) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Iter, func(_ *gocql.Query) *gocql.Iter {
		return &gocql.Iter{}
	})
	// return one row on the first Scan call, then signal end of iteration
	iterCalledCounter := 0
	monkey.Patch((*gocql.Iter).Scan, func(_ *gocql.Iter, dest ...any) bool {
		if iterCalledCounter < len(fakeMessages) {
			message := fakeMessages[iterCalledCounter]
			fields := []interface{}{
				&message.ID,
				&message.Content,
				&message.AuthorID,
				&message.ChannelID,
				&message.CreatedAt,
				&message.UpdatedAt,
				&message.DeletedAt,
			}
			for i, field := range fields {
				if i < len(dest) {
					switch v := field.(type) {
					case *string:
						if p, ok := dest[i].(*string); ok {
							*p = *v
						}
					case *time.Time:
						if p, ok := dest[i].(*time.Time); ok {
							*p = *v
						}
					case **time.Time:
						if p, ok := dest[i].(**time.Time); ok {
							*p = *v
						}
					}
				}
			}
		}
		iterCalledCounter++
		if iterCalledCounter > 4 {
			return false
		}

		return true
	})

	messages := GetMessagesByChannelID(channelID, limit, page)

	assert.Len(t, messages, 2)
	assert.Equal(t, messages[0].ID, "3")
	assert.Equal(t, messages[1].ID, "4")
	assert.Equal(t, messages[0].Content, "Hello, once more!")
	assert.Equal(t, messages[1].Content, "Hello, once more!")
	assert.Equal(t, messages[0].AuthorID, "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF")
	assert.Equal(t, messages[1].AuthorID, "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF")
	assert.Equal(t, messages[0].ChannelID, channelID)
	assert.Equal(t, messages[1].ChannelID, channelID)
	assert.Nil(t, messages[0].DeletedAt)
	assert.Nil(t, messages[1].DeletedAt)
	assert.NotEmpty(t, messages[0].CreatedAt)
	assert.NotEmpty(t, messages[1].CreatedAt)
	assert.NotEmpty(t, messages[0].UpdatedAt)
	assert.NotEmpty(t, messages[1].UpdatedAt)
}

func TestGetMessageByChannelIDAndMessageID(t *testing.T) {
	var channelID = "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	var messageID = "1"

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...any) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		now := time.Now().UTC()
		if len(dest) >= 1 {
			if p, ok := dest[0].(*string); ok {
				*p = messageID
			}
		}
		if len(dest) >= 2 {
			if p, ok := dest[1].(*string); ok {
				*p = "Hello, World!"
			}
		}
		if len(dest) >= 3 {
			if p, ok := dest[2].(*string); ok {
				*p = "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF"
			}
		}
		if len(dest) >= 4 {
			if p, ok := dest[3].(*string); ok {
				*p = channelID
			}
		}
		if len(dest) >= 5 {
			if p, ok := dest[4].(*time.Time); ok {
				*p = now
			}
		}
		if len(dest) >= 6 {
			if p, ok := dest[5].(*time.Time); ok {
				*p = now
			}
		}
		if len(dest) >= 7 {
			if p, ok := dest[6].(**time.Time); ok {
				*p = nil
			}
		}
		return nil
	})
	
	message := GetMessageByChannelIDAndMessageID(channelID, messageID)

	assert.Equal(t, message.ID, messageID)
	assert.Equal(t, message.Content, "Hello, World!")
	assert.Equal(t, message.AuthorID, "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF")
	assert.Equal(t, message.ChannelID, channelID)
	assert.Nil(t, message.DeletedAt)
	assert.NotEmpty(t, message.CreatedAt)
	assert.NotEmpty(t, message.UpdatedAt)
}

func TestGetMessageByChannelIDAndMessageIDNotFound(t *testing.T) {
	var channelID = "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	var messageID = "nonexistent"


	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...any) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		return assert.AnError
	})
	
	message := GetMessageByChannelIDAndMessageID(channelID, messageID)
	assert.Nil(t, message)
}