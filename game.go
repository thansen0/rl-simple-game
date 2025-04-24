package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"simplegame/constants"
	"simplegame/entities"
	"simplegame/spritesheet"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	// the image and position variables for our player
	player               *entities.Player
	playerSpriteSheet    *spritesheet.SpriteSheet
	yellowBatSpriteSheet *spritesheet.SpriteSheet
	enemies              []*entities.Enemy
	tilemapJSON          *TilemapJSON
	tilemapImg           *ebiten.Image
	projectileImg        *ebiten.Image // may not be the best place for this
	cam                  *Camera
}

var projectile_counter int = 0

func (g *Game) Update() error {
	var wg sync.WaitGroup
	var prev_player_X float64
	var prev_player_Y float64
	prev_player_X = g.player.X
	prev_player_Y = g.player.Y

	// we modify projectiles again later one
	wg.Add(1)
	go func() {
		defer wg.Done()
		// update all projectiles
		g.player.UpdateAllProjectiles()
	}()

	// move the player based on keyboar input (left, right, up down)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X = g.tilemapJSON.GetValidXPos(g.player.X, -2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X = g.tilemapJSON.GetValidXPos(g.player.X, 2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y = g.tilemapJSON.GetValidYPos(g.player.Y, -2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y = g.tilemapJSON.GetValidYPos(g.player.Y, 2)
	}

	// needs to be updated for other functions
	g.player.Dx = g.player.X - prev_player_X
	g.player.Dy = g.player.Y - prev_player_Y

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	// add behavior to the enemies
	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			prev_player_X = sprite.X
			prev_player_Y = sprite.Y

			if sprite.X < g.player.X {
				sprite.X = g.tilemapJSON.GetValidXPos(sprite.X, 1)
			} else if sprite.X > g.player.X {
				sprite.X = g.tilemapJSON.GetValidXPos(sprite.X, -1)
			}
			if sprite.Y < g.player.Y {
				sprite.Y = g.tilemapJSON.GetValidYPos(sprite.Y, 1)
			} else if sprite.Y > g.player.Y {
				sprite.Y = g.tilemapJSON.GetValidYPos(sprite.Y, -1)
			}

			sprite.Dx = sprite.X - prev_player_X
			sprite.Dy = sprite.Y - prev_player_Y

			activeAnim := sprite.ActiveAnimation(int(sprite.Dx), int(sprite.Dy))
			if activeAnim == nil {
				// force an up animation
				activeAnim = sprite.ActiveAnimation(0, 2)
			}
			activeAnim.Update()
		}
	}

	// wait on updates before overwriting
	wg.Wait()

	// MouseButtonLeft Create new projectile IsMouseButtonPressed(ebiten.MouseButtonLeft) IsMouseButtonJustPressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursor_x, cursor_y := ebiten.CursorPosition()

		// NOTE: this math was largely done in my head (as all game breaking math should be)
		// and may not be accurate
		gross_diff_x := math.Abs(float64(cursor_x) - g.player.X)
		gross_diff_y := math.Abs(float64(cursor_y) - g.player.Y)

		fmt.Println("x: ", gross_diff_x, ", y: ", gross_diff_y)

		hypot := math.Sqrt(math.Pow(gross_diff_x, 2) + math.Pow(gross_diff_y, 2))
		hypot_adj_factor := math.Pow(hypot/15, 2)

		// fmt.Println("hypot_adj_factor: ", hypot_adj_factor)

		sprite_Dx := gross_diff_x / hypot_adj_factor
		sprite_Dy := gross_diff_y / hypot_adj_factor

		// sprite_Dx = 1
		// sprite_Dy = 1

		g.player.Projectiles[projectile_counter] = &entities.Projectile{
			Sprite: &entities.Sprite{
				Img: g.projectileImg,
				X:   g.player.X,
				Y:   g.player.Y,
				Dx:  sprite_Dx,
				Dy:  sprite_Dy,
			},
			Damage:  10,
			IsAlive: true,
		}

		projectile_counter = (projectile_counter + 1) % constants.NumberOfProjectiles
		fmt.Println("projectile_counter: ", projectile_counter)
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*constants.Tilesize,
		float64(g.tilemapJSON.Layers[0].Height)*constants.Tilesize,
		320,
		240,
	)

	for _, en := range g.enemies {
		caught := PosMatch(g.player.Sprite, en.Sprite)
		if caught {
			fmt.Println("Enemy caught the player!")
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// fill the screen with a nice sky color
	screen.Fill(color.RGBA{120, 180, 255, 255})

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
			g.playerSpriteSheet.Rect(playerFrame),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	// draw projectile(s)
	// Not good practice to have so many copies of the
	// image, and to have to scale each one
	for _, p := range g.player.Projectiles {
		if p != nil && p.IsAlive {
			opts.GeoM.Scale(0.5, 0.5)
			// I think this is right, however my code is wrong :(
			opts.GeoM.Translate(p.X, p.Y)
			// opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(
				p.Sprite.Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Reset()

	// technically only works for yellow bat enemies
	for _, sprite := range g.enemies {
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
				g.yellowBatSpriteSheet.Rect(spriteFrame),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

}

// screen size/layout, not level
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
