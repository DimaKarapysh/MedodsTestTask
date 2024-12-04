package iteractors

import (
	"MedodsTestTask/app/core"
	"MedodsTestTask/domain"
	"MedodsTestTask/tools"
	"context"
	"github.com/pkg/errors"
	"time"
)

type Iter struct {
	log   core.Logger
	repos domain.Repos
	ctx   context.Context
}

func NewIter(log core.Logger, repos domain.Repos, ctx context.Context) *Iter {
	return &Iter{
		log:   log,
		repos: repos,
		ctx:   ctx,
	}
}

func (i *Iter) GetById(id int, ip string) (*domain.RefreshToken, *domain.AccessToken, error) {

	client, err := i.repos.FetchIpById(id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch IP by ID")
	}
	if client == nil {
		return nil, nil, errors.New("user not found")
	}

	if !tools.CheckClientIpHash(ip, client.ClientIp) {
		return nil, nil, errors.New("client IP hash mismatch")
	}

	accessToken, err := tools.GenerateAccessToken(ip)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate access token")
	}

	refreshTokenHash, err := tools.GenerateRefreshRandHash()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate refresh token")
	}

	accessT := &domain.AccessToken{
		Token: accessToken,
	}

	refreshT := &domain.RefreshToken{
		UserID:    id,
		TokenHash: refreshTokenHash,
		ClientIP:  ip,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1000 * time.Minute),
	}

	if err := i.repos.Insert(refreshT); err != nil {
		return nil, nil, errors.Wrap(err, "failed to insert refresh token")
	}

	return refreshT, accessT, nil
}

func (i *Iter) InsertUser(user *domain.User) error {

	err := tools.HashClientIP(user)
	if err != nil {
		return errors.Wrap(err, " Can not HashClientIP")
	}

	err = i.repos.InsertUser(user)
	if err != nil {
		return errors.Wrap(err, "IterInsertErr")
	}
	return nil
}

func (i *Iter) Refresh(refresh, access, ip string) (*domain.RefreshToken, *domain.AccessToken, error) {

	token, err := i.repos.FetchByToken(refresh)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch refresh token")
	}
	if token == nil {
		return nil, nil, errors.New("refresh token not found")
	}

	ok, clientIp, err := tools.ExtractData(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to extract data")
	}
	if !ok {
		return nil, nil, errors.New("invalid access token")
	}
	if clientIp != ip {

		user, err := i.repos.FetchIpById(token.UserID)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to fetch user ID")
		}

		tools.SendEmailWarning(user.Email, clientIp, ip)
		return nil, nil, errors.New("invalid client IP")
	}

	accessToken, err := tools.GenerateAccessToken(token.ClientIP)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate access token")
	}

	refreshTokenHash, err := tools.GenerateRefreshRandHash()

	acc := &domain.AccessToken{
		Token: accessToken,
	}

	refreshToken := &domain.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: refreshTokenHash,
		ClientIP:  token.ClientIP,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1000 * time.Minute),
	}

	err = i.repos.UpdateToken(refreshToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to update refresh token")
	}

	return refreshToken, acc, nil
}
