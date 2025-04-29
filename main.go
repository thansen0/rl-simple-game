package main

import (
	"fmt"
	"log"

	"simplegame/animations"
	"simplegame/constants"
	"simplegame/entities"
	"simplegame/logic"
	"simplegame/spritesheet"
	"simplegame/tilemap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	// ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// load the image from file
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/Noble/SpriteSheet.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	// load the image from file
	yellowBatImg, _, err := ebitenutil.NewImageFromFile("assets/images/YellowsBat/SpriteSheet.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	// load the image from file
	BlueBatImg, _, err := ebitenutil.NewImageFromFile("assets/images/BlueBat/SpriteSheet.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	// load the image from file
	ButterflyImg, _, err := ebitenutil.NewImageFromFile("assets/images/Butterfly/SpriteSheet.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	// load the image from file
	ButterflyBlueImg, _, err := ebitenutil.NewImageFromFile("assets/images/ButterflyBlue/SpriteSheet.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	// load the image from file
	projectileImg, _, err := ebitenutil.NewImageFromFile("assets/images/Projectile/Sprite.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/maps/TileSet/TilesetFloor.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err := tilemap.NewTilemapJSON("assets/maps/TileSet/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	// maps for performing animation
	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)
	EnemySpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img:         playerImg,
				SpriteSheet: playerSpriteSheet,
				X:           50.0,
				Y:           50.0,
				Animations: map[entities.SpriteState]*animations.Animation{
					entities.Up:    animations.NewAnimation(5, 13, 4, 14.0),
					entities.Down:  animations.NewAnimation(4, 12, 4, 14.0),
					entities.Left:  animations.NewAnimation(6, 14, 4, 14.0),
					entities.Right: animations.NewAnimation(7, 15, 4, 14.0),
				},
			},
			Health:      5,
			Projectiles: [constants.NumberOfProjectiles]*entities.Projectile{},
		},
		enemyLogic: &logic.SimpleFollow{},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img:         yellowBatImg,
					SpriteSheet: EnemySpriteSheet,
					X:           100,
					Y:           100.0,
					Animations: map[entities.SpriteState]*animations.Animation{
						entities.Up:    animations.NewAnimation(5, 13, 4, 8.0),
						entities.Down:  animations.NewAnimation(4, 12, 4, 8.0),
						entities.Left:  animations.NewAnimation(6, 14, 4, 8.0),
						entities.Right: animations.NewAnimation(7, 15, 4, 8.0),
					},
				},
				IsAlive:       true,
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img:         BlueBatImg,
					SpriteSheet: EnemySpriteSheet,
					X:           100.0,
					Y:           100.0,
					Animations: map[entities.SpriteState]*animations.Animation{
						entities.Up:    animations.NewAnimation(5, 13, 4, 9.0),
						entities.Down:  animations.NewAnimation(4, 12, 4, 9.0),
						entities.Left:  animations.NewAnimation(6, 14, 4, 9.0),
						entities.Right: animations.NewAnimation(7, 15, 4, 9.0),
					},
				},
				IsAlive:       true,
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img:         ButterflyImg,
					SpriteSheet: EnemySpriteSheet,
					X:           100.0,
					Y:           100.0,
					Animations: map[entities.SpriteState]*animations.Animation{
						entities.Up:    animations.NewAnimation(5, 13, 4, 7.0),
						entities.Down:  animations.NewAnimation(4, 12, 4, 7.0),
						entities.Left:  animations.NewAnimation(6, 14, 4, 7.0),
						entities.Right: animations.NewAnimation(7, 15, 4, 7.0),
					},
				},
				IsAlive:       true,
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img:         ButterflyBlueImg,
					SpriteSheet: EnemySpriteSheet,
					X:           100.0,
					Y:           100.0,
					Animations: map[entities.SpriteState]*animations.Animation{
						entities.Up:    animations.NewAnimation(5, 13, 4, 11.0),
						entities.Down:  animations.NewAnimation(4, 12, 4, 11.0),
						entities.Left:  animations.NewAnimation(6, 14, 4, 11.0),
						entities.Right: animations.NewAnimation(7, 15, 4, 11.0),
					},
				},
				IsAlive:       true,
				FollowsPlayer: true,
			},
		},
		projectileImg: projectileImg,
		tilemapJSON:   tilemapJSON,
		tilemapImg:    tilemapImg,
		cam:           NewCamera(0.0, 0.0),
		gameStats: GameStats{
			AliveEnemies:    4,
			DeadEnemies:     0,
			ProjectilesShot: 0,
		},
	}

	// set random initialization points for enemies
	for i := range game.enemies {
		game.enemies[i].X, game.enemies[i].Y = game.tilemapJSON.GenValidPos()
	}

	var x int = game.tilemapJSON.TileHeight
	fmt.Println(x)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
