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
	apiURL := consts.CreateConversationURL

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

type retrieveConversationRequest struct {
	ConversationID string `form:"conversation_id" binding:"required"`
	ChatID         string `form:"chat_id" binding:"required"`
}

type retrieveConversationResponse struct {
	Status string `json:"status"`
}

func RetrieveConversation(c *gin.Context) {
	var req retrieveConversationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40000,
			Result: "Invalid request parameters",
		})
		return
	}

	client := &http.Client{}
	apiURL := consts.RetrieveConversationURL + "?conversation_id=" + req.ConversationID + "&chat_id=" + req.ChatID

	proxyReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50001,
			Result: "Failed to create request",
		})
		return
	}

	proxyReq.Header.Set("Authorization", "Bearer "+models.CozeToken)
	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50002,
			Result: "Failed to call external API",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50003,
			Result: "Failed to read response",
		})
		return
	}

	/*
		{"code":0,"data":{"bot_id":"7563218003241058343","completed_at":1766909715,"conversation_id":"7588818179242721321","created_at":1766909711,"id":"7588819419779039272","status":"completed","usage":{"input_count":972,"input_tokens_details":{"cached_tokens":0},"output_count":180,"output_tokens_details":{"reasoning_tokens":0},"token_count":1152}},"detail":{"logid":"2025122816153901654CB1627CB59025E7"},"msg":""}
	*/

	type cozeAPIResponse struct {
		Code int `json:"code"`
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	var cozeResp cozeAPIResponse
	if err := json.Unmarshal(body, &cozeResp); err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50004,
			Result: "Failed to parse external response",
		})
		return
	}
	if cozeResp.Code != 0 {
		c.JSON(http.StatusBadGateway, models.Report{
			Code:   cozeResp.Code,
			Result: cozeResp.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, retrieveConversationResponse{
		Status: cozeResp.Data.Status,
	})
}
