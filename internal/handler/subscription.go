package handler

import (
	"context"
	"redis-subscribe-demo/internal/model"
	"redis-subscribe-demo/internal/service"
	"redis-subscribe-demo/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	subService *service.SubscriptionService
}

func NewHandler(subService *service.SubscriptionService) *Handler {
	return &Handler{
		subService: subService,
	}
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id" binding:"required"`
		Service string `json:"service" binding:"required"`
		Branch  string `json:"branch" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	err := h.subService.CreateSubscription(context.Background(), req.UserID, req.Service, req.Branch)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

//func (h *Handler) DeleteSubscription(c *gin.Context) {
//	var req struct {
//		UserID  string `json:"user_id" binding:"required"`
//		Service string `json:"service" binding:"required"`
//		Branch  string `json:"branch" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.Error(c, err)
//		return
//	}
//
//	err := h.subService.DeleteSubscription(req.UserID, req.Service, req.Branch)
//	if err != nil {
//		response.Error(c, err)
//		return
//	}
//
//	response.Success(c, nil)
//}

func (h *Handler) HandleNotification(c *gin.Context) {
	var notification model.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		response.Error(c, err)
		return
	}

	err := h.subService.HandleNotification(context.Background(), &notification)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

//func (h *Handler) ListSubscriptions(c *gin.Context) {
//	userID := c.Query("user_id")
//	if userID == "" {
//		response.Error(c, errors.New("user_id is required"))
//		return
//	}
//
//	subs, err := h.subService.ListSubscriptions(userID)
//	if err != nil {
//		response.Error(c, err)
//		return
//	}
//
//	response.Success(c, subs)
//}
