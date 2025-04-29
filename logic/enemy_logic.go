package logic

import (
	"simplegame/entities"
	"simplegame/tilemap"
)

type Logic interface {
	Init()
	Action(tm *tilemap.TilemapJSON, sprite *entities.Enemy, player_x, player_y float64) (X, Y, delta_X, delta_Y float64)
	Update()
}
