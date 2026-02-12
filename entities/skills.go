package entities

import "fmt"

type Skill struct {
	Name        string
	Description string
	DamageType  string // "Physical", "Magic", "True"
	Cost        int
	CostType    string // "MP", "HP", "None"
	MinDamage   int
	MaxDamage   int
	StatType    string // "Strength", "Intelligence", "Agility"
}

// SkillDatabase is the Master List
var SkillDatabase = map[string]Skill{
	"Bloodlust": {
		Name: "Bloodlust", Cost: 3, CostType: "HP",
		MinDamage: 7, MaxDamage: 9, StatType: "Strength",
	},
	"Fireball": {
		Name: "Arcane Blast", Cost: 4, CostType: "MP",
		MinDamage: 2, MaxDamage: 5, StatType: "Intelligence",
	},
	"Quick Slash": {
		Name: "Quick Slash", Cost: 0, CostType: "None",
		MinDamage: 1, MaxDamage: 2, StatType: "Agility",
	},
}

// GetSkill searches the database for a skill by name
func GetSkill(name string) (Skill, bool) {
	skill, exists := SkillDatabase[name]
	return skill, exists
}

// GetDamage returns the scaling damage based on player stats
func (s Skill) GetDamage(p *Player) (int, int) {
	scaling := 0
	switch s.StatType {
	case "Strength":
		scaling = p.Strength
	case "Intelligence":
		scaling = p.Intelligence
	case "Agility":
		scaling = p.Agility
	}

	// Scaling formula: Stat * 1.5 added to base
	min := s.MinDamage + int(float64(scaling)*1.5)
	max := s.MaxDamage + int(float64(scaling)*1.5)
	return min, max
}

// GetCostString returns a formatted string like "5 MP" or "Free"
func (s Skill) GetCostString() string {
	if s.CostType == "None" || s.Cost == 0 {
		return "Free"
	}
	return fmt.Sprintf("%d %s", s.Cost, s.CostType)
}

func (p *Player) GiveSkill(name string) {
	if s, ok := GetSkill(name); ok {
		p.LearnedSkills = append(p.LearnedSkills, s)
	}
}
