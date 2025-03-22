package mapper

import (
	"github.com/h1sashin/go-app/db"
	"github.com/h1sashin/go-app/graph/public/model"
)

func MapUser(user *db.User) *model.User {
	return &model.User{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  model.Role(user.Role),
	}
}
