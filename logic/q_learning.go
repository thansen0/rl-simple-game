package logic

const (
	states    = 5
	actions   = 4   // 0 = left, 1 = right, 2 = up, 3 = down
	alpha     = 0.1 // learning rate
	gamma     = 0.9 // discount factor
	epsilon   = 0.2 // exploration rate
	episodes  = 1000
	goalState = 4
)

type QLearning struct {
	// *Logic
	Q [][]float64
}

func (ql *QLearning) Init() {
	ql.Q = make([][]float64, states)
	for i := range ql.Q {
		ql.Q[i] = make([]float64, actions)
	}
}

// takes in state and returns movements
func (ql *QLearning) Action() (delta_X, delta_Y float64) {
	return 0, 0
}

func (ql *QLearning) Update() {
}
