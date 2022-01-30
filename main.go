package main

import (
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	var sprite Sprite
	// sprite = LoadSprite("assets/goblin/goblin.png", "assets/goblin/goblin.csv")
	sprite, err = LoadSprite("assets/gopher/gopher.png", "assets/gopher/gopher.csv")
	if err != nil {
		panic(err)
	}
	sprite.Debug()

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		trees        []*pixel.Sprite
		matrices     []pixel.Matrix
	)

	var names []string = []string{"WalkRight", "WalkRight", "WalkUp", "WalkLeft"}
	var nameIdx int = 0
	var frameIdx int = 0

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			tree := pixel.NewSprite(sprite.sheetPic, sprite.animationMap[names[nameIdx]][frameIdx])
			// chosenName := names[nameIdx]
			// fmt.Println(chosenName)
			if frameIdx == len(sprite.animationMap[names[nameIdx]])-2 {
				if nameIdx > len(names)-2 {
					nameIdx = 0
				} else {
					nameIdx = nameIdx + 1
				}
				frameIdx = 0
			} else {
				frameIdx = frameIdx + 1
			}
			trees = append(trees, tree)
			mouse := cam.Unproject(win.MousePosition())
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))

		}
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(colornames.Forestgreen)

		for i, tree := range trees {
			tree.Draw(win, matrices[i])
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
