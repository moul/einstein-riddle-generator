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
	Size       int
	Categories int

	Items map[int]ItemList
	Kinds []int
}

type ItemList []string

var BaseInventory map[int]ItemList

const (
	KindNationality = iota
	KindHouseColor
	KindPet
	KindJob
	KindBeverage
	KindWeapon
	KindTransport
	KindRoom
)

var KindName map[int]string

var Kinds = []int{
	KindNationality, KindHouseColor, KindPet, KindJob,
	KindBeverage, KindWeapon, KindTransport, KindRoom,
}

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
		KindWeapon:      "weapon",
		KindTransport:   "transport",
		KindRoom:        "room",
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
			"scottish",
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
			"fanta",
		},
		KindWeapon: {
			"sword",
			"gun",
			"bazooka",
			"grenade",
			"bomb",
			"assault-rifle",
			"shotgun",
			"knife",
			"lasergun",
		},
		KindTransport: {
			"bus",
			"train",
			"car",
			"bike",
			"plane",
			"roller",
			"motorbike",
			"hoverboard",
		},
		KindRoom: {
			"kitchen",
			"bedroom",
			"lobby",
			"living-room",
			"veranda",
			"garden",
			"pool",
			"restroom",
			"bathroom",
		},
	}
}

func NewInventory(size, categories int) *Inventory {
	inventory := Inventory{
		Size:       size,
		Categories: categories,
		Items:      make(map[int]ItemList, 0),
		Kinds:      make([]int, categories),
	}

	for i := range Kinds {
		j := rand.Intn(i + 1)
		Kinds[i], Kinds[j] = Kinds[j], Kinds[i]
	}

	for i := 0; i < categories; i++ {
		inventory.Kinds[i] = Kinds[i]
	}

	for _, kind := range inventory.Kinds {
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
	for _, kind := range i.Kinds {
		fmt.Fprintf(w, "%s\t%s\n", KindName[kind], strings.Join(i.Items[kind], "\t"))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	inventory := NewInventory(8, 8)
	inventory.Show()
}
