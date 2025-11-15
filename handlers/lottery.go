package handlers

import (
	"lottery_system/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var store *storage.Storage

func Init(s *storage.Storage) {
	store = s
}
func AddParticipant(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供姓名"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "姓名不能为空"})
		return
	}
	if err := store.AddParticipant(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加参与者失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "参与者添加成功"})
}
func GetParticipants(c *gin.Context) {
	participants := store.GetParticipants()
	c.JSON(http.StatusOK, gin.H{
		"participants": participants,
		"count":        len(participants),
	})
}
func DrawWinner(c *gin.Context) {
	winners, err := store.DrawWinner()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "抽奖成功",
		"winners": winners,
	})
}

func GetWinner(c *gin.Context) {
	lottery := store.GetLottery()
	c.JSON(http.StatusOK, gin.H{
		"is_drawn":     lottery.IsDrawn,
		"participants": len(lottery.Participants),
		"winners":      lottery.Winners,
	})
}

func Reset(c *gin.Context) {
	if err := store.Reset(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置抽奖失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "抽奖已重置"})
}
func AddPrize(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供奖品名称"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "奖品名称不能为空"})
		return
	}

	if err := store.AddPrize(req.Name, req.Description, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加奖品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "奖品添加成功"})
}

func GetPrizes(c *gin.Context) {
	prizes := store.GetPrizes()
	c.JSON(http.StatusOK, gin.H{
		"prizes": prizes,
		"count":  len(prizes),
	})
}

// SetDrawTime 设置开奖时间
func SetDrawTime(c *gin.Context) {
	var req struct {
		DrawTime string `json:"draw_time" binding:"required"` // 格式: "2006-01-02T15:04:05Z07:00" 或 RFC3339格式
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供开奖时间"})
		return
	}
	drawTime, err := time.Parse(time.RFC3339, req.DrawTime)
	if err != nil {
		// 尝试其他常见格式
		drawTime, err = time.Parse("2006-01-02 15:04:05", req.DrawTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误，请使用 RFC3339 格式 (例如: 2024-01-01T12:00:00+08:00)"})
			return
		}
	}
	if err := store.SetDrawTime(drawTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "开奖时间设置成功",
		"draw_time": drawTime.Format(time.RFC3339),
	})
}

// GetDrawTime 获取开奖时间
func GetDrawTime(c *gin.Context) {
	drawTime := store.GetDrawTime()
	if drawTime == nil {
		c.JSON(http.StatusOK, gin.H{
			"draw_time": nil,
			"message":   "未设置开奖时间",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"draw_time": drawTime.Format(time.RFC3339),
		"is_passed": time.Now().After(*drawTime),
	})
}
