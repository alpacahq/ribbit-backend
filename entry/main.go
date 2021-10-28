package main

import (
	"fmt"

	alpaca "github.com/alpacahq/ribbit-backend"
)

func main() {
	alpaca.New().
		WithRoutes(&MyServices{}).
		Run()
}

// MyServices implements github.com/alpacahq/ribbit-backend/route.ServicesI
type MyServices struct{}

// SetupRoutes is our implementation of custom routes
func (s *MyServices) SetupRoutes() {
	fmt.Println("set up our custom routes!")
}
