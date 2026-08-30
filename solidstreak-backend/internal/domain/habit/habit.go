package habit

import (
	"time"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
)

type Habit struct {
	ID          int64         `json:"id"`
	Active      bool          `json:"active"`
	Archived    bool          `json:"archived"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Color       Color         `json:"color"`
	CreatorID   int64         `json:"creatorId"`
	IsPublic    bool          `json:"isPublic"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Checks      []*HabitCheck `json:"checks"`
}

type HabitStatus string

const (
	HabitStatusActive   HabitStatus = "active"
	HabitStatusArchived HabitStatus = "archived"
	HabitStatusAny      HabitStatus = "any"
)

var HabitStatusMapping = map[string]HabitStatus{
	string(HabitStatusActive):   HabitStatusActive,
	string(HabitStatusArchived): HabitStatusArchived,
	string(HabitStatusAny):      HabitStatusAny,
}

type Color string

const (
	ColorRed    Color = "red"
	ColorOrange Color = "orange"
	ColorYellow Color = "yellow"
	ColorLime   Color = "lime"
	ColorGreen  Color = "green"
	ColorBlue   Color = "blue"
	ColorPurple Color = "purple"
)

var ColorMapping = map[string]Color{
	string(ColorRed):    ColorRed,
	string(ColorOrange): ColorOrange,
	string(ColorYellow): ColorYellow,
	string(ColorLime):   ColorLime,
	string(ColorGreen):  ColorGreen,
	string(ColorBlue):   ColorBlue,
	string(ColorPurple): ColorPurple,
}

type HabitCheck struct {
	HabitID   int64     `json:"-"`
	UserID    int64     `json:"-"`
	CheckDate date.Date `json:"checkDate"`
	Completed bool      `json:"completed"`
	CheckedAt time.Time `json:"checkedAt"`
}

func NewHabit(title, description string, color Color, creatorID int64, isPublic bool) *Habit {
	return &Habit{
		Active:      true,
		Archived:    false,
		Title:       title,
		Description: description,
		Color:       color,
		IsPublic:    isPublic,
		CreatorID:   creatorID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}
