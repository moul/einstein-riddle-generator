package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/docker/machine/log"
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
	return len(i.Missings()) <= 2
}

func (i *Inventory) PickSamePersonGroup(maxLevel int) []int {
	picked := []int{}

	// arbitrary choosen value to avoid infinite loops
	for try := 0; try < 42; try++ {
		person := rand.Intn(i.Size)

		kinds := make([]int, i.Categories)
		for j := 0; j < i.Categories; j++ {
			kinds[j] = j
		}
		for j := 0; j < i.Categories; j++ {
			k := rand.Intn(i.Categories)
			kinds[j], kinds[k] = kinds[k], kinds[j]
		}
		for _, kind := range kinds {
			idx := kind*i.Size + person
			if i.Vector[idx] <= maxLevel {
				picked = append(picked, idx)
				if len(picked) == i.GroupSize {
					i.Pickeds = append(i.Pickeds, picked)
					for j := 0; j < i.GroupSize; j++ {
						i.Vector[picked[j]]++
					}
					return picked
				}
			}
		}
	}
	panic("should never happen")
	return picked
}

func (i *Inventory) PickAvailableGroup(maxLevel int) []int {
	// FIXME: pick more groups of the same person
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

	if len(group) == 1 {
		currentPerson := items[0].Person
		if currentPerson == 0 {
			return fmt.Sprintf("%s is on the far left", items[0].Name())
		}
		if currentPerson == i.Size-1 {
			return fmt.Sprintf("%s is on the far right", items[0].Name())
		}
		if i.Size%2 == 1 && currentPerson == (i.Size-1)/2 {
			return fmt.Sprintf("%s is in the middle", items[0].Name())
		}
		panic("not implemented")
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
		panic("should never happen")
	case 2:
		// FIXME: sometimes just say "item0 is side by side with item1"
		directLeft := items[0].Person == items[1].Person-1
		directRight := items[0].Person == items[1].Person+1
		if (directLeft || directRight) && rand.Intn(2) > 0 {
			return fmt.Sprintf("%s is on the side of %s", items[0].Name(), items[1].Name())
		}
		if directLeft {
			return fmt.Sprintf("%s is direct on the left of %s", items[0].Name(), items[1].Name())
		}
		if directRight {
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

func (i *Inventory) Missings() []Item {
	missings := []Item{}
	for idx := 0; idx < i.Length(); idx++ {
		if i.Vector[idx] == 0 {
			missings = append(missings, i.At(idx))
		}
	}
	return missings
}

func (i *Inventory) PickItemAtExtremity(targetCounter int) []int {
	picked := make([]int, 1)
	idx := -1

	// limits are arbitrary choosen to avoid infinite loop
	// best effort for now
	for ; targetCounter < 10; targetCounter++ {
		for try := 0; try < 100; try++ {
			direction := rand.Intn(3)
			kind := rand.Intn(i.Categories)
			switch direction {
			case 0:
				idx = kind * i.Size
			case 1:
				idx = kind*i.Size + i.Size - 1
			case 2:
				if i.Size%2 == 1 {
					idx = kind*i.Size + (i.Size-1)/2
				} else {
					continue
				}
			default:
				panic("should never happen")
			}
			if i.Vector[idx] == targetCounter {
				picked[0] = idx
				i.Vector[idx]++
				i.Pickeds = append(i.Pickeds, picked)
				return picked
			}
		}
	}
	panic("should never happen")
	return picked
}

func main() {
	rand.Seed(time.Now().UnixNano())
	inventory := NewInventory(5, 5, 2)

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
			log.Errorf("Invalid riddle: multiple missings item are from the same kind")
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
