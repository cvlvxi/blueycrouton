package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/thoas/go-funk"
)

type Sprite struct {
	sheetPic     pixel.Picture
	animationMap map[string][]pixel.Rect
	frameHeight  float64
	frameWidth   float64
}

func (s *Sprite) Debug() {
	fmt.Println("Image Max X: ", s.sheetPic.Bounds().Max.X)
	fmt.Println("Image Max Y: ", s.sheetPic.Bounds().Max.Y)
	fmt.Println("frameHeight: ", s.frameHeight)
	fmt.Println("frameWidth: ", s.frameWidth)
	fmt.Println("Animation Map:")
	for name, frameList := range s.animationMap {
		fmt.Println("Animation: ", name)
		for idx, frame := range frameList {
			fmt.Println("[", idx, "]: ", frame)
		}
	}
}

func SpriteLoadErr(err error) (Sprite, error) {
	return Sprite{nil, nil, 0, 0}, err
}

func LoadSprite(sheetImagePath string, sheetCsv string) (Sprite, error) {
	var err error
	var sheetImg image.Image
	var sheetDef *SpriteDefinition
	var sheetPix *pixel.PictureData

	sheetFile, err := os.Open(sheetImagePath)
	defer sheetFile.Close()
	if err != nil {
		return SpriteLoadErr(err)
	}
	sheetImg, _, err = image.Decode(sheetFile)

	if err != nil {
		return SpriteLoadErr(err)
	}
	sheetPix = pixel.PictureDataFromImage(sheetImg)
	sheetDef, err = LoadSpriteCsvDefinition(sheetCsv)

	frameHeight := sheetPix.Bounds().Max.Y / float64(sheetDef.gridSize.rows)
	frameWidth := sheetPix.Bounds().Max.X / float64(sheetDef.gridSize.columns)

	// Max Items should be sheetDef.gridSize
	// rows x columns
	var frames [][]pixel.Rect = make([][]pixel.Rect, sheetDef.gridSize.rows)
	for rowIdx := range frames {
		frames[rowIdx] = make([]pixel.Rect, sheetDef.gridSize.columns)
	}

	animations := make(map[string][]pixel.Rect)
	var minX, minY, maxX, maxY float64
	for _, animation := range sheetDef.animations {
		name := animation.name
		if animation.numFrames < 1 {
			continue
		}
		minY = float64(animation.rowIdx) * frameHeight
		minX = float64(animation.colIdxStart) * frameWidth
		for frameIdx := 0; frameIdx <= animation.numFrames; frameIdx++ {
			nextX := float64(frameIdx+1) * frameWidth
			if nextX > sheetPix.Bounds().Max.X {
				minY = minY + frameHeight
				minX = 0
				maxX = frameWidth
			} else {
				maxX = nextX
			}
			maxY = minY + frameHeight
			animations[name] = append(animations[name], pixel.R(minX, minY, maxX, maxY))
		}
	}

	// var minX, minY, maxX, maxY float64
	// for rowIdx := range frames {
	// 	maxY = minY + frameHeight
	// 	for colIdx := range frames[rowIdx] {
	// 		maxX = minX + frameWidth
	// 		frames[rowIdx][colIdx] = pixel.R(minX, minY, maxX, maxY)
	// 		minX = maxX
	// 	}
	// 	minY = maxY
	// }

	// for name, def := range sheetDef.animationMap {
	// 	animations[name] = frames[def.rowIdx][def.colIdxStart:def.colIdxEnd]
	// }
	return Sprite{sheetPix, animations, frameHeight, frameWidth}, nil
}

func ReadCsvLine(r *csv.Reader) []string {
	line, err := r.Read()
	if err == io.EOF {
		fmt.Println("Error reading csv line")
	}
	return line
}

type GridSize struct {
	rows    int
	columns int
}

type AnimationDefinition struct {
	name        string
	rowIdx      int
	colIdxStart int
	numFrames   int
}

type SpriteDefinition struct {
	gridSize   GridSize
	animations []AnimationDefinition
}

// From CSV load the Sprite Definition
// Simple Definition
// Line 1: [num_rows, num_cols] e.g. 5, 11
// Line 2: [name, which_row_idx, start_idx, how_many]
// note: how_many can span multiple rows
func LoadSpriteCsvDefinition(sheetCsv string) (*SpriteDefinition, error) {

	var spriteDef SpriteDefinition

	fh, err := os.Open(sheetCsv)
	reader := csv.NewReader(fh)
	sizes := ReadCsvLine(reader)

	if len(sizes) != 2 {
		err = errors.New("Invalid Resource Size")
		return nil, err
	}
	rows, err := strconv.Atoi(sizes[0])
	columns, err := strconv.Atoi(sizes[1])
	spriteDef.gridSize = GridSize{rows, columns}
	spriteDef.animations = make([]AnimationDefinition, 0)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		normLine := funk.FilterString(line, func(x string) bool {
			return x != ""
		})
		if len(normLine) != 4 {
			err = errors.New(fmt.Sprint("Row definition for sprite invalid: ", line))
			return nil, err
		}

		name := normLine[0]
		rowIdx, err := strconv.Atoi(normLine[1])
		colIdxStart, err := strconv.Atoi(normLine[2])
		colIdxEnd, err := strconv.Atoi(normLine[3])

		def := AnimationDefinition{name, rowIdx, colIdxStart, colIdxEnd}
		spriteDef.animations = append(spriteDef.animations, def)
	}
	// Sort
	sort.Slice(spriteDef.animations, func(i, j int) bool {
		if spriteDef.animations[i].rowIdx != spriteDef.animations[j].rowIdx {
			return spriteDef.animations[i].rowIdx < spriteDef.animations[j].rowIdx
		} else {
			if spriteDef.animations[i].colIdxStart != spriteDef.animations[j].colIdxStart {
				return spriteDef.animations[i].colIdxStart < spriteDef.animations[j].colIdxStart
			} else {
				return true
			}
		}
	})

	return &spriteDef, nil
}
