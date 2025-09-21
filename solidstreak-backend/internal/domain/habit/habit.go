package habit

import "time"

type Habit struct {
	ID          int64     `json:"id"`
	Active      bool      `json:"active"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creatorId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
