package magic

import (
	mag "github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/token"
)

// Service is the interface to our magic service
type Service interface {
	IsValidToken(string) (*token.Token, error)
	GetIssuer(*token.Token) (*mag.UserInfo, error)
}
