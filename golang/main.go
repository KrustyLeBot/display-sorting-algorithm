package main

import (
	_ "image/jpeg"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Pixel struct {
	R      byte
	G      byte
	B      byte
	A      byte
	Weigth int
}

func GetRawBytes(pixels []Pixel) []byte {
	len := len(pixels)
	tmpPixels := make([]byte, 4*len)

	for i := 0; i < len; i++ {
		tmpPixels[(i * 4)] = pixels[i].R
		tmpPixels[(i*4)+1] = pixels[i].G
		tmpPixels[(i*4)+2] = pixels[i].B
		tmpPixels[(i*4)+3] = pixels[i].A
	}

	return tmpPixels
}

type Game struct {
	Img      *ebiten.Image
	NbPixels int
	Pixels   []Pixel
	Init     bool
	Index    int
}

func (g *Game) Update() error {
	if !g.Init {
		g.Init = true
		g.Pixels = make([]Pixel, g.NbPixels)
		tmpPixels := make([]byte, 4*g.NbPixels)
		g.Img.ReadPixels(tmpPixels)

		_weigth := 0
		for i := 0; i < g.NbPixels; i++ {
			g.Pixels[i].R = tmpPixels[(i * 4)]
			g.Pixels[i].G = tmpPixels[(i*4)+1]
			g.Pixels[i].B = tmpPixels[(i*4)+2]
			g.Pixels[i].A = tmpPixels[(i*4)+3]
			g.Pixels[i].Weigth = _weigth
			_weigth++
		}

		// Shuffle the array
		rand.Shuffle(g.NbPixels, func(i, j int) { g.Pixels[i], g.Pixels[j] = g.Pixels[j], g.Pixels[i] })
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(GetRawBytes(g.Pixels))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 1024
}

func NewGame() *Game {
	g := &Game{}

	img, _, err := ebitenutil.NewImageFromFile("../example.jpg")
	if err != nil {
		log.Fatal(err)
	}
	g.NbPixels = 1024 * 1024
	g.Img = img
	g.Init = false
	g.Index = 0
	return g
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(1024, 1024)
	ebiten.SetWindowTitle("Display sorting algorithm")

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
