package entities

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	IsAlive       bool
}
