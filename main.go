package main

import (
	"log"

	"simplegame/animations"
	"simplegame/entities"
	"simplegame/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
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

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/maps/TileSet/TilesetFloor.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/TileSet/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)
	yellowBatSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
				Animations: map[entities.SpriteState]*animations.Animation{
					entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
					entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
					entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
					entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
				},
			},
			Health: 5,
		},
		playerSpriteSheet: playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: yellowBatImg,
					X:   100.0,
					Y:   100.0,
					Animations: map[entities.SpriteState]*animations.Animation{
						entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
						entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
						entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
						entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
					},
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: yellowBatImg,
					X:   150.0,
					Y:   50.0,
					Animations: map[entities.SpriteState]*animations.Animation{
						entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
						entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
						entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
						entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
					},
				},
				FollowsPlayer: true,
			},
		},
		yellowBatSpriteSheet: yellowBatSpriteSheet,
		tilemapJSON:          tilemapJSON,
		tilemapImg:           tilemapImg,
		cam:                  NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
