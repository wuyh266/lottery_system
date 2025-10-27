package models
import "time"

type Participant struct {
	ID int `json:"id"`
	Name string `json:"name"`
	JoinedAt time.Time `json:"joined_at"`
	IsWinner bool `json:"is_winner"`
	Prizes []Prize `json:"prizes"` // 参与者获得的奖品
}
type Lottery struct {
	Participants []Participant `json:"participants"`
	Winners []Participant `json:"winners"` // 所有中奖者
	IsDrawn bool `json:"is_drawn"`
	Prizes []Prize `json:"prizes"`
	UnclaimedPrizes []Prize `json:"unclaimed_prizes"` // 未分配的奖品
}

type Prize struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Quantity int `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}