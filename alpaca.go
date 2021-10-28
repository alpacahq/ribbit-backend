package alpaca

import (
	"github.com/alpacahq/ribbit-backend/cmd"
	"github.com/alpacahq/ribbit-backend/route"
)

// New creates a new Alpaca instance
func New() *Alpaca {
	return &Alpaca{}
}

// Alpaca allows us to specify customizations, such as custom route services
type Alpaca struct {
	RouteServices []route.ServicesI
}

// WithRoutes is the builder method for us to add in custom route services
func (g *Alpaca) WithRoutes(RouteServices ...route.ServicesI) *Alpaca {
	return &Alpaca{RouteServices}
}

// Run executes our alpaca functions or servers
func (g *Alpaca) Run() {
	cmd.Execute(g.RouteServices)
}
