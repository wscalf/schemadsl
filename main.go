package main

import (
	"schemadsl/backend"
	"schemadsl/frontend"
)

func main() {
	service := frontend.LoadFile("inventory.authz")

	backend.WriteOutput("output.zed", []*frontend.Service{service})
}
