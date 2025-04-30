package logic

import (
	"simplegame/entities"
	"simplegame/tilemap"
)

type Logic interface {
	Init()
	Action(tm *tilemap.TilemapJSON, sprite *entities.Enemy, player_x, player_y float64) (en *entities.Enemy)
	Update()
}
