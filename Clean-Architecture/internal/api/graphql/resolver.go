package graphql

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/service"
)

type Resolver struct {
	OrderService service.OrderService
}
