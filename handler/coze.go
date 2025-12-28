package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-se/models"
	"github.com/hewo233/hdu-se/shared/consts"
	"io"
	"net/http"
)

type createConversationRequest struct {
	BotID string `json:"bot_id" binding:"required"`
	Name  string `json:"name" binding:"optional"`
}
type createConversationResponse struct {
	ConversationID string `json:"conversation_id"`
}

func CreateConversation(c *gin.Context) {
	type cozeAPIResponse struct {
		Code int `json:"code"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
	var req createConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40000,
			Result: "Invalid request parameters",
		})
		return
	}
	cozeReqBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50001,
			Result: "Failed to create request body",
		})
		return
	}

	client := &http.Client{}
	apiURL := consts.ConversationURL + "/create"

	proxyReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(cozeReqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	proxyReq.Header.Set("Authorization", "Bearer "+models.CozeToken)
	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call external api"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	var cozeResp cozeAPIResponse
	if err := json.Unmarshal(body, &cozeResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse external response"})
		return
	}

	if cozeResp.Code != 0 {
		c.JSON(http.StatusBadGateway, gin.H{"error": cozeResp.Msg, "code": cozeResp.Code})
		return
	}

	// 返回给客户端
	c.JSON(http.StatusOK, createConversationResponse{
		ConversationID: cozeResp.Data.ID,
	})
}
