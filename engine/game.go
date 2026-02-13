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
	SelectedSkillIndex int    // Track which learned skill is highlighted
	SelectedSlotIndex  int    // Track which active slot is highlighted
	IsPickingSlot      bool   // True if we are choosing a slot for a skill
	Map                *GameMap
	Floor              int // Current dungeon level
}

// DrawText is a helper to render strings to the screen
func DrawText(s tcell.Screen, x, y int, str string, style tcell.Style) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, style)
	}
}

// renderExploration draws the dungeon map, player, and enemies
func (g *Game) renderExploration() {
	// Draw the Map Tiles
	for x := 0; x < g.Map.Width; x++ {
		for y := 0; y < g.Map.Height; y++ {
			tile := g.Map.Tiles[x][y]
			g.Screen.SetContent(x, y, tile.Visual, nil, tcell.StyleDefault.Foreground(tile.Color))
		}
	}

	// Draw Enemies
	for _, e := range g.Enemies {
		style := tcell.StyleDefault.Foreground(e.Color)
		g.Screen.SetContent(e.X, e.Y, e.Visual, nil, style)
	}

	// Draw Player (on top of everything)
	g.Screen.SetContent(g.Player.X, g.Player.Y, g.Player.Visual, nil, tcell.StyleDefault.Foreground(g.Player.Color))

	// UI Overlay
	DrawText(g.Screen, 1, 0, " MODE: Exploring - Press 'M' for Menu ", tcell.StyleDefault.Reverse(true))
}

// renderMenu draws the tabbed interface
func (g *Game) renderMenu() {
	// Top Navigation Bar
	DrawText(g.Screen, 2, 2, " [ 1: Stats ]  [ 2: Inventory ]  [ 3: Skills ] (ESC to Close) ", tcell.StyleDefault.Reverse(true))

	switch g.ActiveTab {
	case 0:
		g.drawStats()
	case 1:
		DrawText(g.Screen, 2, 5, "-- Inventory Slot 1: Empty --", tcell.StyleDefault)
	case 2:
		g.drawSkills()
	}
}

// drawStats renders the detailed player information
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

// drawSkills renders the skill management screen
func (g *Game) drawSkills() {
	// Draw Active Slots
	DrawText(g.Screen, 2, 4, "--- ACTIVE SLOTS ---", tcell.StyleDefault.Foreground(tcell.ColorYellow))

	for i := 0; i < g.Player.MaxActiveSlots; i++ {
		row := 5 + i
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		prefix := fmt.Sprintf("[%d] ", i+1)

		// High contrast style when picking a slot
		if g.IsPickingSlot && i == g.SelectedSlotIndex {
			style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow).Bold(true)
			prefix = " TARGET -> "
		}

		// Check if slot has a skill
		if i < len(g.Player.ActiveSkills) && g.Player.ActiveSkills[i].Name != "Empty" {
			s := g.Player.ActiveSkills[i]
			DrawText(g.Screen, 2, row, prefix+s.Name, style)
		} else {
			if !(g.IsPickingSlot && i == g.SelectedSlotIndex) {
				style = style.Foreground(tcell.ColorGray)
			}
			DrawText(g.Screen, 2, row, prefix+"-- EMPTY --", style)
		}
	}

	// Draw Learned Skills
	listStyle := tcell.StyleDefault.Foreground(tcell.ColorLightCyan)
	if g.IsPickingSlot {
		listStyle = listStyle.Foreground(tcell.ColorGray) // Dim list while picking slot (was causing a problem where player couldnt see what he was selecting)
	}
	DrawText(g.Screen, 2, 11, "--- LEARNED SKILLS (Enter to Equip) ---", listStyle)

	for i, s := range g.Player.LearnedSkills {
		row := 12 + i
		style := tcell.StyleDefault
		prefix := "  "

		// Highlight selection in the learned list
		if !g.IsPickingSlot && i == g.SelectedSkillIndex {
			style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightCyan)
			prefix = "> "
		}
		DrawText(g.Screen, 2, row, prefix+s.Name, style)
	}
}

// Run initializes the game and manages the main loop
func Run() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	// Initialize Player first so we have stats/class ready
	player := entities.NewPlayer("rogue")

	// Initialize Game Instance (without the map yet)
	g := &Game{
		Screen:    screen,
		Player:    player,
		Enemies:   []entities.Enemy{},
		State:     "Exploring",
		ActiveTab: 0,
		Floor:     1, // Start on the first floor
	}

	// Generate Map (Now passing 'g' so it can add enemies to g.Enemies)
	m := NewMap(80, 45)
	m.GenerateLevel(g)
	g.Map = m

	// Position Player at center of first room
	if len(m.Rooms) > 0 {
		g.Player.X, g.Player.Y = m.Rooms[0].Center()
	}

	// Main Game Loop
	for {
		screen.Clear()

		switch g.State {
		case "Exploring":
			g.renderExploration()
		case "Menu":
			g.renderMenu()
		}

		screen.Show()

		// Handle Events
		ev := screen.PollEvent()
		if g.HandleInput(ev) {
			break
		}
	}
}
