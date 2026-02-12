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
				return false // Signal to exit the menu
			}
			return true // Signal to quit the game
		}

		if g.State == "Exploring" {
			return g.handleExplorationInput(ev)
		} else if g.State == "Menu" {
			return g.handleMenuInput(ev)
		}
	}
	return false
}

func (g *Game) handleExplorationInput(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyUp:
		g.Player.Move(0, -1)
	case tcell.KeyDown:
		g.Player.Move(0, 1)
	case tcell.KeyLeft:
		g.Player.Move(-1, 0)
	case tcell.KeyRight:
		g.Player.Move(1, 0)
	}

	if ev.Rune() == 'm' || ev.Rune() == 'M' {
		g.State = "Menu"
	}
	return false
}

func (g *Game) handleMenuInput(ev *tcell.EventKey) bool {
	// Standard Tab Switching
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

	if g.ActiveTab == 2 {
		if !g.IsPickingSlot {
			// Select the Skill
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
				g.IsPickingSlot = true // Move to slot selection
			}
		} else {
			// Select the Slot
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
				g.IsPickingSlot = false // Cancel
			case tcell.KeyEnter:
				g.Player.EquipSkill(g.SelectedSkillIndex, g.SelectedSlotIndex)
				g.IsPickingSlot = false // Done
			}
		}
	}
	return false
}
