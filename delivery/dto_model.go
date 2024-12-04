package delivery

import (
	"MedodsTestTask/domain"
	"net/http"
	"strconv"
)

type TokenPairDTO struct {
	AccessToken  AccessTokenDTO  `json:"accessToken"`  // DTO Access токена
	RefreshToken RefreshTokenDTO `json:"refreshToken"` // DTO Refresh токена
}

type AccessTokenDTO struct {
	Token string `json:"token"` // Сам JWT токен
}

type RefreshTokenDTO struct {
	Token string `json:"token"` // Сам токен (в base64 формате)
}

type TokenPair struct {
	AccessToken  AccessToken  // Access токен
	RefreshToken RefreshToken // Refresh токен
}

type AccessToken struct {
	Token string // Сам JWT токен
}

type RefreshToken struct {
	Token string // Сам токен (в base64 формате)
}

type User struct {
	ID    int64  `json:"id"`    // ID пользователя
	Email string `json:"email"` // Email пользователя
}

func Convert(param string) int {
	p, err := strconv.Atoi(param)
	if err != nil {
		return -1
	}
	return p
}

func (n *User) DTOUser(r *http.Request) *domain.User {
	return &domain.User{
		ID:       n.ID,
		Email:    n.Email,
		ClientIp: r.RemoteAddr,
	}
}

func DTOTokens(dtoR RefreshTokenDTO, dtoA AccessTokenDTO) TokenPairDTO {
	return TokenPairDTO{AccessToken: dtoA, RefreshToken: dtoR}
}

func DTOA(token *domain.AccessToken) AccessTokenDTO {
	return AccessTokenDTO{Token: token.Token}
}

func DTOR(token *domain.RefreshToken) RefreshTokenDTO {
	return RefreshTokenDTO{Token: token.TokenHash}
}

func (d *TokenPairDTO) Dto() TokenPair {
	return TokenPair{RefreshToken: RefreshToken{Token: d.RefreshToken.Token}, AccessToken: AccessToken{Token: d.AccessToken.Token}}
}
