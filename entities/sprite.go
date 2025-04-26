package entities

import (
	"simplegame/animations"
	"simplegame/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteState uint8

const (
	Down SpriteState = iota
	Up
	Left
	Right
)

type Sprite struct {
	Img          *ebiten.Image
	SpriteSheet  *spritesheet.SpriteSheet
	X, Y, Dx, Dy float64
	Animations   map[SpriteState]*animations.Animation
}

func (p *Sprite) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return p.Animations[Right]
	}
	if dx < 0 {
		return p.Animations[Left]
	}
	if dy > 0 {
		return p.Animations[Down]
	}
	if dy < 0 {
		return p.Animations[Up]
	}
	return nil
}
