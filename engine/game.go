package engine

import (
	"fmt"
	"log"
	"rogue-crawler/entities"

	"github.com/gdamore/tcell/v2"
)

type Game struct {
	Screen             tcell.Screen
	Player             *entities.Player
	Enemies            []entities.Enemy
	State              string // "Exploring", "Menu"
	ActiveTab          int    // 0: Stats, 1: Inventory, 2: Skills
	SelectedSkillIndex int    // Track which skill is selected
	SelectedSlotIndex  int
	IsPickingSlot      bool // State toggle
}

// Helper function to render text to the tcell screen
func DrawText(s tcell.Screen, x, y int, str string, style tcell.Style) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, style)
	}
}

// Exploration Mode
func (g *Game) renderExploration() {
	// Render Player
	g.Screen.SetContent(g.Player.X, g.Player.Y, g.Player.Visual, nil, tcell.StyleDefault.Foreground(g.Player.Color))
	DrawText(g.Screen, 1, 1, "MODE: Exploring - Press 'M' for Menu", tcell.StyleDefault)

	// Render enemies
	for i, enemy := range g.Enemies {
		row := 4 + (i * 2) // to space them out

		// Draw the visual with the specific color
		style := tcell.StyleDefault.Foreground(enemy.Color)
		g.Screen.SetContent(2, row, enemy.Visual, nil, style)
	}
}

// Render Menu
func (g *Game) renderMenu() {
	// Draw Menu Border/Background
	DrawText(g.Screen, 2, 2, "[ 1: Stats ]  [ 2: Inventory ]  [ 3: Skills ] (Press ESC to Close)", tcell.StyleDefault.Reverse(true))

	switch g.ActiveTab {
	case 0:
		g.drawStats()
	case 1:
		DrawText(g.Screen, 2, 5, "-- Inventory Slot 1: Empty --", tcell.StyleDefault)
	case 2:
		g.drawSkills()
	}
}

func (g *Game) drawStats() {
	p := g.Player
	DrawText(g.Screen, 2, 5, fmt.Sprintf("Class: %s", p.Class), tcell.StyleDefault)
	DrawText(g.Screen, 2, 6, fmt.Sprintf("Level: %d (XP: %d)", p.Level, p.XP), tcell.StyleDefault)
	DrawText(g.Screen, 2, 8, fmt.Sprintf("HP: %d/%d", p.HP, p.MaxHP), tcell.StyleDefault.Foreground(tcell.ColorRed))
	DrawText(g.Screen, 2, 9, fmt.Sprintf("MP: %d/%d", p.MP, p.MaxMP), tcell.StyleDefault.Foreground(tcell.ColorBlue))
	DrawText(g.Screen, 2, 11, fmt.Sprintf("Strength:     %d", p.Strength), tcell.StyleDefault)
	DrawText(g.Screen, 2, 12, fmt.Sprintf("Intelligence: %d", p.Intelligence), tcell.StyleDefault)
	DrawText(g.Screen, 2, 13, fmt.Sprintf("Agility:      %d", p.Agility), tcell.StyleDefault)
	DrawText(g.Screen, 2, 14, fmt.Sprintf("Luck:         %d", p.Luck), tcell.StyleDefault)
}

func (g *Game) drawSkills() {
	// Draw Active Slots
	DrawText(g.Screen, 2, 4, "--- ACTIVE SLOTS ---", tcell.StyleDefault.Foreground(tcell.ColorYellow))

	for i := 0; i < g.Player.MaxActiveSlots; i++ {
		row := 5 + i

		// DEFAULT STYLE
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		prefix := fmt.Sprintf("[%d] ", i+1)

		// TARGETING STYLE (When picking a slot)
		if g.IsPickingSlot && i == g.SelectedSlotIndex {
			// Explicitly set colors here to avoid the bug where the text is invisible
			style = tcell.StyleDefault.
				Foreground(tcell.ColorBlack).
				Background(tcell.ColorYellow).
				Bold(true)
			prefix = " TARGET -> "
		}

		if i < len(g.Player.ActiveSkills) && g.Player.ActiveSkills[i].Name != "Empty" {
			s := g.Player.ActiveSkills[i]
			DrawText(g.Screen, 2, row, prefix+s.Name, style)
		} else {
			// Make empty slots gray unless they are being targeted
			if !(g.IsPickingSlot && i == g.SelectedSlotIndex) {
				style = style.Foreground(tcell.ColorGray)
			}
			DrawText(g.Screen, 2, row, prefix+"-- EMPTY --", style)
		}
	}

	// Draw Learned Skills (Dimmed if picking a slot)
	listStyle := tcell.StyleDefault.Foreground(tcell.ColorLightCyan)
	if g.IsPickingSlot {
		listStyle = listStyle.Foreground(tcell.ColorGray)
	}
	DrawText(g.Screen, 2, 11, "--- LEARNED SKILLS ---", listStyle)

	for i, s := range g.Player.LearnedSkills {
		row := 12 + i
		style := tcell.StyleDefault
		if !g.IsPickingSlot && i == g.SelectedSkillIndex {
			style = tcell.StyleDefault.Reverse(true)
		}
		DrawText(g.Screen, 2, row, "- "+s.Name, style)
	}
}

func Run() {
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
		Player: entities.NewPlayer("rogue"),
		Enemies: []entities.Enemy{
			entities.NewEnemy("goblin"),
			entities.NewEnemy("spider"),
		},
		State:     "Exploring",
		ActiveTab: 0,
	}

	for {
		screen.Clear()

		// Renders depending on the state
		switch g.State {
		case "Exploring":
			g.renderExploration()
		case "Menu":
			g.renderMenu()
		}

		screen.Show()

		// Handle all the events
		ev := screen.PollEvent()
		if g.HandleInput(ev) {
			break
		}

	}
}
