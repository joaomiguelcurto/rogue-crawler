package engine

import (
	"math/rand"
	"rogue-crawler/entities"

	"github.com/gdamore/tcell/v2"
)

type Tile struct {
	Type   string // "Wall", "Floor", "Void"
	Visual rune
	Color  tcell.Color
}

type Rect struct {
	X1, Y1, X2, Y2 int
}

func (r Rect) Center() (int, int) {
	return (r.X1 + r.X2) / 2, (r.Y1 + r.Y2) / 2
}

type GameMap struct {
	Width, Height int
	Tiles         [][]Tile
	Rooms         []Rect
}

func NewMap(width, height int) *GameMap {
	m := &GameMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]Tile, width),
	}

	for x := 0; x < width; x++ {
		m.Tiles[x] = make([]Tile, height)
		for y := 0; y < height; y++ {
			// Change default to Void: Player cannot walk here
			m.Tiles[x][y] = Tile{Type: "Void", Visual: ' ', Color: tcell.ColorDefault}
		}
	}
	return m
}

func (m *GameMap) GenerateLevel(g *Game) {
	// Clear the enemies from previous floors
	g.Enemies = []entities.Enemy{}

	// Divide map into 3x3 grid, worked, now it avoids spawning on top of each other
	zoneW := m.Width / 3
	zoneH := m.Height / 3
	var roomCenters []struct{ x, y int }

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			// Skip random zones to create "empty" areas or "secret" paths (secret rooms not working yet)
			if rand.Intn(10) < 2 && len(roomCenters) > 1 {
				continue
			}

			// Define zone boundaries
			zX1, zY1 := col*zoneW, row*zoneH

			// Randomize room size within the zone (making them horizontally biased)
			w := rand.Intn(zoneW-10) + 8
			h := rand.Intn(zoneH-8) + 5

			// Randomize position within the zone
			x := zX1 + rand.Intn(zoneW-w-2) + 1
			y := zY1 + rand.Intn(zoneH-h-2) + 1

			newRoom := Rect{x, y, x + w, y + h}
			m.carveRoom(newRoom)
			m.Rooms = append(m.Rooms, newRoom)

			cx, cy := newRoom.Center()
			roomCenters = append(roomCenters, struct{ x, y int }{cx, cy})
		}
	}

	// Connect Rooms in a sequence to ensure every room has an entrance
	for i := 0; i < len(roomCenters)-1; i++ {
		start := roomCenters[i]
		end := roomCenters[i+1]

		// Connect rooms with L-shaped corridors
		if rand.Intn(2) == 0 {
			m.carveHorizontalTunnel(start.x, end.x, start.y)
			m.carveVerticalTunnel(start.y, end.y, end.x)
		} else {
			m.carveVerticalTunnel(start.y, end.y, start.x)
			m.carveHorizontalTunnel(start.x, end.x, end.y)
		}
	}

	for i, room := range m.Rooms {
		// Skip spawning in the starting room (index 0)
		if i == 0 {
			continue
		}

		// Spawn 1-3 enemies in this room
		m.spawnEnemiesInRoom(room, g)
	}

	// Place Portal in the last room
	lastRoom := m.Rooms[len(m.Rooms)-1]
	px, py := lastRoom.Center()
	m.Tiles[px][py] = Tile{Type: "Portal", Visual: 'O', Color: tcell.ColorPurple}
}

func (m *GameMap) spawnEnemiesInRoom(r Rect, g *Game) {
	// Choose number of enemies (1 to 2)
	numEnemies := rand.Intn(2) + 1

	for i := 0; i < numEnemies; i++ {
		// Find a random spot inside the room (not on the walls)
		ex := rand.Intn(r.X2-r.X1-1) + r.X1 + 1
		ey := rand.Intn(r.Y2-r.Y1-1) + r.Y1 + 1

		// Avoid spawning an enemy on top of the player or portal
		if m.Tiles[ex][ey].Type == "Floor" {
			// Pick a random enemy type
			enemyType := "goblin"
			if rand.Intn(2) == 0 {
				enemyType = "spider"
			}

			newEnemy := entities.NewEnemy(enemyType)
			newEnemy.X = ex
			newEnemy.Y = ey

			g.Enemies = append(g.Enemies, newEnemy)
		}
	}
}

func (m *GameMap) carveTunnelWithDeadEnds(pos1, pos2, constant int, horizontal bool) {
	for i := min(pos1, pos2); i <= max(pos1, pos2); i++ {
		x, y := i, constant
		if !horizontal {
			x, y = constant, i
		}

		m.Tiles[x][y] = Tile{Type: "Floor", Visual: ' ', Color: tcell.ColorDefault}
		m.ensureWall(x-1, y)
		m.ensureWall(x+1, y)
		m.ensureWall(x, y-1)
		m.ensureWall(x, y+1)

		// 10% chance to spawn a dead-end bait room
		if rand.Intn(100) < 10 {
			subRoomW, subRoomH := 5, 5
			subRoom := Rect{x, y, x + subRoomW, y + subRoomH}
			// Boundary check
			if subRoom.X2 < m.Width-1 && subRoom.Y2 < m.Height-1 {
				m.carveRoom(subRoom)
			}
		}
	}
}

func (m *GameMap) carveRoom(r Rect) {
	for x := r.X1; x <= r.X2; x++ {
		for y := r.Y1; y <= r.Y2; y++ {
			if x == r.X1 || x == r.X2 || y == r.Y1 || y == r.Y2 {
				// Only edges get the wall visual
				m.Tiles[x][y] = Tile{Type: "Wall", Visual: '#', Color: tcell.ColorGrey}
			} else {
				// Inner area is floor
				m.Tiles[x][y] = Tile{Type: "Floor", Visual: ' ', Color: tcell.ColorDefault}
			}
		}
	}
}

func (m *GameMap) carveHorizontalTunnel(x1, x2, y int) {
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		// Set floor
		m.Tiles[x][y] = Tile{Type: "Floor", Visual: ' ', Color: tcell.ColorDefault}
		// Add walls above and below the corridor if they are currently "Void"
		m.ensureWall(x, y-1)
		m.ensureWall(x, y+1)
	}
}

func (m *GameMap) carveVerticalTunnel(y1, y2, x int) {
	for y := min(y1, y2); y <= max(y1, y2); y++ {
		m.Tiles[x][y] = Tile{Type: "Floor", Visual: ' ', Color: tcell.ColorDefault}
		// Add walls to the left and right of the corridor
		m.ensureWall(x-1, y)
		m.ensureWall(x+1, y)
	}
}

func (m *GameMap) ensureWall(x, y int) {
	// Only place a wall if the tile is currently "Void" (empty space)
	// This prevents corridors from placing walls inside rooms
	if x >= 0 && x < m.Width && y >= 0 && y < m.Height {
		if m.Tiles[x][y].Type == "Void" {
			m.Tiles[x][y] = Tile{Type: "Wall", Visual: '#', Color: tcell.ColorGrey}
		}
	}
}
