package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GatewayService struct {
	hotelInfoURL   string
	hotelSearchURL string
	userBookingURL string
	client         *http.Client
}

func NewGatewayService() *GatewayService {
	return &GatewayService{
		hotelInfoURL:   getEnv("HOTEL_INFO_URL", "http://localhost:8081"),
		hotelSearchURL: getEnv("HOTEL_SEARCH_URL", "http://localhost:8082"),
		userBookingURL: getEnv("USER_BOOKING_URL", "http://localhost:8083"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Auth handlers
func (gs *GatewayService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := gs.forwardRequest("POST", gs.userBookingURL+"/api/auth/login", req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := gs.forwardRequest("POST", gs.userBookingURL+"/api/auth/register", req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

// Hotel handlers
func (gs *GatewayService) SearchHotels(c *gin.Context) {
	params := c.Request.URL.Query()

	url := fmt.Sprintf("%s/api/hotels/search?%s", gs.hotelSearchURL, params.Encode())
	resp, err := gs.forwardRequest("GET", url, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) GetHotel(c *gin.Context) {
	hotelID := c.Param("id")

	url := fmt.Sprintf("%s/api/hotels/%s", gs.hotelInfoURL, hotelID)
	resp, err := gs.forwardRequest("GET", url, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hotel service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) CreateHotel(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("POST", gs.hotelInfoURL+"/api/hotels", req, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) UpdateHotel(c *gin.Context) {
	hotelID := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	url := fmt.Sprintf("%s/api/hotels/%s", gs.hotelInfoURL, hotelID)
	resp, err := gs.forwardRequest("PUT", url, req, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) DeleteHotel(c *gin.Context) {
	hotelID := c.Param("id")

	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	url := fmt.Sprintf("%s/api/hotels/%s", gs.hotelInfoURL, hotelID)
	resp, err := gs.forwardRequest("DELETE", url, nil, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) CheckAvailability(c *gin.Context) {
	hotelID := c.Param("id")
	params := c.Request.URL.Query()

	url := fmt.Sprintf("%s/api/availability/%s?%s", gs.userBookingURL, hotelID, params.Encode())
	resp, err := gs.forwardRequest("GET", url, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Availability service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

// Booking handlers
func (gs *GatewayService) CreateBooking(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("POST", gs.userBookingURL+"/api/bookings", req, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) GetUserBookings(c *gin.Context) {
	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("GET", gs.userBookingURL+"/api/bookings/user", nil, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) GetAllBookings(c *gin.Context) {
	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("GET", gs.userBookingURL+"/api/bookings", nil, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) UpdateBookingStatus(c *gin.Context) {
	bookingID := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	url := fmt.Sprintf("%s/api/bookings/%s/status", gs.userBookingURL, bookingID)
	resp, err := gs.forwardRequest("PATCH", url, req, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

// User handlers
func (gs *GatewayService) GetProfile(c *gin.Context) {
	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("GET", gs.userBookingURL+"/api/users/profile", nil, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) UpdateProfile(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("PUT", gs.userBookingURL+"/api/users/profile", req, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

func (gs *GatewayService) GetAllUsers(c *gin.Context) {
	headers := map[string]string{
		"Authorization": c.GetHeader("Authorization"),
	}

	resp, err := gs.forwardRequest("GET", gs.userBookingURL+"/api/users", nil, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User service unavailable"})
		return
	}

	c.JSON(resp.StatusCode, resp.Data)
}

// Helper methods
type ServiceResponse struct {
	StatusCode int
	Data       interface{}
}

func (gs *GatewayService) forwardRequest(method, url string, body interface{}, headers map[string]string) (*ServiceResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if len(responseBody) > 0 {
		json.Unmarshal(responseBody, &data)
	}

	return &ServiceResponse{
		StatusCode: resp.StatusCode,
		Data:       data,
	}, nil
}

// Request/Response types
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}
