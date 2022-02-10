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

var Board = []Sq{
	Sq{Name: "GO"}.nop(),
	Sq{Name: "Hmm(os) Road", Price: 60, Rent: []int{2, 10, 30, 90, 160, 250}, Set: 1}.def(),
	Sq{Name: "Community Chest"}.nop(),
	Sq{Name: "Well yes but actually no Road", Price: 60, Rent: []int{4, 20, 60, 180, 360, 450}, Set: 1}.def(),
	Sq{Name: "Tax"}.nop(),
	Sq{Name: "This is a Station", Price: 200, Rent: []int{100}, Set: 0}.def(),
	Sq{Name: "Not The Angel Islington", Price: 100, Rent: []int{6, 30, 90, 270, 400, 550}, Set: 2}.def(),
	Sq{Name: "Chance"}.nop(),
	Sq{Name: "${roadname} Road", Price: 100, Rent: []int{6, 30, 90, 270, 400, 550}, Set: 2}.def(),
	Sq{Name: "Pen Road", Price: 100, Rent: []int{8, 40, 100, 300, 450, 600}, Set: 2}.def(),
	Sq{Name: "IN JAIL"}.nop(),
	Sq{Name: "Small Mall", Price: 140, Rent: []int{10, 50, 150, 450, 625, 750}, Set: 3}.def(),
	Sq{Name: "Electric Company", Price: 150, Rent: []int{30}, Set: 0}.def(),
	Sq{Name: "Nohall", Price: 140, Rent: []int{10, 50, 150, 450, 625, 750}, Set: 3}.def(),
	Sq{Name: "Notaland Avenue", Price: 160, Rent: []int{12, 60, 180, 500, 700, 900}, Set: 3}.def(),
	Sq{Name: "${railroad.name} Railroad", Price: 200, Rent: []int{100}, Set: 0}.def(),
	Sq{Name: "Electogenius", Price: 180, Rent: []int{14, 70, 200, 550, 750, 950}, Set: 4}.def(),
}
