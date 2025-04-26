package entities

import "simplegame/animations"

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	IsAlive       bool
}

func (en *Enemy) CreateNewEnemy(g_x, g_y float64) *Enemy {
	return &Enemy{
		Sprite: &Sprite{
			Img:         en.Img,
			SpriteSheet: en.SpriteSheet,
			X:           g_x,
			Y:           g_y,
			Animations: map[SpriteState]*animations.Animation{
				Up:    animations.NewAnimation(5, 13, 4, 8.0),
				Down:  animations.NewAnimation(4, 12, 4, 8.0),
				Left:  animations.NewAnimation(6, 14, 4, 8.0),
				Right: animations.NewAnimation(7, 15, 4, 8.0),
			},
		},
		IsAlive:       true,
		FollowsPlayer: true,
	}
}
