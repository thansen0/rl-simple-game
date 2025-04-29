package logic

import (
	"simplegame/entities"
	"simplegame/tilemap"
)

type SimpleFollow struct{}

// Nothing to update
func (f *SimpleFollow) Init() {}

func (f *SimpleFollow) Action(tm *tilemap.TilemapJSON, sprite *entities.Enemy, player_x, player_y float64) (X, Y, delta_X, delta_Y float64) {
	prev_player_X := sprite.X
	prev_player_Y := sprite.Y

	if prev_player_X < player_x {
		prev_player_X = tm.GetValidXPos(prev_player_X, 1)
	} else if prev_player_X > player_x {
		prev_player_X = tm.GetValidXPos(prev_player_X, -1)
	}
	if prev_player_Y < player_y {
		prev_player_Y = tm.GetValidYPos(prev_player_Y, 1)
	} else if prev_player_Y > player_y {
		prev_player_Y = tm.GetValidYPos(prev_player_Y, -1)
	}

	dx := sprite.X - prev_player_X
	dy := sprite.Y - prev_player_Y

	activeAnim := sprite.ActiveAnimation(int(dx), int(dy))
	if activeAnim == nil {
		// force an up animation
		activeAnim = sprite.ActiveAnimation(0, 2)
	}
	activeAnim.Update()

	return prev_player_X, prev_player_Y, dx, dy
}

// Update really doesn't matter here, since there is no
// matrix or model to update
func (f *SimpleFollow) Update() {}
