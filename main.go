package main

import (
	"log"
	"polygon-pools/client"
	"polygon-pools/pools"
	"polygon-pools/reserves"
)

func main() {
	newPools := pools.NewPools()
	if err := newPools.Start(); err != nil {
		log.Fatal(err)
	}

	newReserves := reserves.NewReserves(newPools)
	if err := newReserves.Start(); err != nil {
		log.Fatal(err)
	}

	newClient := client.NewClient(newReserves)
	if err := newClient.Start(); err != nil {
		log.Fatal(err)
	}
}
