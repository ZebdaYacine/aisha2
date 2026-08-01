package httpapi

import (
	"errors"

	"github.com/aisha-platform/aisha/backend/core/security"
	"github.com/aisha-platform/aisha/backend/features/auth"
	"github.com/aisha-platform/aisha/backend/features/users"
	"github.com/gofiber/fiber/v3"
)

type CustomerHandler struct {
	service   *customer.Service
	validator *RequestValidator
}

type profileUpdateRequest struct {
	DisplayName string `json:"displayName" validate:"required,min=2,max=120"`
	Phone       string `json:"phone" validate:"omitempty,min=6,max=32"`
}

type addressRequest struct {
	FullName   string `json:"fullName" validate:"required,min=2,max=120"`
	Phone      string `json:"phone" validate:"omitempty,min=6,max=32"`
	Line1      string `json:"line1" validate:"required,max=200"`
	Line2      string `json:"line2" validate:"max=200"`
	City       string `json:"city" validate:"required,max=100"`
	PostalCode string `json:"postalCode" validate:"required,max=32"`
	Country    string `json:"country" validate:"required,max=100"`
	IsDefault  bool   `json:"isDefault"`
}

type profileResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Phone       string `json:"phone"`
}

type addressResponse struct {
	ID         string `json:"id"`
	FullName   string `json:"fullName"`
	Phone      string `json:"phone"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	IsDefault  bool   `json:"isDefault"`
}

func NewCustomerHandler(service *customer.Service, validator *RequestValidator) *CustomerHandler {
	return &CustomerHandler{service: service, validator: validator}
}

func customerPrincipal(c fiber.Ctx) (auth.Principal, error) {
	p, ok := c.Locals(principalLocal).(auth.Principal)
	if !ok {
		return auth.Principal{}, NewAPIError(CodeAuthenticationRequired, "Authentication is required.", nil)
	}
	return p, nil
}

func (h *CustomerHandler) Profile(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	result, err := h.service.Profile(c.Context(), p)
	if err != nil {
		return customerAPIError(err)
	}
	return c.JSON(profileDTO(result))
}

func (h *CustomerHandler) UpdateProfile(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request profileUpdateRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	result, err := h.service.UpdateProfile(c.Context(), p, request.DisplayName, request.Phone)
	if err != nil {
		return customerAPIError(err)
	}
	return c.JSON(profileDTO(result))
}

func (h *CustomerHandler) Addresses(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	items, err := h.service.Addresses(c.Context(), p)
	if err != nil {
		return customerAPIError(err)
	}
	response := make([]addressResponse, len(items))
	for i, item := range items {
		response[i] = addressDTO(item)
	}
	return c.JSON(response)
}

func (h *CustomerHandler) CreateAddress(c fiber.Ctx) error { return h.writeAddress(c, "") }
func (h *CustomerHandler) UpdateAddress(c fiber.Ctx) error { return h.writeAddress(c, c.Params("id")) }

func (h *CustomerHandler) writeAddress(c fiber.Ctx, id string) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request addressRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	input := customer.AddressInput{FullName: request.FullName, Phone: request.Phone, Line1: request.Line1, Line2: request.Line2, City: request.City, PostalCode: request.PostalCode, Country: request.Country, Default: request.IsDefault}
	var result customer.Address
	if id == "" {
		result, err = h.service.CreateAddress(c.Context(), p, input)
	} else {
		result, err = h.service.UpdateAddress(c.Context(), p, id, input)
	}
	if err != nil {
		return customerAPIError(err)
	}
	status := fiber.StatusOK
	if id == "" {
		status = fiber.StatusCreated
	}
	return c.Status(status).JSON(addressDTO(result))
}

func (h *CustomerHandler) DeleteAddress(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteAddress(c.Context(), p, c.Params("id")); err != nil {
		return customerAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func profileDTO(p customer.Profile) profileResponse {
	return profileResponse{ID: p.ID, Email: p.Email, DisplayName: p.DisplayName, Phone: p.Phone}
}
func addressDTO(a customer.Address) addressResponse {
	return addressResponse{ID: a.ID, FullName: a.FullName, Phone: a.Phone, Line1: a.Line1, Line2: a.Line2, City: a.City, PostalCode: a.PostalCode, Country: a.Country, IsDefault: a.Default}
}
func customerAPIError(err error) error {
	switch {
	case errors.Is(err, customer.ErrValidation):
		return NewAPIError(CodeValidationError, "The request is invalid.", nil)
	case errors.Is(err, customer.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The requested resource was not found.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
