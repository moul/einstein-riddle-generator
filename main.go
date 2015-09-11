package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type Inventory struct {
	Items map[int]ItemList
	Size  int
}

type ItemList []string

var BaseInventory map[int]ItemList

const (
	KindNationality = iota
	KindHouseColor
	KindPet
	KindJob
	KindBeverage
)

var KindName map[int]string

var Kinds = []int{KindNationality, KindHouseColor, KindPet, KindJob, KindBeverage}

func (il *ItemList) Shuffle() {
	for i := range *il {
		j := rand.Intn(i + 1)
		(*il)[i], (*il)[j] = (*il)[j], (*il)[i]
	}
}

func init() {
	KindName = map[int]string{
		KindNationality: "nationality",
		KindHouseColor:  "house-color",
		KindPet:         "pet",
		KindJob:         "job",
		KindBeverage:    "beverage",
	}
	BaseInventory = map[int]ItemList{
		KindNationality: {
			"french",
			"english",
			"norvegian",
			"american",
			"portuguese",
			"spannish",
			"german",
		},
		KindHouseColor: {
			"red",
			"yellow",
			"pink",
			"orange",
			"purple",
			"magenta",
			"blue",
			"green",
		},
		KindPet: {
			"dog",
			"cat",
			"horse",
			"poney",
			"fish",
			"whale",
			"beaver",
			"bird",
			"shark",
			"snake",
		},
		KindJob: {
			"teacher",
			"architect",
			"nurse",
			"scientist",
			"student",
			"cop",
			"designer",
			"docker",
		},
		KindBeverage: {
			"beer",
			"wine",
			"water",
			"long-island-ice-tea",
			"coca-cola",
			"dr-pepper",
			"blue-lagoon",
		},
	}
}

func NewInventory(size int) *Inventory {
	inventory := Inventory{
		Size:  size,
		Items: make(map[int]ItemList, 0),
	}

	for _, kind := range Kinds {
		itemList := BaseInventory[kind]
		itemList.Shuffle()
		inventory.Items[kind] = ItemList{}
		for i := 0; i < inventory.Size; i++ {
			inventory.Items[kind] = append(inventory.Items[kind], itemList[i])
		}
	}

	return &inventory
}

func (i *Inventory) Show() {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	header := ""
	for j := 0; j < i.Size; j++ {
		header += fmt.Sprintf("\t%d", j+1)
	}
	fmt.Fprintf(w, "%s\n", header)
	for _, kind := range Kinds {
		fmt.Fprintf(w, "%s\t%s\n", KindName[kind], strings.Join(i.Items[kind], "\t"))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	inventory := NewInventory(5)
	inventory.Show()
}
