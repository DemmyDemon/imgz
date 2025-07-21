package main

import (
	"fmt"
	"hash/fnv"
	"image/color"
	"math/rand/v2"
	"os"
	"runtime"

	"github.com/DemmyDemon/imgz/internal/do"
	"github.com/DemmyDemon/imgz/loader"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	StepSize      = 2.5
	PulseInterval = 45
)
const (
	AutoModeOff = iota
	AutoModeOn
	AutoModeRandom
)

var (
	hasher = fnv.New32a()
)

type Settings struct {
	OffsetX float64
	OffsetY float64
	Scale   float64
}

func getSum(path string) uint32 {
	hasher.Reset()
	hasher.Write([]byte(path))
	return hasher.Sum32()
}

type Game struct {
	auto      int
	autoDelay int
	tick      int
	lastPulse int
	image     loader.GalleryImage
	ScreenW   int
	ScreenH   int
	settings  map[uint32]*Settings
	paths     []string
	index     int
}

func NewGame(paths []string, title string) *Game {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("imgz " + title)
	ebiten.SetFullscreen(true)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if len(paths) == 0 {
		fmt.Println("No relevant images found")
		os.Exit(1)
	}
	imgz, err := loader.Load(paths[0])
	do.Fuck(err)
	settings := make(map[uint32]*Settings, len(paths))
	return &Game{
		paths:     paths,
		image:     imgz,
		settings:  settings,
		index:     0,
		autoDelay: 100,
	}
}

func (g *Game) currentSettings() *Settings {
	path := g.paths[g.index]
	sum := getSum(path)
	if set, ok := g.settings[sum]; ok {
		return set
	}
	do.Verbose("No settings entry for", path)
	image := g.image.Get()
	if image == nil {
		panic("Okay, somehow the image getter returned a nil image")
	}
	scale := float64(g.ScreenW) / float64(image.Bounds().Dx())
	height := float64(g.ScreenH) / float64(image.Bounds().Dy())
	fullWidth := true
	if height < scale {
		scale = height
		fullWidth = false
	}

	offsetX := 0.0
	offsetY := 0.0
	if fullWidth {
		origHeight := image.Bounds().Dy()
		scaledHeight := (float64(origHeight) * scale)
		halfScreen := float64(g.ScreenH) / 2
		halfScaledImg := scaledHeight / 2
		offsetY = halfScreen - halfScaledImg
		do.Verbose("> Full width. Vertical offset", offsetY, "Height", scaledHeight, "Originally", origHeight, "Scale", scale)
	} else {
		origWidth := image.Bounds().Dx()
		scaledWidth := (float64(origWidth) * scale)
		halfScreen := float64(g.ScreenW) / 2
		halfScaledImg := scaledWidth / 2
		offsetX = halfScreen - halfScaledImg
		do.Verbose("> Not full width. Horizontal offset", offsetX, "Width", scaledWidth, "Originally", origWidth, "Scale", scale)
	}

	fresh := &Settings{
		OffsetX: offsetX,
		OffsetY: offsetY,
		Scale:   scale,
	}
	g.settings[sum] = fresh
	return fresh
}

func (g *Game) reframe() {
	path := g.paths[g.index]
	sum := getSum(path)
	delete(g.settings, sum)
	g.currentSettings()
}

func (g *Game) drawSegment(screen *ebiten.Image, start, end int, col color.RGBA) {
	bottom := screen.Bounds().Dy()
	for i := start; i <= end; i++ {
		for j := 0; j <= 4; j++ {
			screen.Set(i, bottom-j, col)
		}
	}
}

func (g *Game) drawSegments(screen *ebiten.Image) {
	numImg := len(g.paths)
	col := color.RGBA{128, 128, 128, 128}
	segWidth := screen.Bounds().Dx() / numImg
	offset := (screen.Bounds().Dx() - (segWidth * numImg)) / 2
	for i := 0; i < numImg; i++ {
		if i == g.index {
			switch g.auto {
			case AutoModeOff:
				col = color.RGBA{255, 0, 0, 255}
			case AutoModeOn:
				col = color.RGBA{0, 255, 0, 255}
			case AutoModeRandom:
				col = color.RGBA{255, 255, 0, 255}
			default:
				col = color.RGBA{255, 255, 255, 255}
			}
		} else if i > g.index {
			col = color.RGBA{72, 72, 72, 128}
		}
		g.drawSegment(screen, offset+i*segWidth, offset+(i*segWidth)+segWidth, col)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	set := g.currentSettings()
	opt := &ebiten.DrawImageOptions{
		GeoM: ebiten.GeoM{},
	}

	opt.GeoM.Scale(set.Scale, set.Scale)
	opt.GeoM.Translate(set.OffsetX, set.OffsetY)
	img := g.image.Get()
	screen.DrawImage(img, opt)
	g.drawSegments(screen)
	if do.Verbosity {
		switch g.auto {
		case AutoModeOn:
			ebitenutil.DebugPrint(screen, fmt.Sprintf("(auto %d) %d/%d %s", g.autoDelay, g.index+1, len(g.paths), g.paths[g.index]))
		case AutoModeRandom:
			ebitenutil.DebugPrint(screen, fmt.Sprintf("RANDOM %d %s %d", g.index+1, g.paths[g.index], g.autoDelay))
		default:
			ebitenutil.DebugPrint(screen, fmt.Sprintf("%d/%d %s", g.index+1, len(g.paths), g.paths[g.index]))
		}
	}
}

func (g *Game) Update() error {
	g.tick++

	if g.auto > 0 && g.tick%g.autoDelay == 0 {
		switch g.auto {
		case AutoModeOn:
			g.nextImage(true)
		case AutoModeRandom:
			g.randomImage()
		default:
			g.auto = 0
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if runtime.GOARCH != "wasm" { // Because this will exit the program, but not exit fullscreen or clear the canvas...
			os.Exit(0)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		if ebiten.IsFullscreen() {
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.reframe()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.auto++
		if g.auto > 2 {
			g.auto = 0
		}
		do.Verbose("Auto mode: ", g.auto)
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			set := g.currentSettings()
			set.OffsetY -= StepSize
			do.Verbose("Up", set.OffsetY)
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			set := g.currentSettings()
			set.OffsetY += StepSize
			do.Verbose("Down", set.OffsetY)
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			set := g.currentSettings()
			set.OffsetX -= StepSize
			do.Verbose("Left", set.OffsetX)
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			set := g.currentSettings()
			set.OffsetX += StepSize
			do.Verbose("Right", set.OffsetX)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		set := g.currentSettings()
		set.Scale -= StepSize / 250
		if set.Scale < 0.05 {
			set.Scale = 0.05
		}
		do.Verbose("Out", set.Scale)
	}
	if ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		set := g.currentSettings()
		set.Scale += StepSize / 250
		if set.Scale > 8.0 {
			set.Scale = 8.0
		}
		do.Verbose("In", set.Scale)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.nextImage(inpututil.IsKeyJustPressed(ebiten.KeyArrowRight))
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.prevImage(inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft))
	}
	if g.auto > 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.autoDelay -= 10
			if g.autoDelay < 10 {
				g.autoDelay = 10
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.autoDelay += 10
			if g.autoDelay > 1000 {
				g.autoDelay = 1000
			}
		}
	}
	return nil
}

func (g *Game) randomImage() {
	if len(g.paths) == 1 {
		do.Verbose("Picking a random one out of ONE image is easy! Done!")
		return
	}
	idx := g.index
	for idx == g.index {
		idx = rand.IntN(len(g.paths))
	}
	g.index = idx
	do.Verbose("Picked", idx, "at random")
	g.mustLoad()
}

func (g *Game) nextImage(skipPulse bool) {
	if !skipPulse && g.tick < g.lastPulse+PulseInterval {
		return
	}
	g.lastPulse = g.tick
	g.index++
	if g.index >= len(g.paths) {
		g.index = 0
	}
	do.Verbose("Next!")
	g.mustLoad()
}

func (g *Game) prevImage(skipPulse bool) {
	if !skipPulse && g.tick < g.lastPulse+PulseInterval {
		return
	}
	g.lastPulse = g.tick
	g.index--
	if g.index < 0 {
		g.index = len(g.paths) - 1
	}
	do.Verbose("Previous!")
	g.mustLoad()
}

func (g *Game) mustLoad() {
	imgz, err := loader.Load(g.paths[g.index])
	do.Fuck(err)
	g.image = imgz
}

func (g *Game) Layout(width, height int) (int, int) {
	g.ScreenH = height
	g.ScreenW = width
	return width, height
}
