package main

import (
	"schemadsl/backend"
	"schemadsl/frontend"
)

func main() {
	inventory := frontend.LoadFile("inventory.authz")
	dispatcher := frontend.LoadFile("dispatcher.authz")

	backend.WriteOutput("output.zed", []*frontend.Service{inventory, dispatcher})
}
