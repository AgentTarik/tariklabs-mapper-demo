package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tariklabs/mapper"

	"go-test/go-test/internal/core/domain"
	"go-test/go-test/internal/core/service"
)

// GeoRequest uses mapconv for string-to-float conversion demonstration
type GeoRequest struct {
	Lat string `json:"lat" map:"Latitude" mapconv:"float64"`
	Lng string `json:"lng" map:"Longitude" mapconv:"float64"`
}

// AddressRequest with nested struct
type AddressRequest struct {
	Street  string     `json:"street"`
	Suite   string     `json:"suite"`
	City    string     `json:"city"`
	Zipcode string     `json:"zipcode"`
	Geo     GeoRequest `json:"geo"`
}

// CompanyRequest with map tags for field aliasing
type CompanyRequest struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catch_phrase" map:"CatchPhrase"`
	BS          string `json:"bs" map:"BS"`
}

// CreateUserRequest demonstrates combining map and mapconv tags
type CreateUserRequest struct {
	Name     string         `json:"name" binding:"required"`
	Username string         `json:"username"`
	Email    string         `json:"email" binding:"required"`
	Phone    string         `json:"phone"`
	Website  string         `json:"website"`
	Address  AddressRequest `json:"address"`
	Company  CompanyRequest `json:"company"`
}

// PatchUserRequest for partial updates - uses WithIgnoreZeroSource
type PatchUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

// GeoResponse for nested response
type GeoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AddressResponse with nested struct
type AddressResponse struct {
	Street  string      `json:"street"`
	Suite   string      `json:"suite"`
	City    string      `json:"city"`
	Zipcode string      `json:"zipcode"`
	Geo     GeoResponse `json:"geo"`
}

// CompanyResponse for nested response
type CompanyResponse struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catch_phrase" map:"CatchPhrase"`
	BS          string `json:"bs" map:"BS"`
}

// UserResponse demonstrates nested struct mapping
type UserResponse struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Website  string          `json:"website"`
	Address  AddressResponse `json:"address"`
	Company  CompanyResponse `json:"company"`
}

// ErrorResponse for mapper errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Field   string `json:"field,omitempty"`
	SrcType string `json:"src_type,omitempty"`
	DstType string `json:"dst_type,omitempty"`
}

type UserHandler interface {
	GetUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response UserResponse
	// Using MapWithOptions with multiple options:
	// - WithStrictMode: ensures all destination fields have matching source fields
	// - WithMaxDepth: sets maximum nesting depth for recursive mapping
	if err := mapper.MapWithOptions(&response, user, mapper.WithStrictMode(),
		mapper.WithMaxDepth(10),
	); err != nil {
		// Demonstrate MappingError handling
		var mappingErr *mapper.MappingError
		if errors.As(err, &mappingErr) {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "mapping failed",
				Field:   mappingErr.FieldPath,
				SrcType: mappingErr.SrcType,
				DstType: mappingErr.DstType,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to map response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user domain.User
	// Using MapWithOptions for request-to-domain mapping
	if err := mapper.MapWithOptions(&user, req,
		mapper.WithIgnoreZeroSource(),
	); err != nil {
		var mappingErr *mapper.MappingError
		if errors.As(err, &mappingErr) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid request mapping",
				Field:   mappingErr.FieldPath,
				SrcType: mappingErr.SrcType,
				DstType: mappingErr.DstType,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to map request"})
		return
	}

	createdUser, err := h.service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// First map the original request to response (preserves nested data like address, geo)
	var response UserResponse
	if err := mapper.Map(&response, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to map request to response"})
		return
	}

	// Then merge the API response using WithIgnoreZeroSource
	// This only updates non-zero fields (like ID) from the created user
	// while preserving the nested data from the original request
	if err := mapper.MapWithOptions(&response, createdUser,
		mapper.WithIgnoreZeroSource(),
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to merge response"})
		return
	}

	c.JSON(http.StatusCreated, response)
}
