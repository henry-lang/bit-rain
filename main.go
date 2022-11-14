package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	fps := 1

	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defer func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}()

	events := make(chan tcell.Event)
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	c := 0

	go screen.ChannelEvents(events, nil)

outer:
	for {
		select {
		case event := <-events:
			{
				switch event.(type) {
				case *tcell.EventKey:
					{
						break outer
					}
				}
			}
		case <-ticker.C:
			{
				c += 1
				screen.Clear()
				screen.SetContent(int(c), 0, 'c', nil, tcell.StyleDefault)
				screen.Show()
			}
		}
	}
}
