package entities

type Projectile struct {
	*Sprite
	Damage  uint
	IsAlive bool
}
