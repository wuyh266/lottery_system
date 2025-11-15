package storage

import (
	"encoding/json"
	"fmt"
	"lottery_system/models"
	"math/rand"
	"os"
	"time"
)

const dataFile = "data/lottery.json"

type Storage struct {
	lottery *models.Lottery
}

func NewStorage() *Storage {
	return &Storage{
		lottery: &models.Lottery{
			Participants:    make([]models.Participant, 0),
			Winners:         make([]models.Participant, 0),
			IsDrawn:         false,
			Prizes:          make([]models.Prize, 0),
			UnclaimedPrizes: make([]models.Prize, 0),
			DrawTime:        nil,
		},
	}
}
func (s *Storage) Load() error {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s.lottery)
}
func (s *Storage) Save() error {
	if _err := os.MkdirAll("data", 0755); _err != nil {
		return _err
	}
	data, err := json.MarshalIndent(s.lottery, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func (s *Storage) AddParticipant(name string) error {
	if s.lottery.IsDrawn {
		return fmt.Errorf("cannot add participant after lottery is drawn")
	}
	newID := len(s.lottery.Participants) + 1
	participant := models.Participant{
		ID:       newID,
		Name:     name,
		JoinedAt: time.Now(),
		IsWinner: false,
		Prizes:   make([]models.Prize, 0),
	}
	s.lottery.Participants = append(s.lottery.Participants, participant)
	return s.Save()
}
func (s *Storage) GetParticipants() []models.Participant {
	return s.lottery.Participants
}
func (s *Storage) GetLottery() *models.Lottery {
	return s.lottery
}

func (s *Storage) DrawWinner() ([]models.Participant, error) {
	if s.lottery.IsDrawn {
		return nil, fmt.Errorf("lottery already drawn")
	}
	if len(s.lottery.Participants) == 0 {
		return nil, fmt.Errorf("no participants to draw from")
	}
	if len(s.lottery.Prizes) == 0 {
		return nil, fmt.Errorf("no prizes available")
	}
	rand.Seed(time.Now().UnixNano())
	participants := make([]models.Participant, len(s.lottery.Participants))
	copy(participants, s.lottery.Participants)
	for i := len(participants) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		participants[i], participants[j] = participants[j], participants[i]
	}
	totalPrizes := 0
	for _, prize := range s.lottery.Prizes {
		totalPrizes += prize.Quantity
	}

	numParticipants := len(participants)
	winners := make([]models.Participant, 0)
	availablePrizes := make([]models.Prize, len(s.lottery.Prizes))
	copy(availablePrizes, s.lottery.Prizes)

	if numParticipants == 1 {
		participants[0].IsWinner = true
		participants[0].Prizes = availablePrizes
		winners = append(winners, participants[0])
	} else {
		prizeList := make([]models.Prize, 0)
		for _, prize := range availablePrizes {
			for q := 0; q < prize.Quantity; q++ {
				prizeList = append(prizeList, prize)
			}
		}

		for i := len(prizeList) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			prizeList[i], prizeList[j] = prizeList[j], prizeList[i]
		}

		// 给参与者分配奖品
		for i := 0; i < numParticipants; i++ {
			participants[i].IsWinner = true
			if i < len(prizeList) {
				participants[i].Prizes = append(participants[i].Prizes, prizeList[i])
			}
			winners = append(winners, participants[i])
		}
	}

	s.lottery.Winners = winners
	s.lottery.Participants = participants
	s.lottery.IsDrawn = true
	s.Save()
	return winners, nil
}
func (s *Storage) Reset() error {
	s.lottery = &models.Lottery{
		Participants:    make([]models.Participant, 0),
		Winners:         make([]models.Participant, 0),
		Prizes:          make([]models.Prize, 0),
		UnclaimedPrizes: make([]models.Prize, 0),
		IsDrawn:         false,
		DrawTime:        nil,
	}
	return s.Save()
}
func (s *Storage) AddPrize(name, description string, quantity int) error {
	newID := len(s.lottery.Prizes) + 1
	prize := models.Prize{
		ID:          newID,
		Name:        name,
		Description: description,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
	}
	s.lottery.Prizes = append(s.lottery.Prizes, prize)
	return s.Save()
}
func (s *Storage) GetPrizes() []models.Prize {
	return s.lottery.Prizes
}
func (s *Storage) SetDrawTime(drawTime time.Time) error {
	if s.lottery.IsDrawn {
		return fmt.Errorf("cannot set draw time after lottery is drawn")
	}
	s.lottery.DrawTime = &drawTime
	return s.Save()
}


func (s *Storage) GetDrawTime() *time.Time {
	return s.lottery.DrawTime
}

func (s *Storage) CheckAndAutoDraw() error {
	if s.lottery.IsDrawn {
		return nil 
	}
	if s.lottery.DrawTime == nil {
		return nil 
	}
	now := time.Now()
	if now.After(*s.lottery.DrawTime) || now.Equal(*s.lottery.DrawTime) {
		
		_, err := s.DrawWinner()
		return err
	}
	return nil
}
