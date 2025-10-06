package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`          // Не сериализуем в JSON
	AvatarURL *string   `json:"avatar_url"` // Указатель для nullable поля
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	Author    *User     `json:"author,omitempty"`
	ParentID  *int      `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Like struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}
