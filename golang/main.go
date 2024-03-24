package main

import (
	_ "image/jpeg"
	"log"
	"math/rand"
	"sync"

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

func QuickSort(g *Game, pixels []Pixel, start int, stop int) {
	if len(pixels) <= 1 {
		return
	}

	pivot := pixels[len(pixels)/2]

	left := make([]Pixel, 0)
	middle := make([]Pixel, 0)
	right := make([]Pixel, 0)

	for i := 0; i < len(pixels); i++ {
		if pixels[i].Weigth < pivot.Weigth {
			left = append(left, pixels[i])
		} else if pixels[i].Weigth == pivot.Weigth {
			middle = append(middle, pixels[i])
		} else {
			right = append(right, pixels[i])
		}
	}

	len_l := len(left)

	start_m := start + len_l
	len_m := len(middle)

	start_r := start + len_l + len_m
	len_r := len(right)

	g.Mutex.Lock()
	for i := 0; i < len_l; i++ {
		g.Pixels[start+i] = left[i]
	}
	for i := 0; i < len_m; i++ {
		g.Pixels[start_m+i] = middle[i]
	}
	for i := 0; i < len_r; i++ {
		g.Pixels[start_r+i] = right[i]
	}
	g.Mutex.Unlock()

	// The 2 sub-sort can be launched in other goroutines with "go QuickSort" to improve speed
	// But to better see the sorting, keep^the execution in the same goroutine
	QuickSort(g, left, start, start+len_l)
	QuickSort(g, right, start_r, stop)
}

type Game struct {
	Img      *ebiten.Image
	NbPixels int
	Pixels   []Pixel
	Init     bool
	Index    int
	Mutex    sync.Mutex
}

func (g *Game) Update() error {
	if !g.Init {
		g.Mutex.Lock()
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
		g.Mutex.Unlock()

		go QuickSort(g, g.Pixels, 0, g.NbPixels)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Mutex.Lock()
	screen.WritePixels(GetRawBytes(g.Pixels))
	g.Mutex.Unlock()
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
