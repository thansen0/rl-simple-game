package main

import (
	"encoding/json"
	"os"
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
	TileWidth  int                `json:"tilewidth"`
	TileHeight int                `json:"tileheight"`
}

func clamp(value, min, max float64) float64 {
	// assume character is width 2
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
	return clamp(cur_pos+move_command, 2, float64(tm.Layers[0].Width*tm.TileWidth)-18)
}

// calculates a valid Y position for a sprite
func (tm *TilemapJSON) GetValidYPos(cur_pos, move_command float64) float64 {
	// Assume all sprites are 16 pixels tall, 2 pixel buffer
	return clamp(cur_pos+move_command, 2, float64(tm.Layers[0].Height*tm.TileHeight)-18)
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
