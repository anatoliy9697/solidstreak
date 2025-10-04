package habit

import (
	"time"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
)

type Habit struct {
	ID          int64         `json:"id"`
	Active      bool          `json:"active"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatorID   int64         `json:"creatorId"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Checks      []*HabitCheck `json:"checks,omitempty"`
}

type HabitCheck struct {
	HabitID   int64     `json:"-"`
	UserID    int64     `json:"-"`
	CheckDate date.Date `json:"checkDate"`
	Completed bool      `json:"completed"`
	CheckedAt time.Time `json:"checkedAt"`
}

func NewHabit(title, description string, creatorID int64) *Habit {
	return &Habit{
		Active:      true,
		Title:       title,
		Description: description,
		CreatorID:   creatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
