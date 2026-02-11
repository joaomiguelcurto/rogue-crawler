package engine

import (
	"log"
	"rogue-crawler/entities"

	"github.com/gdamore/tcell/v2"
)

type Game struct {
	Screen  tcell.Screen
	Enemies []entities.Enemy
}

// Helper function to render text to the tcell screen
func DrawText(s tcell.Screen, x, y int, str string, style tcell.Style) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, style)
	}
}

// Header
func RenderHeader(s tcell.Screen) {
	DrawText(s, 1, 1, "Rogue Crawler Debug View", tcell.StyleDefault)
	DrawText(s, 1, 2, "Press ESC to quit", tcell.StyleDefault)
}

func StartGame() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	// Initialize game structure
	g := &Game{
		Screen: screen,
		Enemies: []entities.Enemy{
			entities.NewEnemy("goblin"),
			entities.NewEnemy("spider"),
		},
	}

	for {
		screen.Clear()

		RenderHeader(screen)

		// Render enemies
		for i, enemy := range g.Enemies {
			row := 4 + (i * 2) // to space them out

			// Draw the symbol with the specific color
			style := tcell.StyleDefault.Foreground(enemy.Color)
			screen.SetContent(2, row, enemy.Symbol, nil, style)
		}

		screen.Show()

		ev := screen.PollEvent()
		if g.HandleInput(ev) {
			break
		}

	}
}
