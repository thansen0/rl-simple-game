package logic

import (
	"simplegame/entities"
	"simplegame/tilemap"
)

type SimpleFollow struct{}

// Nothing to update
func (f *SimpleFollow) Init() {}

func (f *SimpleFollow) Action(tm *tilemap.TilemapJSON, cur_sprite *entities.Enemy, player_x, player_y float64) (en *entities.Enemy) {
	prev_player_X := cur_sprite.X
	prev_player_Y := cur_sprite.Y

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

	dx := cur_sprite.X - prev_player_X
	dy := cur_sprite.Y - prev_player_Y

	activeAnim := cur_sprite.ActiveAnimation(int(dx), int(dy))
	if activeAnim == nil {
		// force an up animation
		activeAnim = cur_sprite.ActiveAnimation(0, 2)
	}
	activeAnim.Update()

	enn := &entities.Enemy{
		Sprite: &entities.Sprite{
			Img:         cur_sprite.Img,
			SpriteSheet: cur_sprite.SpriteSheet,
			X:           prev_player_X,
			Y:           prev_player_Y,
			Animations:  cur_sprite.Animations,
		},
		IsAlive:       cur_sprite.IsAlive,
		FollowsPlayer: true,
	}

	return enn
}

// Update really doesn't matter here, since there is no
// matrix or model to update
func (f *SimpleFollow) Update() {}
