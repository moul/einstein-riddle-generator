package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/moul/einstein-riddle-generator"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	inventory := einsteinriddle.NewInventory(5, 5, 2)

	// inventory.Show()
	for i := 0; i < inventory.Length()/3; i++ {
		inventory.PickSamePersonGroup(0)
	}

	secretQuantity := 2

	// pick at least each item one time
	for len(inventory.Missings()) > secretQuantity+1 {
		inventory.PickAvailableGroup(1)
	}

	for len(inventory.Missings()) > secretQuantity {
		inventory.PickItemAtExtremity(0)
	}

	// pick again some items
	for i := 0; i < 3; i++ {
		inventory.PickAvailableGroup(2)
	}

	// pick groups of 1 item on an extremity
	for i := 0; i < 3; i++ {
		inventory.PickItemAtExtremity(1)
	}

	inventory.Show()

	missingsKind := make(map[int]bool, 0)
	for _, missing := range inventory.Missings() {
		if missingsKind[missing.Kind] {
			fmt.Errorf("Invalid riddle: multiple missings item are from the same kind")
			return
		}
		missingsKind[missing.Kind] = true
	}

	for _, group := range inventory.Pickeds {
		fmt.Printf("- %s\n", inventory.GroupString(group))
	}

	fmt.Println("")

	for _, item := range inventory.Missings() {
		fmt.Printf("- where is %s ?\n", item.Name())
	}
}
