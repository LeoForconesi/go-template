package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	useruc "github.com/LeonardoForconesi/go-template/internal/usecase/user"
)

type UserHandlers struct {
	Create useruc.CreateUsecase
	Get    useruc.GetUsecase
	List   useruc.ListUsecase
	Update useruc.UpdateUsecase
	Delete useruc.DeleteUsecase
	Notify useruc.NotifyUsecase
}

func (h *UserHandlers) Register(rg *gin.RouterGroup) {
	rg.POST("/users", h.createUser)
	rg.GET("/users/:id", h.getUser)
	rg.GET("/users", h.listUsers)
	rg.PUT("/users/:id", h.updateUser)
	rg.DELETE("/users/:id", h.deleteUser)

	rg.POST("/notify", h.notifyUser)
}

type createReq struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
}

func (h *UserHandlers) createUser(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	out, err := h.Create.Execute(c.Request.Context(), useruc.CreateInput{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "resource already exists" || err.Error() == "user: email required" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *UserHandlers) getUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	out, err := h.Get.Execute(c.Request.Context(), useruc.GetInput{ID: id})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "resource not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *UserHandlers) listUsers(c *gin.Context) {
	// el paginado esta hardcodeado, aca podria recibir por parametros el page y size
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	out, err := h.List.Execute(c.Request.Context(), useruc.ListInput{Page: page, Size: size})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

type updateReq struct {
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
}

func (h *UserHandlers) updateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateReq
	if err := c.ShouldBindJSON(&req); err != nil || (req.Name == nil && req.Phone == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	out, err := h.Update.Execute(c.Request.Context(), useruc.UpdateInput{
		ID:    id,
		Name:  req.Name,
		Phone: req.Phone,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "resource not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *UserHandlers) deleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Delete.Execute(c.Request.Context(), useruc.DeleteInput{ID: id}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type notifyReq struct {
	UserID  string `json:"user_id" binding:"required,uuid"`
	Message string `json:"message" binding:"required"`
}

func (h *UserHandlers) notifyUser(c *gin.Context) {
	var req notifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	uid, _ := uuid.Parse(req.UserID)
	if err := h.Notify.Execute(c.Request.Context(), useruc.NotifyInput{
		UserID:  uid,
		Message: req.Message,
	}); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "resource not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
