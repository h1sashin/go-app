package mapper

import (
	"github.com/h1sashin/go-app/graph/public/model"
	"github.com/h1sashin/go-app/service"
)

func MapTokens(tokens *service.Tokens) *model.Tokens {
	return &model.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
