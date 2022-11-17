package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jessevdk/go-flags"
)

type Bit struct {
	Value rune
	X     int
	Y     float32
	Z     uint8
}

var opts struct {
	Fps int64 `long:"fps" description:"Controls speed of animation." default:"20"`
}

func createBits(w, h int) (bits []Bit) {
	bits = make([]Bit, w*h/8)
	for i := range bits {
		bits[i] = Bit{
			Value: rune(rand.Intn(2)) + '0',
			X:     rand.Intn(w),
			Y:     rand.Float32() * float32(h),
			Z:     uint8(rand.Intn(220) + 35),
		}
	}

	return
}

func zero(slice []uint8) {
	for i := 0; i < len(slice); i++ {
		slice[i] = 0
	}
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		switch err := err.(type) {
		case flags.ErrorType:
			if err == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	screen, err := tcell.NewScreen()

	if err != nil {
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		os.Exit(1)
	}

	w, h := screen.Size()

	bits := createBits(w, h)

	depth := make([]uint8, w*h)

	defer func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}()

	events := make(chan tcell.Event)
	ticker := time.NewTicker(time.Second / time.Duration(opts.Fps))

	go screen.ChannelEvents(events, nil)

outer:
	for {
		select {
		case event := <-events:
			{
				switch event := event.(type) {
				case *tcell.EventKey:
					{
						break outer
					}
				case *tcell.EventResize:
					{
						w, h = event.Size()
						bits = createBits(w, h)
						depth = make([]uint8, w*h)
					}
				}
			}
		case <-ticker.C:
			{
				screen.Clear()
				zero(depth)

				for i := range bits {
					bits[i].Y += float32(bits[i].Z) / 255

					x := bits[i].X
					y := int(bits[i].Y)
					z := bits[i].Z

					if !(y >= 0 && y < h) {
						bits[i].Y = 0
						continue
					}

					if depth[x+y*w] >= z {
						continue
					}

					depth[x+y*w] = z

					color := tcell.NewRGBColor(0, int32(z), 0)
					style := tcell.StyleDefault.Foreground(color)

					depth[x+y*w] = z

					screen.SetContent(
						x,
						y,
						bits[i].Value,
						nil,
						style,
					)
				}

				screen.Show()
			}
		}
	}
}
