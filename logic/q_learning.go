package logic

import (
	"math"
	"math/rand/v2"
	"simplegame/constants"
	"simplegame/entities"
)

// layers Width * Height
const (
	states  = 50 * 50
	actions = 4   // 0 = left, 1 = right, 2 = up, 3 = down
	alpha   = 0.1 // learning rate
	gamma   = 0.9 // discount factor
	epsilon = 0.2 // exploration rate
)

type QLearning struct {
	Q [][]float64
}

func (ql *QLearning) Init() {
	ql.Q = make([][]float64, states)
	for i := range ql.Q {
		ql.Q[i] = make([]float64, actions)
	}
	// ql.Q should be all 0's at this point
}

func calcState(enemy_x, enemy_y float64) int {
	tile_e_x := int(math.Round(enemy_x / constants.Tilesize))
	tile_e_y := int(math.Round(enemy_y / constants.Tilesize))
	return tile_e_x * tile_e_y
}

// takes in state and returns movements
func (ql *QLearning) Action(enemy_x, enemy_y float64) entities.SpriteState {
	if rand.Float64() < epsilon {
		return entities.SpriteState(rand.IntN(4))
	}

	// calculate best action from Q matrix
	qVals := ql.Q[calcState(enemy_x, enemy_y)]
	// TODO this fails in initial state of 0
	best := maxIndex(qVals)
	// I'm not using a % here because I want to ensure
	// errors are caught
	return entities.SpriteState(best)
}

func (ql *QLearning) Update(enemy_x, enemy_y float64) (delta_X, delta_Y uint8) {
	// get tile positions, used for states
	// cur_state := calcState(enemy_x, enemy_y)

	// reward := 0

	// new_state := cur_state + alpha*(reward+gamma*argMax(ql.Q[cur_state]))

	return 0, 0
}

func maxIndex(qVals []float64) int {
	bestIndex := 0
	for i := 1; i < actions; i++ {
		if qVals[i] > qVals[bestIndex] {
			bestIndex = i
		}
	}
	return bestIndex
}

func argMax(vals []float64) float64 {
	maxIndex := 0
	maxVal := vals[0]
	for i, v := range vals {
		if v > maxVal {
			maxVal = v
			maxIndex = i
		}
	}
	return float64(maxIndex)
}
