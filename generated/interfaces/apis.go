package interfaces

import (
	"context"
	"encoding/json"
)

// SessionInfo represents a normalized session object across freeze versions
type SessionInfo struct {
	SessionID string `json:"sessionID,omitempty"`
	Status    string `json:"status,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// PlanSchedule describes a plan execution schedule
type PlanSchedule struct {
	StartAt   string `json:"startAt,omitempty"`
	EndAt     string `json:"endAt,omitempty"`
	Frequency string `json:"frequency,omitempty"`
}

// PlanPayload represents the core structured payload for a plan
type PlanPayload struct {
	Type       string                 `json:"type,omitempty"`
	Schedule   PlanSchedule           `json:"schedule,omitempty"`
	Targets    []string               `json:"targets,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Raw        json.RawMessage        `json:"raw,omitempty"`
}

// PlanDefinition represents a normalized plan payload across versions
type PlanDefinition struct {
	PlanID      string      `json:"planID,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Payload     PlanPayload `json:"payload,omitempty"`
}

// CollectionItem represents a single piece of collected plan data
type CollectionItem struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// CollectionData represents normalized collected data for a plan
type CollectionData struct {
	PlanID  string           `json:"planID,omitempty"`
	Summary string           `json:"summary,omitempty"`
	Items   []CollectionItem `json:"items,omitempty"`
	Raw     json.RawMessage  `json:"raw,omitempty"`
}

// CollectionStatus represents the current status of a collection plan
type CollectionStatus struct {
	PlanID      string                 `json:"planID,omitempty"`
	State       string                 `json:"state,omitempty"`
	LastUpdated string                 `json:"lastUpdated,omitempty"`
	Progress    float64                `json:"progress,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// MetadataAPI defines common metadata operations available in both freeze versions
type MetadataAPI interface {
	// Equipment Metadata operations
	GetEquipmentStructure(ctx context.Context) error
	GetEquipmentNodeDescriptions(ctx context.Context) error

	// Type operations
	GetTypeDefinitions(ctx context.Context) error
	GetStateMachines(ctx context.Context) error
	GetSEMIObjTypes(ctx context.Context) error

	// Utility methods
	GetUnits(ctx context.Context) error
	GetExceptions(ctx context.Context) error
}

// SecurityAPI defines common security/authentication operations
type SecurityAPI interface {
	// Authentication and session management
	Login(ctx context.Context, username, password string) (sessionID string, err error)
	Logout(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, sessionID string) (valid bool, err error)

	// Permission checks
	CheckPermission(ctx context.Context, sessionID string, resource string) (allowed bool, err error)
}

// SessionAPI defines common session management operations
type SessionAPI interface {
	// Session lifecycle
	CreateSession(ctx context.Context) (sessionID string, err error)
	DestroySession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (sessionInfo SessionInfo, err error)

	// Session operations
	GetActiveSessions(ctx context.Context) (sessions []string, err error)
	RefreshSession(ctx context.Context, sessionID string) error
}

// DataCollectionAPI defines common data collection operations
type DataCollectionAPI interface {
	// Plan management
	DefinePlan(ctx context.Context, planDef PlanDefinition) error
	GetDefinedPlanIds(ctx context.Context) (ids []string, err error)
	GetPlanDefinition(ctx context.Context, planID string) (planDef PlanDefinition, err error)
	DeletePlan(ctx context.Context, planID string) error

	// Plan activation
	ActivatePlan(ctx context.Context, planID string) error
	GetActivePlanIds(ctx context.Context) (ids []string, err error)
	DeactivatePlan(ctx context.Context, planID string) error

	// Data collection
	CollectData(ctx context.Context, planID string) (data CollectionData, err error)
	GetCollectionStatus(ctx context.Context, planID string) (status CollectionStatus, err error)
}

// DiscoveryAPI defines service discovery operations (primarily freeze2)
type DiscoveryAPI interface {
	// Interface discovery
	DiscoverInterfaces(ctx context.Context) (interfaces []ServiceInterface, err error)
	GetServiceInfo(ctx context.Context, serviceName string) (info ServiceInterface, err error)
}

// ServiceInterface describes a discovered service
type ServiceInterface struct {
	Name      string
	Version   string
	Port      string
	Namespace string
	Methods   []MethodInfo
}

// MethodInfo describes a service method
type MethodInfo struct {
	Name       string
	InputType  string
	OutputType string
	Fault      string
}

// UnifiedServiceClient is the main entry point that provides access to all unified APIs
type UnifiedServiceClient interface {
	// Get unified APIs
	Metadata() MetadataAPI
	Security() SecurityAPI
	Session() SessionAPI
	DataCollection() DataCollectionAPI
	Discovery() DiscoveryAPI

	// Lifecycle
	Close() error
}
