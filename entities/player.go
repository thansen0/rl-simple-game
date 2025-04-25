package entities

import "simplegame/constants"

type Player struct {
	*Sprite
	Projectiles [constants.NumberOfProjectiles]*Projectile
	Health      uint
}

func (p *Player) UpdateAllProjectiles() {
	for i, proj := range p.Projectiles {
		if p.Projectiles[i] != nil && proj.IsAlive {
			p.Projectiles[i].X += p.Projectiles[i].Dx
			p.Projectiles[i].Y += p.Projectiles[i].Dy

			if p.Projectiles[i].X < 0 || p.Projectiles[i].Y < 0 {
				// out of screen area
				p.Projectiles[i].IsAlive = false
			}

			// ignore projectiles that go out of bounds, we'll just write over them
		}
	}
}
