package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrValidation        = errors.New("warehouse validation failed")
	ErrNotFound          = errors.New("warehouse resource not found")
	ErrDuplicate         = errors.New("warehouse reference already exists")
	ErrInvalidTransition = errors.New("warehouse resource is in an invalid state")
	ErrQuantityMismatch  = errors.New("inspection quantities do not match received quantity")
	ErrEvidenceRequired  = errors.New("inspection evidence is required")
	ErrAlreadyInspected  = errors.New("reception has already been inspected")
)

type ReceptionInput struct {
	ProductID        string
	ReceivedQuantity int64
	ReferenceKey     string
	SupplierName     string
	ParcelReference  string
	Notes            string
	ReceivedAt       *time.Time
}

type InspectionInput struct {
	AcceptedQuantity    int64
	RejectedQuantity    int64
	QuarantinedQuantity int64
	DamagedQuantity     int64
	Reason              string
}

type EvidenceInput struct {
	ObjectKey        string
	OriginalFilename string
	MediaType        string
	SizeBytes        int64
	Checksum         string
}

type Evidence struct {
	ID               string    `json:"id"`
	ReceptionID      string    `json:"receptionId"`
	ObjectKey        string    `json:"-"`
	OriginalFilename string    `json:"originalFilename"`
	MediaType        string    `json:"mediaType"`
	SizeBytes        int64     `json:"sizeBytes"`
	Checksum         string    `json:"checksum"`
	URL              string    `json:"url,omitempty"`
	UploadedBy       string    `json:"uploadedBy"`
	CreatedAt        time.Time `json:"createdAt"`
}

type Inspection struct {
	ID                  string    `json:"id"`
	ReceptionID         string    `json:"receptionId"`
	AcceptedQuantity    int64     `json:"acceptedQuantity"`
	RejectedQuantity    int64     `json:"rejectedQuantity"`
	QuarantinedQuantity int64     `json:"quarantinedQuantity"`
	DamagedQuantity     int64     `json:"damagedQuantity"`
	Reason              string    `json:"reason"`
	InspectedBy         string    `json:"inspectedBy"`
	InspectedAt         time.Time `json:"inspectedAt"`
}

type Reception struct {
	ID               string      `json:"id"`
	ProductID        string      `json:"productId"`
	ProductName      string      `json:"productName"`
	ArtisanProfileID string      `json:"artisanProfileId"`
	ArtisanName      string      `json:"artisanName"`
	WorkshopID       string      `json:"workshopId"`
	WorkshopName     string      `json:"workshopName"`
	SupplierName     string      `json:"supplierName"`
	ReceivedQuantity int64       `json:"receivedQuantity"`
	ReferenceKey     string      `json:"referenceKey"`
	ParcelReference  string      `json:"parcelReference"`
	Notes            string      `json:"notes"`
	Status           string      `json:"status"`
	ReceivedAt       time.Time   `json:"receivedAt"`
	ReceivedBy       string      `json:"receivedBy"`
	InspectedAt      *time.Time  `json:"inspectedAt,omitempty"`
	Inspection       *Inspection `json:"inspection,omitempty"`
	Evidence         []Evidence  `json:"evidence"`
}

type Repository interface {
	CreateReception(context.Context, string, ReceptionInput) (Reception, error)
	ListReceptions(context.Context, string, string, int, int) ([]Reception, int, error)
	Inspect(context.Context, string, string, InspectionInput) (Inspection, error)
	AddEvidence(context.Context, string, string, EvidenceInput) (Evidence, error)
	ListEvidence(context.Context, string) ([]Evidence, error)
}

type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
