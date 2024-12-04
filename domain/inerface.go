package domain

type Iter interface {
	GetById(int, string) (*RefreshToken, *AccessToken, error)

	Refresh(string, string, string) (*RefreshToken, *AccessToken, error)
	InsertUser(*User) error
}

type Repos interface {
	FetchByToken(string) (*RefreshToken, error)

	UpdateToken(*RefreshToken) error
	FetchIpById(int) (*User, error)
	Insert(*RefreshToken) error
	InsertUser(*User) error
}
