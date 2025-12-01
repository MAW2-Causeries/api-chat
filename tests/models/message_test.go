package tests

import (
	"MessagesService/models"
	"testing"
	"time"

	"github.com/bouk/monkey"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

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
	
	var message = models.NewMessage(authorID, channelID, content)

	assert.Equal(t, message.Content, content)
	assert.Equal(t, message.AuthorID, authorID)
	assert.Equal(t, message.ChannelID, channelID)
	assert.Equal(t, "1", message.ID)
	assert.Nil(t, message.DeletedAt)
	assert.NotEmpty(t, message.CreatedAt)
	assert.NotEmpty(t, message.UpdatedAt)
}

func TestGetMessagesByChannelId(t *testing.T) {
	var channelID = "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	var limit = 50

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...any) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Iter, func(_ *gocql.Query) *gocql.Iter {
		return &gocql.Iter{}
	})
	// return one row on the first Scan call, then signal end of iteration
	iterCalled := false
	monkey.Patch((*gocql.Iter).Scan, func(_ *gocql.Iter, dest ...any) bool {
		if !iterCalled {
			if p, ok := dest[0].(*string); ok {
				*p = "1"
			}
			if p, ok := dest[1].(*string); ok {
				*p = "Hello, World!"
			}
			if p, ok := dest[2].(*string); ok {
				*p = "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF"
			}
			if p, ok := dest[3].(*string); ok {
				*p = channelID
			}
			if p, ok := dest[4].(*time.Time); ok {
				*p = time.Now().UTC()
			}
			if p, ok := dest[5].(*time.Time); ok {
				*p = time.Now().UTC()
			}
			if p, ok := dest[6].(**time.Time); ok {
				*p = nil
			}
			iterCalled = true
			return true
		}
		return false // end of iteration after one record
	})

	messages, err := models.GetMessagesByChannelID(channelID, limit)

	assert.Nil(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, messages[0].ChannelID, channelID)
	assert.Equal(t, messages[0].Content, "Hello, World!")
	assert.Equal(t, messages[0].AuthorID, "A71F1694-D9AD-40E4-94D2-A59D07E9D9AF")
	assert.Equal(t, messages[0].ID, "1")
	assert.Nil(t, messages[0].DeletedAt)
	assert.NotEmpty(t, messages[0].CreatedAt)
	assert.NotEmpty(t, messages[0].UpdatedAt)
}