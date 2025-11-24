package models

import (
	"MessagesService/databases"
)

// Message represents a message in the system
type Message struct {
	ID      	  	string
	Content 	  	string
	AuthorID 		string
	ChannelID 		string
	CreatedAt 		string
	UpdatedAt 		string
	DeletedAt 		string
}

// ToMap converts a Message to a map[string]any
func (m *Message) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"author_id":  m.AuthorID,
		"channel_id": m.ChannelID,
		"content":    m.Content,
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
		"deleted_at": m.DeletedAt,
	}
}

// NewMessage creates a new Message and saves it to the database
func NewMessage(authorID, channelID, content string) *Message {
	var id, createdAt, updatedAt, deletedAt string


	if err := databases.Session.Query(`
		INSERT INTO messages (id, content, author_id, channel_id, created_at, updated_at, deleted_at)
		VALUES (uuid(), ?, ?, ?, toTimestamp(now()), toTimestamp(now()), null)`,
		content, authorID, channelID,
	).Exec(); err != nil {
		return nil
	}

	_ = databases.Session.Query(`
		SELECT id FROM messages
		WHERE content = ? AND author_id = ? AND channel_id = ? AND created_at = toTimestamp(now())
		LIMIT 1`,
		content, authorID, channelID,
	).Scan(&id, &createdAt, &updatedAt, &deletedAt)

	return &Message{
		ID:       	id,
		Content:   	content,
		AuthorID: 	authorID,
		ChannelID:	channelID,
		CreatedAt:	createdAt,
		UpdatedAt:	updatedAt,
		DeletedAt:	deletedAt,
	}
}