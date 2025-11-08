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
	Active   HabitStatus = "active"
	Archived HabitStatus = "archived"
	Any      HabitStatus = "any"
)

var HabitStatusMapping = map[string]HabitStatus{
	string(Active):   Active,
	string(Archived): Archived,
	string(Any):      Any,
}

type Color string

const (
	Red    Color = "red"
	Orange Color = "orange"
	Yellow Color = "yellow"
	Lime   Color = "lime"
	Green  Color = "green"
	Blue   Color = "blue"
	Purple Color = "purple"
)

var ColorMapping = map[string]Color{
	string(Red):    Red,
	string(Orange): Orange,
	string(Yellow): Yellow,
	string(Lime):   Lime,
	string(Green):  Green,
	string(Blue):   Blue,
	string(Purple): Purple,
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
