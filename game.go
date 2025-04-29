package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"simplegame/constants"
	"simplegame/entities"
	"simplegame/logic"
	"simplegame/tilemap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameStats struct {
	AliveEnemies    uint32
	DeadEnemies     uint32
	ProjectilesShot uint32
}

type Game struct {
	gameStats     GameStats
	enemyLogic    logic.Logic
	player        *entities.Player
	enemies       []*entities.Enemy
	tilemapJSON   *tilemap.TilemapJSON
	tilemapImg    *ebiten.Image
	projectileImg *ebiten.Image // may not be the best place for this
	cam           *Camera
}

var projectile_counter int = 0

func (g *Game) Update() error {
	var prev_player_X float64
	var prev_player_Y float64
	prev_player_X = g.player.X
	prev_player_Y = g.player.Y

	// we modify projectiles again later one
	// update all projectiles
	g.player.UpdateAllProjectiles()

	// move the player based on keyboar input (left, right, up down)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X = g.tilemapJSON.GetValidXPos(prev_player_X, -2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X = g.tilemapJSON.GetValidXPos(prev_player_X, 2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y = g.tilemapJSON.GetValidYPos(prev_player_Y, -2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y = g.tilemapJSON.GetValidYPos(prev_player_Y, 2)
	}

	// needs to be updated for other functions
	g.player.Dx = prev_player_X - prev_player_X
	g.player.Dy = prev_player_Y - prev_player_Y

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	// add behavior to the enemies
	for _, sprite := range g.enemies {
		// really should be returning a whole enemy struct to
		// replace the existing one in the slice so we can use a
		// channel
		sprite.X, sprite.Y, sprite.Dx, sprite.Dy = g.enemyLogic.Action(g.tilemapJSON, sprite, g.player.X, g.player.Y)

		// sprite is a pointer, so we're updating in place
		// prev_player_X = sprite.X
		// prev_player_Y = sprite.Y

		// if sprite.X < g.player.X {
		// 	sprite.X = g.tilemapJSON.GetValidXPos(sprite.X, 1)
		// } else if sprite.X > g.player.X {
		// 	sprite.X = g.tilemapJSON.GetValidXPos(sprite.X, -1)
		// }
		// if sprite.Y < g.player.Y {
		// 	sprite.Y = g.tilemapJSON.GetValidYPos(sprite.Y, 1)
		// } else if sprite.Y > g.player.Y {
		// 	sprite.Y = g.tilemapJSON.GetValidYPos(sprite.Y, -1)
		// }

		// sprite.Dx = sprite.X - prev_player_X
		// sprite.Dy = sprite.Y - prev_player_Y

		// activeAnim := sprite.ActiveAnimation(int(sprite.Dx), int(sprite.Dy))
		// if activeAnim == nil {
		// 	// force an up animation
		// 	activeAnim = sprite.ActiveAnimation(0, 2)
		// }
		// activeAnim.Update()
	}

	////////////////////////////////////////////////
	// start new thread to check whether any projectiles overlap with an enemy
	// check for interactions with enemies
	for i, proj := range g.player.Projectiles {
		if g.player.Projectiles[i] == nil {
			continue
		}
		for j, en := range g.enemies {
			if g.enemies[j] == nil {
				continue
			}

			if math.Abs((en.X+4)-proj.X) < 4 && math.Abs((en.Y+4)-proj.Y) < 4 {
				// overlap, kill both
				g.enemies[j] = en.CreateNewEnemy(g.tilemapJSON.GenValidPosOutsideCamera(g.player.X, g.player.Y))
				g.player.Projectiles[i].IsAlive = false

				// Create new enemy; each death creates new new ones
				g.enemies = append(g.enemies, en.CreateNewEnemy(g.tilemapJSON.GenValidPosOutsideCamera(g.player.X, g.player.Y)))

				g.gameStats.AliveEnemies += 2
				g.gameStats.DeadEnemies += 1
			}
		}
	}

	// var wg sync.WaitGroup

	// type CollisionResult struct {
	// 	ProjIndex  int
	// 	EnemyIndex int
	// }

	// // Channel to collect results safely
	// collisionCh := make(chan CollisionResult, 10000)

	// // Split work into chunks (one chunk per projectile, or batch them)
	// for i, proj := range g.player.Projectiles {
	// 	if proj == nil {
	// 		continue
	// 	}

	// 	wg.Add(1)
	// 	go func(projIndex int, proj *entities.Projectile) {
	// 		defer wg.Done()

	// 		for j, en := range g.enemies {
	// 			if en == nil {
	// 				continue
	// 			}

	// 			if math.Abs((en.X+4)-proj.X) < 4 && math.Abs((en.Y+4)-proj.Y) < 4 {
	// 				// Report collision (defer the actual modification)
	// 				collisionCh <- CollisionResult{ProjIndex: projIndex, EnemyIndex: j}
	// 				break // Once a projectile hits, it usually disappears
	// 			}
	// 		}
	// 	}(i, proj)
	// }

	// // Close channel when all goroutines are done
	// go func() {
	// 	wg.Wait()
	// 	close(collisionCh)
	// }()

	// // Handle all the enemy collisions sequentially
	// for collision := range collisionCh {
	// 	i := collision.ProjIndex
	// 	j := collision.EnemyIndex

	// 	if g.player.Projectiles[i] != nil {
	// 		// Kill projectile
	// 		// TODO change back to false, just for fun
	// 		// g.player.Projectiles[i].IsAlive = false
	// 	}

	// 	if g.enemies[j] != nil {
	// 		// Respawn enemy
	// 		g.enemies[j] = g.enemies[j].CreateNewEnemy(g.tilemapJSON.GenValidPosOutsideCamera(g.player.X, g.player.Y))

	// 		// Create an extra new enemy
	// 		newEnemy := g.enemies[j].CreateNewEnemy(g.tilemapJSON.GenValidPosOutsideCamera(g.player.X, g.player.Y))
	// 		g.enemies = append(g.enemies, newEnemy)

	// 		g.gameStats.AliveEnemies += 2
	// 		g.gameStats.DeadEnemies += 1
	// 	}
	// }
	////////////////////////////////////////////////////

	// MouseButtonLeft Create new projectile
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		raw_cursor_x, raw_cursor_y := ebiten.CursorPosition()
		// must convert cursor pos to game pos
		cursor_x := float64(raw_cursor_x) - g.cam.X
		cursor_y := float64(raw_cursor_y) - g.cam.Y

		gross_diff_x := float64(cursor_x) - g.player.X
		gross_diff_y := float64(cursor_y) - g.player.Y

		// also math.Hypot(gross_diff_x, gross_diff_y)
		hypot := math.Sqrt(math.Pow(gross_diff_x, 2) + math.Pow(gross_diff_y, 2))
		if hypot == 0 {
			// consider doing something better?
			hypot = 1
			gross_diff_x = 1
		}

		speed := 240.0 // px / second

		sprite_Dx := (gross_diff_x / hypot) * speed
		sprite_Dy := (gross_diff_y / hypot) * speed

		g.player.Projectiles[projectile_counter] = &entities.Projectile{
			Sprite: &entities.Sprite{
				Img: g.projectileImg,
				X:   g.player.X + 5,
				Y:   g.player.Y + 12,
				Dx:  sprite_Dx / 60, // 60 fps/tps
				Dy:  sprite_Dy / 60,
			},
			Damage:  10,
			IsAlive: true,
		}

		g.gameStats.ProjectilesShot += 1
		projectile_counter = (projectile_counter + 1) % constants.NumberOfProjectiles
		// fmt.Println("projectile_counter: ", projectile_counter)
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, constants.CameraWidth, constants.CameraHeight)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*constants.Tilesize,
		float64(g.tilemapJSON.Layers[0].Height)*constants.Tilesize,
		constants.CameraWidth,
		constants.CameraHeight,
	)

	// Check if the player has been caught
	// for _, en := range g.enemies {
	// 	if en.IsAlive {
	// 		caught := tilemap.PosMatch(g.player.Sprite, en.Sprite)
	// 		if caught {
	// 			// fmt.Println("Enemy caught the player!")
	// 			fmt.Print("")
	// 		}
	// 	}
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// fill the screen with a nice sky color
	// screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	// loop over the layers
	for _, layer := range g.tilemapJSON.Layers {
		// loop over the tiles in the layer data
		for index, id := range layer.Data {

			// get the tile position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// convert the tile position to pixel position
			x *= 16
			y *= 16

			// get the position on the image where the tile id is
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			// convert the src tile pos to pixel src position
			srcX *= 16
			srcY *= 16

			// set the drawimageoptions to draw the tile at x, y
			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			// draw the tile
			screen.DrawImage(
				// cropping out the tile that we want from the spritesheet
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	// set the translation of our drawImageOptions to the player's position
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	playerFrame := 0
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		playerFrame = activeAnim.Frame()
	}

	// draw the player
	screen.DrawImage(
		// grab a subimage of the spritesheet
		g.player.Img.SubImage(
			g.player.SpriteSheet.Rect(playerFrame),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	// draw projectile(s)
	// Not good practice to have so many copies of the
	// image, and to have to scale each one
	for _, p := range g.player.Projectiles {
		if p != nil && p.IsAlive {
			// really should scale this before loading
			opts.GeoM.Scale(0.5, 0.5)
			opts.GeoM.Translate(p.X, p.Y)
			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(
				p.Sprite.Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Reset()

	// prints out enemies, but also doesn't?
	for _, sprite := range g.enemies {
		if sprite.IsAlive {
			opts.GeoM.Translate(sprite.X, sprite.Y)
			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			spriteFrame := 0
			activeAnim := sprite.ActiveAnimation(int(sprite.Dx), int(sprite.Dy))
			if activeAnim == nil {
				// force an up animation
				activeAnim = sprite.ActiveAnimation(0, 2)
			}
			spriteFrame = activeAnim.Frame()

			// draw the player
			screen.DrawImage(
				// grab a subimage of the spritesheet
				sprite.Img.SubImage(
					sprite.SpriteSheet.Rect(spriteFrame),
				).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Reset()

	// add transparent HUD
	g.drawHUD(screen)

}

// screen size/layout, not level
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	// Create a semi-transparent rectangle
	hud := ebiten.NewImage(140, 80)
	hud.Fill(color.RGBA{0, 0, 0, 128}) // black with 50% transparency (128/255)

	// Draw the rectangle at top-left corner
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(5, 5)
	// op.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(hud, op)

	// Now draw text (health, enemies, etc.)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Enemies: %d", g.gameStats.AliveEnemies), 20, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Kill Count: %d", g.gameStats.DeadEnemies), 20, 40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Projectiles: %d", g.gameStats.ProjectilesShot), 20, 60)
}
