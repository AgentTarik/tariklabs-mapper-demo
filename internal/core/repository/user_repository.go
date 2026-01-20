package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tariklabs/mapper"

	"go-test/go-test/internal/core/domain"
	"go-test/go-test/internal/engine"
)

const baseURL = "https://jsonplaceholder.typicode.com"

// GeoDTO uses map tag for field aliasing (lat -> Latitude, lng -> Longitude)
// and mapconv tag for string-to-float64 conversion
type GeoDTO struct {
	Lat string `json:"lat" map:"Latitude" mapconv:"float64"`
	Lng string `json:"lng" map:"Longitude" mapconv:"float64"`
}

// AddressDTO with nested struct to demonstrate recursive mapping
type AddressDTO struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     GeoDTO `json:"geo"`
}

// CompanyDTO uses map tag to alias catchPhrase -> CatchPhrase and bs -> BS
type CompanyDTO struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase" map:"CatchPhrase"`
	BS          string `json:"bs" map:"BS"`
}

// GetUserResponseDTO demonstrates nested structs and field aliasing
type GetUserResponseDTO struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Website  string     `json:"website"`
	Address  AddressDTO `json:"address"`
	Company  CompanyDTO `json:"company"`
}

// CreateUserRequestDTO for outgoing requests
type CreateUserRequestDTO struct {
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Website  string     `json:"website"`
	Address  AddressDTO `json:"address"`
	Company  CompanyDTO `json:"company"`
}

// CreateUserResponseDTO for parsing create response
type CreateUserResponseDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRepository interface {
	GetByID(id int) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
}

type userRepository struct {
	engine engine.HTTPEngine
}

func NewUserRepository(engine engine.HTTPEngine) UserRepository {
	return &userRepository{
		engine: engine,
	}
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	url := fmt.Sprintf("%s/users/%d", baseURL, id)

	respBody, statusCode, err := r.engine.MakeRequest(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if statusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var dto GetUserResponseDTO
	err = json.Unmarshal(respBody, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var user domain.User
	// Using MapWithOptions with WithStrictMode to ensure all fields are mapped
	// This demonstrates nested struct mapping with field aliasing and type conversion
	if err := mapper.MapWithOptions(&user, dto, mapper.WithStrictMode()); err != nil {
		return nil, fmt.Errorf("failed to map response to domain: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	url := fmt.Sprintf("%s/users", baseURL)

	var requestDTO CreateUserRequestDTO
	// Using MapWithOptions with WithIgnoreZeroSource to skip empty fields
	// This is useful for partial updates where we don't want to overwrite with zero values
	if err := mapper.MapWithOptions(&requestDTO, user, mapper.WithIgnoreZeroSource()); err != nil {
		return nil, fmt.Errorf("failed to map domain to request: %w", err)
	}

	body, err := json.Marshal(requestDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	respBody, statusCode, err := r.engine.MakeRequest(url, http.MethodPost, body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if statusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var responseDTO CreateUserResponseDTO
	if err := json.Unmarshal(respBody, &responseDTO); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Using basic Map for simple struct mapping
	var createdUser domain.User
	if err := mapper.Map(&createdUser, responseDTO); err != nil {
		return nil, fmt.Errorf("failed to map response to domain: %w", err)
	}

	return &createdUser, nil
}
