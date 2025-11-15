package models

import "time"

type Participant struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	JoinedAt time.Time `json:"joined_at"`
	IsWinner bool      `json:"is_winner"`
	Prizes   []Prize   `json:"prizes"`
}
type Lottery struct {
	Participants    []Participant `json:"participants"`
	Winners         []Participant `json:"winners"`
	IsDrawn         bool          `json:"is_drawn"`
	Prizes          []Prize       `json:"prizes"`
	UnclaimedPrizes []Prize       `json:"unclaimed_prizes"`
	DrawTime        *time.Time    `json:"draw_time,omitempty"` 
}

type Prize struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
}
