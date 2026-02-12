package entities

import (
	"github.com/gdamore/tcell/v2"
)

type Enemy struct {
	Name   string
	Visual rune
	Color  tcell.Color
	HP     int
	Power  int
}

func NewEnemy(enemyType string) Enemy {
	switch enemyType {
	case "goblin":
		return Enemy{
			Name:   "Goblin",
			Visual: 'g',
			Color:  tcell.ColorGreen,
			HP:     10,
			Power:  5,
		}
	case "spider":
		return Enemy{
			Name:   "Spider",
			Visual: 's',
			Color:  tcell.ColorGray,
			HP:     5,
			Power:  3,
		}
	default:
		return Enemy{
			Name:   "Unknown",
			Visual: '?',
			Color:  tcell.ColorWhite,
			HP:     1,
			Power:  1,
		}
	}
}
