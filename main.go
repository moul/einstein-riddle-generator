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
	Vector     []int
	Pickeds    []PickedGroup
	GroupSize  int

	Items map[int]ItemList
	Kinds []int
}

type PickedGroup []int

type ItemList []string

type Item struct {
	Kind   int
	Value  string
	Idx    int
	Person int
}

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

func NewInventory(size, categories, groupSize int) *Inventory {
	inventory := Inventory{
		Size:       size,
		Categories: categories,
		Items:      make(map[int]ItemList, 0),
		Kinds:      make([]int, categories),
		GroupSize:  groupSize,
	}

	inventory.Vector = make([]int, inventory.Length())
	inventory.Pickeds = make([]PickedGroup, 0)

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
	for y, kind := range i.Kinds {
		row := ""
		for j := 0; j < i.Size; j++ {
			row += fmt.Sprintf("%d: %s\t", i.Vector[y*i.Size+j], i.Items[kind][j])
		}
		fmt.Fprintf(w, "%s\t%s\n", KindName[kind], row)
	}
	fmt.Fprintln(w, "\n")
}

func (i *Inventory) Length() int {
	return i.Size * i.Categories
}

func (i *Inventory) EOF() bool {
	return (len(i.Pickeds)+1)*i.GroupSize >= i.Length()
}

func (i *Inventory) PickAvailableGroup(maxLevel int) []int {
	// FIXME: pick more groups of the same person
	// FIXME: pick groups of one item on full left or full right
	length := i.Length()
	picked := make([]int, i.GroupSize)

	for j := 0; j < i.GroupSize; j++ {
		max := rand.Intn(length) + 1
		k := 0
		l := 0
		idx := -1
		for k < max {
			for {
				if i.Vector[(k+l)%length] == maxLevel-1 {
					idx = (k + l) % length
					break
				}
				l++
			}
			k++
		}
		picked[j] = idx
		i.Vector[idx]++
	}
	i.Pickeds = append(i.Pickeds, picked)
	return picked
}

func (i *Inventory) At(idx int) Item {
	item := Item{
		Kind:   i.Kinds[idx/i.Size],
		Idx:    idx,
		Person: idx % i.Size,
	}
	item.Value = i.Items[item.Kind][item.Person]
	return item
}
func (i *Item) Name() string {
	return fmt.Sprintf("%s:%s", KindName[i.Kind], i.Value)
}

func (i *Inventory) GroupString(group PickedGroup) string {
	items := []Item{}

	fullNames := []string{}
	for _, idx := range group {
		item := i.At(idx)
		items = append(items, item)
		fullNames = append(fullNames, item.Name())
	}

	sameKind := true
	samePerson := true

	currentKind := items[0].Kind
	currentPerson := items[0].Person
	for j := 1; j < i.GroupSize; j++ {
		if currentKind != items[j].Kind {
			sameKind = false
		}
		if currentPerson != items[j].Person {
			samePerson = false
		}
	}

	if samePerson {
		return strings.Join(fullNames, " == ")
	}

	switch len(group) {
	case 1:
		if currentPerson == 0 {
			return fmt.Sprintf("%s is on the left", items[0].Name())
		}
		if currentPerson == i.Size {
			return fmt.Sprintf("%s is on the right", items[0].Name())
		}
		if i.Size%2 == 1 && currentPerson == (i.Size-1)/2 {
			return fmt.Sprintf("%s is in the middle", items[0].Name())
		}
		panic("not implemented")
	case 2:
		if items[0].Person == items[1].Person-1 {
			return fmt.Sprintf("%s is direct on the left of %s", items[0].Name(), items[1].Name())
		}
		if items[0].Person == items[1].Person+1 {
			return fmt.Sprintf("%s is direct on the right of %s", items[0].Name(), items[1].Name())
		}
		if items[0].Person < items[1].Person {
			return fmt.Sprintf("%s is on the left of %s", items[0].Name(), items[1].Name())
		}
		if items[0].Person > items[1].Person {
			return fmt.Sprintf("%s is on the left of %s", items[0].Name(), items[1].Name())
		}
	default:
		panic("not implemented")
	}

	return fmt.Sprintf("%v %v %v", sameKind, samePerson, items)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	inventory := NewInventory(5, 5, 2)

	inventory.Show()
	// pick at least each item one time
	for !inventory.EOF() {
		inventory.PickAvailableGroup(1)
	}

	// pick again some items
	for i := 0; i < 3; i++ {
		inventory.PickAvailableGroup(2)
	}

	inventory.Show()

	for _, group := range inventory.Pickeds {
		fmt.Printf("- %s\n", inventory.GroupString(group))
	}
}
