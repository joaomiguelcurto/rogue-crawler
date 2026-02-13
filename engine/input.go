package engine

import (
	"github.com/gdamore/tcell/v2"
)

// HandleInput processes keyboard events
func (g *Game) HandleInput(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		// Check if exit
		if ev.Key() == tcell.KeyEscape {
			if g.State == "Menu" {
				g.State = "Exploring"
				g.IsPickingSlot = false // Reset slot picking if exiting menu
				return false
			}
			return true
		}

		if g.State == "Exploring" {
			return g.handleExplorationInput(ev)
		} else if g.State == "Menu" {
			return g.handleMenuInput(ev)
		}
	}
	return false
}

// Move handles the actual coordinate update with collision detection
func (g *Game) MovePlayer(dx, dy int) {
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	if newX >= 0 && newX < g.Map.Width && newY >= 0 && newY < g.Map.Height {
		tile := g.Map.Tiles[newX][newY]

		if tile.Type == "Floor" {
			g.Player.X = newX
			g.Player.Y = newY
		} else if tile.Type == "Portal" {
			// Move to the next floor
			g.NextFloor()
		}
	}
}

// Add this to engine/game.go or input.go
func (g *Game) NextFloor() {
	g.Floor++
	// Re-generate map
	g.Map = NewMap(80, 45)
	g.Map.GenerateLevel(g)

	// Reset player position to the new room 0
	if len(g.Map.Rooms) > 0 {
		g.Player.X, g.Player.Y = g.Map.Rooms[0].Center()
	}
}

func (g *Game) handleExplorationInput(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyUp:
		g.MovePlayer(0, -1)
	case tcell.KeyDown:
		g.MovePlayer(0, 1)
	case tcell.KeyLeft:
		g.MovePlayer(-1, 0)
	case tcell.KeyRight:
		g.MovePlayer(1, 0)
	}

	if ev.Rune() == 'm' || ev.Rune() == 'M' {
		g.State = "Menu"
	}
	return false
}

func (g *Game) handleMenuInput(ev *tcell.EventKey) bool {
	// Standard Tab Switching (Disabled while picking a slot for a cleaner UI experience)
	if !g.IsPickingSlot {
		switch ev.Rune() {
		case '1':
			g.ActiveTab = 0
			return false
		case '2':
			g.ActiveTab = 1
			return false
		case '3':
			g.ActiveTab = 2
			return false
		}
	}

	// Skill Tab Logic
	if g.ActiveTab == 2 {
		if !g.IsPickingSlot {
			switch ev.Key() {
			case tcell.KeyUp:
				if g.SelectedSkillIndex > 0 {
					g.SelectedSkillIndex--
				}
			case tcell.KeyDown:
				if g.SelectedSkillIndex < len(g.Player.LearnedSkills)-1 {
					g.SelectedSkillIndex++
				}
			case tcell.KeyEnter:
				if len(g.Player.LearnedSkills) > 0 {
					g.IsPickingSlot = true
				}
			}
		} else {
			switch ev.Key() {
			case tcell.KeyUp:
				if g.SelectedSlotIndex > 0 {
					g.SelectedSlotIndex--
				}
			case tcell.KeyDown:
				if g.SelectedSlotIndex < g.Player.MaxActiveSlots-1 {
					g.SelectedSlotIndex++
				}
			case tcell.KeyEscape:
				g.IsPickingSlot = false
			case tcell.KeyEnter:
				g.Player.EquipSkill(g.SelectedSkillIndex, g.SelectedSlotIndex)
				g.IsPickingSlot = false
			}
		}
	}
	return false
}
