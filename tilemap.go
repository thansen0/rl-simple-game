package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"simplegame/constants"
	"simplegame/entities"
)

// data we want for one layer in our list of layers
type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers     []TilemapLayerJSON `json:"layers"`
	TileWidth  int                `json:"tilewidth"`  // x position
	TileHeight int                `json:"tileheight"` // y position
}

// makes sure a float64 is between a min and max range
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// calculates a valid X position for a sprite
func (tm *TilemapJSON) GetValidXPos(cur_pos, move_command float64) float64 {
	// fmt.Printf("w: %i, tile width: %i \n", tm.Layers[0].Width, tm.TileWidth)
	// Assume all sprites are 16 pixels tall, 2 pixel buffer
	return clamp(cur_pos+move_command, constants.CharacterBuffer,
		float64(tm.Layers[0].Width*tm.TileWidth)-(constants.Tilesize+constants.CharacterBuffer))
}

// calculates a valid Y position for a sprite
func (tm *TilemapJSON) GetValidYPos(cur_pos, move_command float64) float64 {
	// Assume all sprites are 16 pixels tall, 2 pixel buffer
	return clamp(cur_pos+move_command, constants.CharacterBuffer,
		float64(tm.Layers[0].Height*tm.TileHeight)-(constants.Tilesize+constants.CharacterBuffer))
}

// This may fail as map constraints grow
func (tm *TilemapJSON) GenValidXPos() float64 {
	var rand_x_pos float64 = float64(rand.Intn(tm.Layers[0].Width * tm.TileWidth))

	return clamp(rand_x_pos, constants.CharacterBuffer,
		float64(tm.Layers[0].Width*tm.TileWidth)-(constants.Tilesize+constants.CharacterBuffer))
}

func (tm *TilemapJSON) GenValidYPos(_ float64) float64 {
	var rand_y_pos float64 = float64(rand.Intn(tm.Layers[0].Height * tm.TileHeight))

	return clamp(rand_y_pos, constants.CharacterBuffer,
		float64(tm.Layers[0].Height*tm.TileHeight)-(constants.Tilesize+constants.CharacterBuffer))
}

func (tm *TilemapJSON) GenValidPos() (float64, float64) {
	var rand_x_pos float64 = tm.GenValidXPos()
	fmt.Println("X:", rand_x_pos)
	return rand_x_pos, tm.GenValidYPos(rand_x_pos)
}

func PosMatch(s1, s2 *entities.Sprite) bool {
	var x_pos_match bool = math.Abs(s1.X-s2.X) < constants.Tilesize/2
	var y_pos_match bool = math.Abs(s1.Y-s2.Y) < constants.Tilesize/2

	return x_pos_match && y_pos_match
}

// opens the file, parses it, and returns the json object + potential error
func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
