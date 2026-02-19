package api

import (
	"database/sql"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"cadastral-service/internal/config"
	"cadastral-service/internal/models"
	"cadastral-service/internal/repository"
	"cadastral-service/internal/service"
)

type Handler struct {
	repo    *repository.Repository
	service *service.Service
	config  *config.Config
}

type QueryRequest struct {
	CadastralNumber string  `json:"cadastral_number" binding:"required"`
	Latitude        float64 `json:"latitude" binding:"required"`
	Longitude       float64 `json:"longitude" binding:"required"`
}

type QueryResponse struct {
	ID              string    `json:"id"`
	CadastralNumber string    `json:"cadastral_number"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Status          string    `json:"status"`
	Result          *bool     `json:"result,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     time.Time `json:"completed_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewHandler(db *sql.DB, cfg *config.Config) *Handler {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg)
	return &Handler{
		repo:    repo,
		service: svc,
		config:  cfg,
	}
}

// ping checking server
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"time":    time.Now().UTC(),
	})
}

// CreateQuery is creating new process
func (h *Handler) CreateQuery(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// data validataion
	if req.CadastralNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cadastral_number is required"})
		return
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude must be between -90 and 90"})
		return
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "longitude must be between -180 and 180"})
		return
	}

	// taken user ID from context if auth exist
	var userID string
	if claims, exists := c.Get("userClaims"); exists {
		if userClaims, ok := claims.(*Claims); ok {
			userID = userClaims.UserID
		}
	}

	// make request
	query := &models.Query{
		ID:              generateID(),
		CadastralNumber: req.CadastralNumber,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		Status:          "pending",
		UserID:          userID,
		CreatedAt:       time.Now(),
	}

	if err := h.repo.CreateQuery(c.Request.Context(), query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create query"})
		return
	}

	// run asynchronous processing
	go h.service.ProcessQuery(query)

	// return answer
	response := QueryResponse{
		ID:              query.ID,
		CadastralNumber: query.CadastralNumber,
		Latitude:        query.Latitude,
		Longitude:       query.Longitude,
		Status:          query.Status,
		CreatedAt:       query.CreatedAt,
	}

	c.JSON(http.StatusAccepted, response)
}

// GetHistory is take history of request
func (h *Handler) GetHistory(c *gin.Context) {
	ctx := c.Request.Context()

	// taken parametrs of pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// take user ID from context if auth exist
	var userID string
	if claims, exists := c.Get("userClaims"); exists {
		if userClaims, ok := claims.(*Claims); ok {
			userID = userClaims.UserID
		}
	}

	queries, err := h.repo.GetQueries(ctx, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queries"})
		return
	}

	// transform in response
	responses := make([]QueryResponse, len(queries))
	for i, query := range queries {
		responses[i] = QueryResponse{
			ID:              query.ID,
			CadastralNumber: query.CadastralNumber,
			Latitude:        query.Latitude,
			Longitude:       query.Longitude,
			Status:          query.Status,
			Result:          query.Result,
			CreatedAt:       query.CreatedAt,
			CompletedAt:     query.CompletedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// GetHistoryByCadastral need to take history by cadastral number
func (h *Handler) GetHistoryByCadastral(c *gin.Context) {
	ctx := c.Request.Context()
	cadastralNumber := c.Param("cadastral_number")

	if cadastralNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cadastral_number is required"})
		return
	}

	// take user ID from context if auth exist
	var userID string
	if claims, exists := c.Get("userClaims"); exists {
		if userClaims, ok := claims.(*Claims); ok {
			userID = userClaims.UserID
		}
	}

	queries, err := h.repo.GetQueriesByCadastral(ctx, cadastralNumber, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queries"})
		return
	}

	// transform in response
	responses := make([]QueryResponse, len(queries))
	for i, query := range queries {
		responses[i] = QueryResponse{
			ID:              query.ID,
			CadastralNumber: query.CadastralNumber,
			Latitude:        query.Latitude,
			Longitude:       query.Longitude,
			Status:          query.Status,
			Result:          query.Result,
			CreatedAt:       query.CreatedAt,
			CompletedAt:     query.CompletedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// ProcessResult its external sever emulation
func (h *Handler) ProcessResult(c *gin.Context) {
	// imitation of processing until 60 sec
	delay := time.Duration(rand.Intn(60)) * time.Second
	time.Sleep(delay)

	// random result
	result := rand.Intn(2) == 1

	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"delay":  delay.Seconds(),
	})
}

// Login its auth user
func (h *Handler) Login(c *gin.Context) {
	if !h.config.Auth.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authentication is disabled"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Пcheck a password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// generate a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.config.Auth.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// Register its register user
func (h *Handler) Register(c *gin.Context) {
	if !h.config.Auth.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "registration is disabled"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// create password
	user := &models.User{
		ID:           generateID(),
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := h.repo.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

// SwaggerHandler recoil Swagger docs
func (h *Handler) SwaggerHandler(c *gin.Context) {
	// here can recoil generate Swagger docs
	c.JSON(http.StatusOK, gin.H{
		"message": "Swagger documentation",
		// add here your Swagger specification
	})
}

// helpful function
func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
