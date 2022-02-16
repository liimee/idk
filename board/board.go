package board

import "github.com/liimee/idk/user"

type Sq struct {
	Name  string
	Price int
	Rent  []int
	Set   int
}

func (s Sq) nop() Sq {
	s.Price = 0
	s.Set = 0
	s.Rent = []int{0, 0, 0, 0, 0, 0}
	return s
}

type Ca struct {
	Str string
	Fun func(user.User) user.User
}

var Chance = Sq{Name: "Chance"}.nop()
var CommunityChest = Sq{Name: "Community Chest"}.nop()

var Board = []Sq{
	Sq{Name: "GO"}.nop(),
	{Name: "Hmm(os) Road", Price: 60, Rent: []int{2, 10, 30, 90, 160, 250}, Set: 1},
	CommunityChest,
	{Name: "Well yes but actually no Road", Price: 60, Rent: []int{4, 20, 60, 180, 360, 450}, Set: 1},
	Sq{Name: "Tax", Price: 100}.nop(),
	{Name: "This is a Station", Price: 200, Rent: []int{100}, Set: 0},
	{Name: "Not The Angel Islington", Price: 100, Rent: []int{6, 30, 90, 270, 400, 550}, Set: 2},
	Chance,
	{Name: "${roadname} Road", Price: 100, Rent: []int{6, 30, 90, 270, 400, 550}, Set: 2},
	{Name: "Pen Road", Price: 100, Rent: []int{8, 40, 100, 300, 450, 600}, Set: 2},
	Sq{Name: "IN JAIL"}.nop(),
	{Name: "Small Mall", Price: 140, Rent: []int{10, 50, 150, 450, 625, 750}, Set: 3},
	{Name: "Electric Company", Price: 150, Rent: []int{30}, Set: 0},
	{Name: "Nohall", Price: 140, Rent: []int{10, 50, 150, 450, 625, 750}, Set: 3},
	{Name: "Notaland Avenue", Price: 160, Rent: []int{12, 60, 180, 500, 700, 900}, Set: 3},
	{Name: "${railroad.name} Railroad", Price: 200, Rent: []int{100}, Set: 0},
	{Name: "Electogenius", Price: 180, Rent: []int{14, 70, 200, 550, 750, 950}, Set: 4},
	CommunityChest,
	{Name: "(Ad)", Price: 180, Rent: []int{14, 70, 200, 550, 750, 950}, Set: 4},
	{Name: "Avenue", Price: 200, Rent: []int{16, 80, 220, 600, 800, 1000}, Set: 4},
	Sq{Name: "Free Parking"}.nop(),
	{Name: "(368533454).toString(32)", Price: 220, Rent: []int{18, 90, 250, 700, 875, 1050}, Set: 5},
	Chance,
	{Name: "EGCode", Price: 220, Rent: []int{18, 90, 250, 700, 875, 1050}, Set: 5},
	{Name: "Something Road", Price: 240, Rent: []int{20, 100, 300, 750, 925, 1100}, Set: 5},
	{Name: "BBQ Railroad", Price: 200, Rent: []int{100}, Set: 0},
	{Name: "...", Price: 260, Rent: []int{22, 110, 330, 800, 975, 1150}, Set: 6},
	{Name: "Somewhere", Price: 260, Rent: []int{22, 110, 330, 800, 975, 1150}, Set: 6},
	{Name: "Water Works", Price: 150, Rent: []int{30}, Set: 0},
	{Name: "Please go to jail please p%:", Price: 280, Rent: []int{24, 120, 360, 850, 1025, 1200}, Set: 6},
	Sq{Name: "Go to Jail :)"}.nop(),
	{Name: "Pacific Ocean", Price: 300, Rent: []int{26, 130, 390, 900, 1100, 1275}, Set: 7},
	{Name: "${12086..toString(26)} Avenue", Price: 300, Rent: []int{26, 130, 390, 900, 1100, 1275}, Set: 7},
	CommunityChest,
	{Name: "Not Pacific Ocean But A Bit More Expensive", Price: 320, Rent: []int{28, 150, 450, 1000, 1200, 1400}, Set: 7},
	{Name: "I honestly don't know Station", Price: 200, Rent: []int{100}, Set: 0},
	Chance,
	{Name: "Yes we're almost done!!!!", Price: 350, Rent: []int{35, 175, 500, 1100, 1300, 1500}, Set: 8},
	Sq{Name: "Tax(i)", Price: 200}.nop(),
	{Name: "Undardese Electartica", Price: 400, Rent: []int{50, 200, 600, 1400, 1700, 2000}, Set: 8},
}

var ChanceCards = []Ca{
	{
		Str: "Make General Repairs on All Your Property. For each house pay $25. For each hotel $100.",
		Fun: func(w user.User) user.User {
			for _, v := range w.Ho {
				if v > 0 {
					if v == 5 {
						w.Money -= 100
					} else {
						w.Money -= 25
					}
				}
			}

			return w
		},
	},
	{
		Str: "Speeding fine $15.",
		Fun: func(w user.User) user.User {
			w.Money -= 15

			return w
		},
	},
	{
		Str: "Go back three spaces.",
		Fun: func(w user.User) user.User {
			w.Pos -= 3
			if w.Pos < 0 {
				w.Pos += 40
			}

			return w
		},
	},
	{
		Str: "Bank pays you dividend of $50.",
		Fun: func(w user.User) user.User {
			w.Money += 50

			return w
		},
	},
	{
		Str: "Pay poor tax of $15.",
		Fun: func(w user.User) user.User {
			w.Money -= 15

			return w
		},
	},
	{
		Str: "Take a trip to This Is A Station. If you pass \"GO\" collect $200.",
		Fun: func(w user.User) user.User {
			if w.Pos > 5 {
				w.Money += 200
			}

			w.Pos = 5

			return w
		},
	},
	{
		Str: "ADVANCE to Undardese Electartica.",
		Fun: func(w user.User) user.User {
			w.Pos = 39

			return w
		},
	},
	{
		Str: "ADVANCE to Something Road. If you pass \"GO\" collect $200.",
		Fun: func(w user.User) user.User {
			if w.Pos > 24 {
				w.Money += 200
			}

			w.Pos = 24

			return w
		},
	},
	{
		Str: "Your building loan matures. Collect $150.",
		Fun: func(w user.User) user.User {
			w.Money += 150

			return w
		},
	},
	{
		Str: "ADVANCE to Small Mall. If you pass \"GO\" collect $200.",
		Fun: func(w user.User) user.User {
			if w.Pos > 11 {
				w.Money += 200
			}

			w.Pos = 11

			return w
		},
	},
	{
		Str: "Go to Jail. Go Directly to Jail. Do not pass \"GO\". Do not collect $200.",
		Fun: func(w user.User) user.User {
			w.Pos = 10
			w.InJail = true

			return w
		},
	},
}
