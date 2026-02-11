package engine

import (
	"github.com/gdamore/tcell/v2"
)

// HandleInput processes keyboard events
func (g *Game) HandleInput(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		// Check if exit
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			return true // Signal to quit the game
		}

		// 2. Check for Movement (Arrows)
		switch ev.Key() {
		case tcell.KeyUp:
			// g.Player.Move(0, -1)
		case tcell.KeyDown:
			// g.Player.Move(0, 1)
		case tcell.KeyLeft:
			// g.Player.Move(-1, 0)
		case tcell.KeyRight:
			// g.Player.Move(1, 0)
		}

		// 3. Check for Actions (Runes/Letters)
		switch ev.Rune() {
		case 'a':
			// Start Attack Menu
		case 'i':
			// Open Inventory
		}

	case *tcell.EventResize:
		g.Screen.Sync()
	}

	return false // Do not quit
}
