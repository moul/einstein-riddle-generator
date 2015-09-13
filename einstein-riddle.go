package einsteinriddle

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type Generator struct {
	Options Options

	Vector  []int
	Pickeds []PickedGroup

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

var BaseGenerator map[int]ItemList

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
	BaseGenerator = map[int]ItemList{
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

type Options struct {
	Size             int
	Categories       int
	GroupSize        int
	Secrets          int
	Seed             int64
	SamePersonGroups int
}

func (o *Options) ApplyDefaults() {
	if o.Size == 0 {
		o.Size = 5
	}
	if o.Categories == 0 {
		o.Categories = 5
	}
	if o.Secrets == 0 {
		o.Secrets = 2
	}
	if o.GroupSize == 0 {
		o.GroupSize = 2
	}
	if o.Seed == 0 {
		o.Seed = time.Now().UnixNano()
	}
	if o.SamePersonGroups == 0 {
		o.SamePersonGroups = o.Size * o.Categories / 3
	}
}

func NewGenerator(options Options) *Generator {
	options.ApplyDefaults()

	rand.Seed(options.Seed)

	generator := Generator{
		Options: options,
		Items:   make(map[int]ItemList, 0),
		Kinds:   make([]int, options.Categories),
	}
	generator.Vector = make([]int, generator.Length())
	generator.Pickeds = make([]PickedGroup, 0)

	for i := range Kinds {
		j := rand.Intn(i + 1)
		Kinds[i], Kinds[j] = Kinds[j], Kinds[i]
	}

	for i := 0; i < options.Categories; i++ {
		generator.Kinds[i] = Kinds[i]
	}

	for _, kind := range generator.Kinds {
		itemList := BaseGenerator[kind]
		itemList.Shuffle()
		generator.Items[kind] = ItemList{}
		for i := 0; i < options.Size; i++ {
			generator.Items[kind] = append(generator.Items[kind], itemList[i])
		}
	}

	return &generator
}

func (g *Generator) Show() {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	header := ""
	for j := 0; j < g.Options.Size; j++ {
		header += fmt.Sprintf("\t%d", j+1)
	}
	fmt.Fprintf(w, "%s\n", header)
	for y, kind := range g.Kinds {
		row := ""
		for j := 0; j < g.Options.Size; j++ {
			row += fmt.Sprintf("%d: %s\t", g.Vector[y*g.Options.Size+j], g.Items[kind][j])
		}
		fmt.Fprintf(w, "%s\t%s\n", KindName[kind], row)
	}
	fmt.Fprintln(w, "\n")
}

func (g *Generator) Length() int {
	return g.Options.Size * g.Options.Categories
}

func (g *Generator) EOF() bool {
	return len(g.Missings()) <= 2
}

func (g *Generator) PickSamePersonGroup(maxLevel int) []int {
	picked := []int{}

	// arbitrary choosen value to avoid infinite loops
	for try := 0; try < 42; try++ {
		person := rand.Intn(g.Options.Size)

		kinds := make([]int, g.Options.Categories)
		for j := 0; j < g.Options.Categories; j++ {
			kinds[j] = j
		}
		for j := 0; j < g.Options.Categories; j++ {
			k := rand.Intn(g.Options.Categories)
			kinds[j], kinds[k] = kinds[k], kinds[j]
		}
		for _, kind := range kinds {
			idx := kind*g.Options.Size + person
			if g.Vector[idx] <= maxLevel {
				picked = append(picked, idx)
				if len(picked) == g.Options.GroupSize {
					g.Pickeds = append(g.Pickeds, picked)
					for j := 0; j < g.Options.GroupSize; j++ {
						g.Vector[picked[j]]++
					}
					return picked
				}
			}
		}
	}
	panic("should never happen")
	return picked
}

func (g *Generator) PickAvailableGroup(maxLevel int) []int {
	// FIXME: pick more groups of the same person
	length := g.Length()
	picked := make([]int, g.Options.GroupSize)

	for j := 0; j < g.Options.GroupSize; j++ {
		max := rand.Intn(length) + 1
		k := 0
		l := 0
		idx := -1
		for k < max {
			for {
				if g.Vector[(k+l)%length] == maxLevel-1 {
					idx = (k + l) % length
					break
				}
				l++
			}
			k++
		}
		picked[j] = idx
		g.Vector[idx]++
	}
	g.Pickeds = append(g.Pickeds, picked)
	return picked
}

func (g *Generator) At(idx int) Item {
	item := Item{
		Kind:   g.Kinds[idx/g.Options.Size],
		Idx:    idx,
		Person: idx % g.Options.Size,
	}
	item.Value = g.Items[item.Kind][item.Person]
	return item
}
func (i *Item) Name() string {
	return fmt.Sprintf("%s:%s", KindName[i.Kind], i.Value)
}

func (g *Generator) GroupString(group PickedGroup) string {
	items := []Item{}

	fullNames := []string{}
	for _, idx := range group {
		item := g.At(idx)
		items = append(items, item)
		fullNames = append(fullNames, item.Name())
	}

	if len(group) == 1 {
		currentPerson := items[0].Person
		if currentPerson == 0 {
			return fmt.Sprintf("%s is on the far left", items[0].Name())
		}
		if currentPerson == g.Options.Size-1 {
			return fmt.Sprintf("%s is on the far right", items[0].Name())
		}
		if g.Options.Size%2 == 1 && currentPerson == (g.Options.Size-1)/2 {
			return fmt.Sprintf("%s is in the middle", items[0].Name())
		}
		panic("not implemented")
	}

	sameKind := true
	samePerson := true

	currentKind := items[0].Kind
	currentPerson := items[0].Person
	for j := 1; j < g.Options.GroupSize; j++ {
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
			return fmt.Sprintf("%s is on the right of %s", items[0].Name(), items[1].Name())
		}
	default:
		panic("not implemented")
	}

	return fmt.Sprintf("%v %v %v", sameKind, samePerson, items)
}

func (g *Generator) Missings() []Item {
	missings := []Item{}
	for idx := 0; idx < g.Length(); idx++ {
		if g.Vector[idx] == 0 {
			missings = append(missings, g.At(idx))
		}
	}
	return missings
}

func (g *Generator) PickItemAtExtremity(targetCounter int) []int {
	picked := make([]int, 1)
	idx := -1

	// limits are arbitrary choosen to avoid infinite loop
	// best effort for now
	for ; targetCounter < 10; targetCounter++ {
		for try := 0; try < 100; try++ {
			direction := rand.Intn(3)
			kind := rand.Intn(g.Options.Categories)
			switch direction {
			case 0:
				idx = kind * g.Options.Size
			case 1:
				idx = kind*g.Options.Size + g.Options.Size - 1
			case 2:
				if g.Options.Size%2 == 1 {
					idx = kind*g.Options.Size + (g.Options.Size-1)/2
				} else {
					continue
				}
			default:
				panic("should never happen")
			}
			if g.Vector[idx] == targetCounter {
				picked[0] = idx
				g.Vector[idx]++
				g.Pickeds = append(g.Pickeds, picked)
				return picked
			}
		}
	}
	panic("should never happen")
	return picked
}

func (g *Generator) Shazam() error {
	for i := 0; i < g.Options.SamePersonGroups; i++ {
		g.PickSamePersonGroup(0)
	}

	// pick at least each item one time
	for len(g.Missings()) > g.Options.Secrets+1 {
		g.PickAvailableGroup(1)
	}

	for len(g.Missings()) > g.Options.Secrets {
		g.PickItemAtExtremity(0)
	}

	// pick again some items
	for i := 0; i < 3; i++ {
		g.PickAvailableGroup(2)
	}

	// pick groups of 1 item on an extremity
	for i := 0; i < 3; i++ {
		g.PickItemAtExtremity(1)
	}

	missingsKind := make(map[int]bool, 0)
	for _, missing := range g.Missings() {
		if missingsKind[missing.Kind] {
			return fmt.Errorf("Invalid riddle: multiple missings item are from the same kind")
		}
		missingsKind[missing.Kind] = true
	}

	return nil
}
