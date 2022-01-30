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
	var spriteRes SpriteResource
	// spriteRes = LoadSprite("assets/goblin/goblin.png", "assets/goblin/goblin.csv")
	spriteRes, err = LoadSpriteResource("assets/gopher/gopher.png", "assets/gopher/gopher.csv")
	if err != nil {
		panic(err)
	}
	spriteRes.Debug()

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		sprites      []*pixel.Sprite
		matrices     []pixel.Matrix
	)

	var currIdx int = 0
	var chosenAnimation string = spriteRes.animationNames[currIdx]
	var numAnimations = len(spriteRes.animationNames)
	var currFrame int = 0

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			numFrames := len(spriteRes.animationMap[chosenAnimation])
			sprite := pixel.NewSprite(spriteRes.sheetPic, spriteRes.animationMap[chosenAnimation][currFrame])
			sprites = append(sprites, sprite)
			mouse := cam.Unproject(win.MousePosition())
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))

			if (currFrame + 1) > (numFrames - 1) {
				if (currIdx + 1) > (numAnimations - 1) {
					currIdx = 0
				} else {
					currIdx = currIdx + 1
				}
				currFrame = 0
			} else {
				currFrame = currFrame + 1
			}

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

		for i, tree := range sprites {
			tree.Draw(win, matrices[i])
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
