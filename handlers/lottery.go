package handlers

import (
	"lottery_system/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

var store *storage.Storage

// Init 初始化存储
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

// GetPrizes 获取奖品列表
func GetPrizes(c *gin.Context) {
	prizes := store.GetPrizes()
	c.JSON(http.StatusOK, gin.H{
		"prizes": prizes,
		"count":  len(prizes),
	})
}
