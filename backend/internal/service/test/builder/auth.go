package builder

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
)

type LogInCredentialsBuilder struct {
	obj ports.LogInCredentials
}

func NewLogInCredentialsBuilder() *LogInCredentialsBuilder {
	return new(LogInCredentialsBuilder)
}

func (b *LogInCredentialsBuilder) Build() ports.LogInCredentials {
	return b.obj
}

func (b *LogInCredentialsBuilder) Default() *LogInCredentialsBuilder {
	b.obj = ports.LogInCredentials{
		Name:     "user",
		Password: "password",
	}
	return b
}

func (b *LogInCredentialsBuilder) SetName(name string) *LogInCredentialsBuilder {
	b.obj.Name = name
	return b
}

func (b *LogInCredentialsBuilder) SetPassword(password string) *LogInCredentialsBuilder {
	b.obj.Password = password
	return b
}

type LogInCredentialsMother struct {
	name     string
	password string
	b        *LogInCredentialsBuilder
}

func NewLoginCredentialsMother(name, password string) *LogInCredentialsMother {
	return &LogInCredentialsMother{
		name:     name,
		password: password,
		b:        NewLogInCredentialsBuilder(),
	}
}

func (m *LogInCredentialsMother) Create() ports.LogInCredentials {
	return m.b.Default().SetName(m.name).SetPassword(m.password).Build()
}

type PayloadBuilder struct {
	obj domain.Payload
}

func NewPayloadBuilder() *PayloadBuilder {
	return new(PayloadBuilder)
}

func (b *PayloadBuilder) Build() domain.Payload {
	return b.obj
}

func (b *PayloadBuilder) Default() *PayloadBuilder {
	b.obj = domain.Payload{
		UserID: uuid.New(),
		Role:   domain.UserRole,
	}
	return b
}

func (b *PayloadBuilder) SetID(id uuid.UUID) *PayloadBuilder {
	b.obj.UserID = id
	return b
}

func (b *PayloadBuilder) SetRole(role int) *PayloadBuilder {
	b.obj.Role = role
	return b
}

type TokenPairBuilder struct {
	obj domain.TokenPair
}

func NewTokenPairBuilder() *TokenPairBuilder {
	return new(TokenPairBuilder)
}

func (b *TokenPairBuilder) Build() domain.TokenPair {
	return b.obj
}

func (b *TokenPairBuilder) Default() *TokenPairBuilder {
	b.obj = domain.TokenPair{
		AccessToken:  "accesstoken",
		RefreshToken: "refreshtoken",
	}
	return b
}

func (b *TokenPairBuilder) SetAccessToken(token string) *TokenPairBuilder {
	b.obj.AccessToken = token
	return b
}

func (b *TokenPairBuilder) SetRefreshToken(token string) *TokenPairBuilder {
	b.obj.RefreshToken = token
	return b
}
