package models

import (
	"MessagesService/databases"
	"time"
)

// Message represents a message in the system
type Message struct {
	ID        string
	Content   string
	AuthorID  string
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// ToMap converts a Message to a map[string]any
func (m *Message) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"author_id":  m.AuthorID,
		"channel_id": m.ChannelID,
		"content":    m.Content,
		"created_at": m.CreatedAt.Format(time.RFC3339),
		"updated_at": m.UpdatedAt.Format(time.RFC3339),
		"deleted_at": func() any {
			if m.DeletedAt == nil {
				return nil
			}
			return m.DeletedAt.Format(time.RFC3339)
		}(),
	}
}

// NewMessage creates a new Message and saves it to the database
func NewMessage(authorID, channelID, content string) *Message {
	var id string
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	if err := databases.Session.Query(`
		INSERT INTO messages (id, content, author_id, channel_id, created_at, updated_at, deleted_at)
		VALUES (uuid(), ?, ?, ?, toTimestamp(now()), toTimestamp(now()), null)`,
		content, authorID, channelID,
	).Exec(); err != nil {
		return nil
	}

	// Retrieve the newly created message's details
	_ = databases.Session.Query(`
		SELECT id, created_at, updated_at, deleted_at FROM messages
		WHERE channel_id = ?
		LIMIT 1`,
		channelID,
	).Scan(&id, &createdAt, &updatedAt, &deletedAt)

	return &Message{
		ID:        id,
		Content:   content,
		AuthorID:  authorID,
		ChannelID: channelID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

// GetMessagesByChannelID retrieves messages for a given channel ID
func GetMessagesByChannelID(channelID string, pageSize, pageNumber int) ([]*Message, error) {
	var messages []*Message

	q := databases.Session.Query(`
		SELECT id, content, author_id, channel_id, created_at, updated_at, deleted_at
		FROM messages
		WHERE channel_id = ? LIMIT ?`,
		channelID, pageSize * pageNumber,
	)

	iter := q.Iter()

	var id, content, authorID, chID string
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time
	var currentIndex int

	for iter.Scan(&id, &content, &authorID, &chID, &createdAt, &updatedAt, &deletedAt) {
		currentIndex++
		if currentIndex <= pageSize*(pageNumber-1) {
			continue
		}
		message := &Message{
			ID:        id,
			Content:   content,
			AuthorID:  authorID,
			ChannelID: chID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		}
		messages = append(messages, message)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}