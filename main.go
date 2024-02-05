package main

import (
	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal"
)

func main() {
	drone.New(internal.New).Boot()
}
