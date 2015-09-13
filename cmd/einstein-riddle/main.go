package main

import (
	"fmt"

	"github.com/moul/einstein-riddle-generator"
)

func main() {
	options := einsteinriddle.Options{
		Size:       5,
		Categories: 5,
		Secrets:    2,
	}
	generator := einsteinriddle.NewGenerator(options)

	// Shazam
	generator.Shazam()

	// Print map
	generator.Show()

	// Print riddle
	for _, group := range generator.Pickeds {
		fmt.Printf("- %s\n", generator.GroupString(group))
	}
	fmt.Println("")
	for _, item := range generator.Missings() {
		fmt.Printf("- where is %s ?\n", item.Name())
	}
}
