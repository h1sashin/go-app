package resolver

import (
	"github.com/h1sashin/go-app/config"
	"github.com/h1sashin/go-app/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Cfg *config.Config
	*service.Service
}
