package mock

import (
	mag "github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/token"
)

type Magic struct {
	IsValidtokenFn func(string) (*token.Token, error)
	GetIssuerFn    func(*token.Token) (*mag.UserInfo, error)
}

func (m *Magic) IsValidToken(token string) (*token.Token, error) {
	return m.IsValidtokenFn(token)
}

func (m *Magic) GetIssuer(tok *token.Token) (*mag.UserInfo, error) {
	return m.GetIssuerFn(tok)
}
