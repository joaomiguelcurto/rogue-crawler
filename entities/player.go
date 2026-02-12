package entities

import (
	"github.com/gdamore/tcell/v2"
)

type Player struct {
	X, Y         int
	Visual       rune
	Color        tcell.Color
	Class        string
	Level        int
	XP           int
	HP, MaxHP    int
	MP, MaxMP    int
	Strength     int
	Intelligence int
	Agility      int
	Luck         int
	// Equipment    []Equipment
	// Inventory    []Item
	MaxActiveSlots int
	LearnedSkills  []Skill
	ActiveSkills   []Skill
}

func NewPlayer(class string) *Player {
	// Shared defaults
	p := &Player{
		X:              10,
		Y:              10,
		Visual:         '@',
		Color:          tcell.ColorWhite,
		Level:          1,
		XP:             0,
		MaxActiveSlots: 4,
		LearnedSkills:  []Skill{},
		ActiveSkills:   []Skill{},
	}

	// Class-specific stat modifiers
	switch class {
	case "warrior":
		p.Class = "Warrior"
		p.HP, p.MaxHP = 30, 30
		p.MP, p.MaxMP = 5, 5
		p.Strength = 8
		p.Intelligence = 2
		p.Agility = 4
		p.Luck = 3
		p.GiveSkill("Bloodlust")

	case "mage":
		p.Class = "Mage"
		p.HP, p.MaxHP = 15, 15
		p.MP, p.MaxMP = 25, 25
		p.Strength = 4
		p.Intelligence = 8
		p.Agility = 3
		p.Luck = 6
		p.GiveSkill("Arcane Blast")

	case "rogue":
		p.Class = "Rogue"
		p.HP, p.MaxHP = 20, 20
		p.MP, p.MaxMP = 10, 10
		p.Strength = 3
		p.Intelligence = 4
		p.Agility = 8
		p.Luck = 7
		p.GiveSkill("Quick Slash")

	// Default "Novice" stats if class is unknown
	default:
		p.Class = "Novice"
		p.HP, p.MaxHP = 20, 20
		p.MP, p.MaxMP = 10, 10
		p.Strength = 3
		p.Intelligence = 3
		p.Agility = 3
		p.Luck = 3
	}

	return p
}

func (p *Player) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

func (p *Player) EquipSkill(skillIndex int, slotIndex int) {
	if skillIndex < 0 || skillIndex >= len(p.LearnedSkills) {
		return
	}
	if slotIndex < 0 || slotIndex >= p.MaxActiveSlots {
		return
	}

	newSkill := p.LearnedSkills[skillIndex]

	// Remove the skill if it exists in ANY other slot
	for i, s := range p.ActiveSkills {
		if s.Name == newSkill.Name {
			// Remove it by slice manipulation (weird stuff from stackoverflow)
			p.ActiveSkills = append(p.ActiveSkills[:i], p.ActiveSkills[i+1:]...)
		}
	}

	// Makes sure ActiveSkills slice is long enough to reach the slotIndex
	// This fills the gap if you try to equip slot 3 when only slot 0 is full
	for len(p.ActiveSkills) <= slotIndex {
		p.ActiveSkills = append(p.ActiveSkills, Skill{Name: "Empty"})
	}

	// Place the skill in the specific slot
	p.ActiveSkills[slotIndex] = newSkill
}
