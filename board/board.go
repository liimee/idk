package board

type Sq struct {
	Name  string
	Price int
	Rent  []int
	House int
	Hotel int
	Set   int
}

func (s Sq) nop() Sq {
	s = s.def()
	s.Price = 0
	s.Set = 0
	s.Rent = []int{0, 0, 0, 0, 0, 0}
	return s
}

func (s Sq) def() Sq {
	s.House = 0
	s.Hotel = 0
	return s
}

var Chance = Sq{Name: "Chance"}.nop()
var CommunityChest = Sq{Name: "Community Chest"}.nop()

var Board = []Sq{
	Sq{Name: "GO"}.nop(),
	Sq{Name: "Hmm(os) Road", Price: 60, Rent: []int{2, 10, 30, 90, 160, 250}, Set: 1}.def(),
	CommunityChest,
	Sq{Name: "Well yes but actually no Road", Price: 60, Rent: []int{4, 20, 60, 180, 360, 450}, Set: 1}.def(),
	Sq{Name: "Tax"}.nop(),
	Sq{Name: "This is a Station", Price: 200, Rent: []int{100}, Set: 0}.def(),
	Sq{Name: "Not The Angel Islington", Price: 100, Rent: []int{6, 30, 90, 270, 400, 550}, Set: 2}.def(),
	Chance,
	Sq{Name: "${roadname} Road", Price: 100, Rent: []int{6, 30, 90, 270, 400, 550}, Set: 2}.def(),
	Sq{Name: "Pen Road", Price: 100, Rent: []int{8, 40, 100, 300, 450, 600}, Set: 2}.def(),
	Sq{Name: "IN JAIL"}.nop(),
	Sq{Name: "Small Mall", Price: 140, Rent: []int{10, 50, 150, 450, 625, 750}, Set: 3}.def(),
	Sq{Name: "Electric Company", Price: 150, Rent: []int{30}, Set: 0}.def(),
	Sq{Name: "Nohall", Price: 140, Rent: []int{10, 50, 150, 450, 625, 750}, Set: 3}.def(),
	Sq{Name: "Notaland Avenue", Price: 160, Rent: []int{12, 60, 180, 500, 700, 900}, Set: 3}.def(),
	Sq{Name: "${railroad.name} Railroad", Price: 200, Rent: []int{100}, Set: 0}.def(),
	Sq{Name: "Electogenius", Price: 180, Rent: []int{14, 70, 200, 550, 750, 950}, Set: 4}.def(),
	CommunityChest,
	Sq{Name: "(Ad)", Price: 180, Rent: []int{14, 70, 200, 550, 750, 950}, Set: 4}.def(),
	Sq{Name: "Avenue", Price: 200, Rent: []int{16, 80, 220, 600, 800, 1000}, Set: 4}.def(),
	Sq{Name: "Free Parking"}.nop(),
	Sq{Name: "(368533454).toString(32)", Price: 220, Rent: []int{18, 90, 250, 700, 875, 1050}, Set: 5}.def(),
	Chance,
	Sq{Name: "EGCode", Price: 220, Rent: []int{18, 90, 250, 700, 875, 1050}, Set: 5}.def(),
	Sq{Name: "Something Road", Price: 240, Rent: []int{20, 100, 300, 750, 925, 1100}, Set: 5}.def(),
	Sq{Name: "BBQ Railroad", Price: 200, Rent: []int{100}, Set: 0}.def(),
	Sq{Name: "...", Price: 260, Rent: []int{22, 110, 330, 800, 975, 1150}, Set: 6}.def(),
	Sq{Name: "Somewhere", Price: 260, Rent: []int{22, 110, 330, 800, 975, 1150}, Set: 6}.def(),
	Sq{Name: "Water Works", Price: 150, Rent: []int{30}, Set: 0}.def(),
	Sq{Name: "Please go to jail please 🥺", Price: 280, Rent: []int{24, 120, 360, 850, 1025, 1200}, Set: 6}.def(),
	Sq{Name: "Go to Jail :)"}.nop(),
	Sq{Name: "Pacific Ocean", Price: 300, Rent: []int{26, 130, 390, 900, 1100, 1275}, Set: 7}.def(),
	Sq{Name: "${12086..toString(26)} Avenue", Price: 300, Rent: []int{26, 130, 390, 900, 1100, 1275}, Set: 7}.def(),
	CommunityChest,
	Sq{Name: "Not Pacific Ocean But A Bit More Expensive", Price: 320, Rent: []int{28, 150, 450, 1000, 1200, 1400}, Set: 7}.def(),
	Sq{Name: "I honestly don't know Station", Price: 200, Rent: []int{100}, Set: 0}.def(),
	Chance,
	Sq{Name: "Yes we're almost done!!!!", Price: 350, Rent: []int{35, 175, 500, 1100, 1300, 1500}, Set: 8}.def(),
	Sq{Name: "Tax(i)"}.nop(),
	Sq{Name: "Undardese Electartica", Price: 400, Rent: []int{50, 200, 600, 1400, 1700, 2000}, Set: 8}.def(),
}
