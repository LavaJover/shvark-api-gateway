package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	ProfileClient *client.ProfileClient
}

func NewProfileHandler(addr string) (*ProfileHandler, error) {
	profileClient, err := client.NewProfileClient(addr)
	if err != nil {
		return nil, err
	}
	return &ProfileHandler{
		ProfileClient: profileClient,
	}, nil
}

// @Summary Get profile by uuid
// @Description Get profile by uuid
// @Tags profiles
// @Accept json
// @Produce json
// @Param uuid path string true "Profile uuid"
// @Success 200 {object} GetProfileByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /profiles/{uuid} [get]
func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	profileID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.ProfileClient.GetProfileByID(profileID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile_id": response.ProfileId,
		"avatar_url": response.AvatarUrl,
		"tg_link": response.TgLink,
		"user_id": response.UserId,
	})
}

type GetProfileByIDResponse struct {
	ProfileID string `json:"profile_id"`
	AvatarURL string `json:"avatar_url"`
	TgLink 	  string `json:"tg_link"`
	UserID    string `json:"user_id"`
}