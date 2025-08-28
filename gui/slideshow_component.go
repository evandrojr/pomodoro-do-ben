package gui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// SlideshowComponent is a reusable and responsive slideshow component.
type SlideshowComponent struct {
	container *fyne.Container
	img *canvas.Image
	pics []string
	current int
	ticker *time.Ticker
}

// NewSlideshowComponent creates a new SlideshowComponent.
func NewSlideshowComponent(imagePaths []string) *SlideshowComponent {
	img := canvas.NewImageFromFile(imagePaths[0])
	img.FillMode = canvas.ImageFillContain

	sc := &SlideshowComponent{
		img: img,
		pics: imagePaths,
		current: 0,
	}
	sc.container = container.NewMax(img)

	sc.startSlideshow()

	return sc
}

// GetContent returns the Fyne CanvasObject for the slideshow.
func (sc *SlideshowComponent) GetContent() fyne.CanvasObject {
	return sc.container
}

func (sc *SlideshowComponent) startSlideshow() {
	sc.ticker = time.NewTicker(3 * time.Second)
	go func() {
		defer sc.ticker.Stop()
		for range sc.ticker.C {
			sc.current = (sc.current + 1) % len(sc.pics)
			sc.img.File = sc.pics[sc.current]
			sc.img.Refresh()
		}
	}()
}

// StopSlideshow stops the slideshow animation.
func (sc *SlideshowComponent) StopSlideshow() {
	if sc.ticker != nil {
		sc.ticker.Stop()
	}
}
