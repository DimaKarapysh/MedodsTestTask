package repository

import (
	"MedodsTestTask/app/core"
	"MedodsTestTask/domain"
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type Repos struct {
	log core.Logger
	db  *sql.DB
	ctx context.Context
}

func NewRepos(log core.Logger, db *sql.DB, ctx context.Context) *Repos {
	return &Repos{
		log: log,
		db:  db,
		ctx: ctx,
	}
}

func (r *Repos) FetchByToken(t string) (*domain.RefreshToken, error) {
	token := &domain.RefreshToken{}

	query := `SELECT id, user_id, token_hash, client_ip, issued_at, expires_at, created_at, updated_at FROM token WHERE token_hash=$1`

	err := r.db.QueryRow(query, t).Scan(&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ClientIP,
		&token.IssuedAt, // Добавлено для корректного сканирования
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt)

	switch err {
	case nil:
		return token, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err

	}

}

func (r *Repos) Insert(token *domain.RefreshToken) error {
	query := `INSERT INTO token (user_id, token_hash, client_ip, issued_at, expires_at, created_at, updated_at)VALUES ($1, $2, $3, $4, $5, $6, $7)ON CONFLICT (id) DO NOTHING;`

	_, err := r.db.Exec(query,
		token.UserID,
		token.TokenHash,
		token.ClientIP,
		token.IssuedAt,
		token.ExpiresAt,
		token.CreatedAt,
		token.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "ReposInsertPostgresErr")
	}
	return nil
}

func (r *Repos) InsertUser(user *domain.User) error {

	query := `INSERT INTO users (id, email,tokenip) VALUES ($1,$2,$3)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.ClientIp)
	if err != nil {
		return errors.Wrap(err, "ReposInsertPostgresErr")
	}

	return err
}

func (r *Repos) FetchIpById(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email,tokenip FROM users WHERE id=$1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Email, &user.ClientIp)

	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err

	}
}

func (r *Repos) UpdateToken(token *domain.RefreshToken) error {
	query := `
    UPDATE token
    SET
        user_id = $2,
        token_hash = $3,
        client_ip = $4,
        issued_at = $5,
        expires_at = $6,
        created_at = $7,
        updated_at = now()
    WHERE id = $1
`

	r.log.Debug("Updating token with ID: %v, UserID: %v, TokenHash: %v, ClientIP: %v\n",
		token.ID, token.UserID, token.TokenHash, token.ClientIP)
	_, err := r.db.Exec(query, token.ID, token.UserID, token.TokenHash, token.ClientIP, token.IssuedAt, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "ReposUpdatePostgresErr")
	}
	return nil
}
