package loader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DemmyDemon/imgz/internal/do"

	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/webp"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GalleryImage struct {
	Images []*ebiten.Image
	Delays []int
	index  int
	last   int64
}

func (g *GalleryImage) Get() *ebiten.Image {
	if len(g.Images) == 1 {
		return g.Images[0]
	}
	now := time.Now().UnixMilli() / 10
	if now >= g.last+int64(g.Delays[g.index]) {
		g.last = now
		g.index++
		if g.index >= len(g.Delays) {
			g.index = 0
		}
	}
	// do.Verbose("Images len:", len(g.Images), "Index:", g.index)
	return g.Images[g.index]
}

func Load(path string) (GalleryImage, error) {
	ext := filepath.Ext(strings.ToLower(path))
	fd, err := os.Open(path)
	if err != nil {
		return GalleryImage{}, err
	}
	defer fd.Close()

	switch ext {
	case ".jpg", ".jpeg":
		return LoadJPEG(fd)
	case ".gif":
		return LoadGIF(fd)
	case ".png":
		return LoadPNG(fd)
	case ".webp":
		return LoadWEBP(fd)
	}

	return GalleryImage{}, fmt.Errorf("could not guess filetype from filename: %s", path)
}

func LoadGIF(r io.Reader) (GalleryImage, error) {
	img, err := gif.DecodeAll(r)
	if err != nil {
		return GalleryImage{}, err
	}
	frames := make([]*ebiten.Image, len(img.Image))
	delays := make([]int, len(img.Delay))
	copy(delays, img.Delay)
	do.Verbose("Loaded GIF, trying to make ebiten Image of", len(img.Image), "frames")
	var prevFrame *ebiten.Image
	for i, frame := range img.Image {
		if frame == nil {
			break
		}
		eImg := ebiten.NewImageFromImage(frame)
		if eImg == nil {
			return GalleryImage{}, fmt.Errorf("failed to make *ebiten.Image out of frame %d", i)
		}
		if prevFrame != nil {
			prevFrame.DrawImage(eImg, nil)
		} else {
			prevFrame = eImg
		}
		frames[i] = ebiten.NewImageFromImage(prevFrame)
		// do.Verbose("Frame", i, "done")
	}
	do.Verbose("All done frobnicating GIF data")
	return GalleryImage{
		Images: frames,
		Delays: delays,
	}, nil
}
func LoadJPEG(r io.Reader) (GalleryImage, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return GalleryImage{}, err
	}
	eImg := ebiten.NewImageFromImage(img)
	return GalleryImage{
		Images: []*ebiten.Image{eImg},
		Delays: []int{0},
	}, nil
}
func LoadPNG(r io.Reader) (GalleryImage, error) {
	img, err := png.Decode(r)
	if err != nil {
		return GalleryImage{}, err
	}
	eImg := ebiten.NewImageFromImage(img)
	return GalleryImage{
		Images: []*ebiten.Image{eImg},
		Delays: []int{0},
	}, nil
}
func LoadWEBP(r io.Reader) (GalleryImage, error) {
	img, err := webp.Decode(r)
	if err != nil {
		return GalleryImage{}, err
	}
	eImg := ebiten.NewImageFromImage(img)
	return GalleryImage{
		Images: []*ebiten.Image{eImg},
		Delays: []int{0},
	}, nil
}

func LoadOld(path string) (*ebiten.Image, error) {
	do.Verbose(path)
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		fmt.Printf("%s: %s\n", path, err)
		return nil, err
	}
	return img, nil
}
