package domain

import "time"

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ClientIP  string    `db:"client_ip"`
	IssuedAt  time.Time `db:"issued_at"`
	ExpiresAt time.Time `db:"expires_at"`

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type AccessToken struct {
	Token string `json:"token"` // Сам JWT токен

}

type User struct {
	ID       int64  `json:"id"`    // ID пользователя
	Email    string `json:"email"` // Email пользователя
	ClientIp string

	CreatedAt time.Time  `json:"created_at"` // Время создания записи
	UpdatedAt time.Time  `json:"updated_at"` // Время последнего обновления записи
	DeletedAt *time.Time `json:"deleted_at"` // Время удаления записи (мягкое удаление), может быть nil
}
