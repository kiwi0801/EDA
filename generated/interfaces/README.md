# Unified WSDL-to-Go API Layer

## Overview

This package provides a unified Go API layer that abstracts both **Freeze1** and **Freeze2** WSDL-based web services. The unified interface allows your code to work seamlessly with either version without modification.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Your Application Code                         │
│              (Uses Unified Interface APIs)                        │
└────────────┬────────────────────────────────────────────────┬────┘
             │                                                │
    ┌────────▼──────────────┐             ┌─────────▼────────────┐
    │  Freeze1 Adapter      │             │  Freeze2 Adapter     │
    │ (freeze1_adapter.go)  │             │ (freeze2_adapter.go) │
    └────────┬──────────────┘             └─────────┬────────────┘
             │                                       │
    ┌────────▼──────────────┐             ┌─────────▼────────────┐
    │ Freeze1 WSDL Clients  │             │ Freeze2 WSDL Clients │
    │ (Generated Go Code)   │             │ (Generated Go Code)  │
    └────────┬──────────────┘             └─────────┬────────────┘
             │                                       │
    ┌────────▼──────────────┐             ┌─────────▼────────────┐
    │  Freeze1 SOAP Web     │             │  Freeze2 SOAP Web    │
    │    Services           │             │    Services          │
    └───────────────────────┘             └──────────────────────┘
```

## Key Interfaces

### 1. **MetadataAPI** - Equipment and Type Definitions
```go
type MetadataAPI interface {
    GetEquipmentStructure(ctx context.Context) error
    GetEquipmentNodeDescriptions(ctx context.Context) error
    GetTypeDefinitions(ctx context.Context) error
    GetStateMachines(ctx context.Context) error
    GetSEMIObjTypes(ctx context.Context) error
    GetUnits(ctx context.Context) error
    GetExceptions(ctx context.Context) error
}
```

### 2. **SecurityAPI** - Authentication & Authorization
```go
type SecurityAPI interface {
    Login(ctx context.Context, username, password string) (sessionID string, err error)
    Logout(ctx context.Context, sessionID string) error
    ValidateSession(ctx context.Context, sessionID string) (valid bool, err error)
    CheckPermission(ctx context.Context, sessionID string, resource string) (allowed bool, err error)
}
```

### 3. **SessionAPI** - Session Management
```go
type SessionAPI interface {
    CreateSession(ctx context.Context) (sessionID string, err error)
    DestroySession(ctx context.Context, sessionID string) error
    GetSession(ctx context.Context, sessionID string) (sessionInfo interface{}, err error)
    GetActiveSessions(ctx context.Context) (sessions []string, err error)
    RefreshSession(ctx context.Context, sessionID string) error
}
```

### 4. **DataCollectionAPI** - Data Plans & Collection
```go
type DataCollectionAPI interface {
    DefinePlan(ctx context.Context, planDef interface{}) error
    GetDefinedPlanIds(ctx context.Context) (ids []string, err error)
    GetPlanDefinition(ctx context.Context, planID string) (planDef interface{}, err error)
    DeletePlan(ctx context.Context, planID string) error
    ActivatePlan(ctx context.Context, planID string) error
    GetActivePlanIds(ctx context.Context) (ids []string, err error)
    DeactivatePlan(ctx context.Context, planID string) error
    CollectData(ctx context.Context, planID string) (data interface{}, err error)
    GetCollectionStatus(ctx context.Context, planID string) (status interface{}, err error)
}
```

### 5. **DiscoveryAPI** - Service Discovery (Freeze2 only)
```go
type DiscoveryAPI interface {
    DiscoverInterfaces(ctx context.Context) (interfaces []ServiceInterface, err error)
    GetServiceInfo(ctx context.Context, serviceName string) (info ServiceInterface, err error)
}
```

### 6. **UnifiedServiceClient** - Main Entry Point
```go
type UnifiedServiceClient interface {
    Metadata() MetadataAPI
    Security() SecurityAPI
    Session() SessionAPI
    DataCollection() DataCollectionAPI
    Discovery() DiscoveryAPI
    Close() error
}
```

## Usage Examples

### Basic Client Initialization (Freeze1)

```go
package main

import (
    "context"
    "log"
    "interfaces"
)

func main() {
    config := interfaces.ClientConfig{
        MetadataURL:       "http://equipment-service:8080/services/metadata",
        SecurityURL:       "http://security-service:8080/services/security",
        SessionURL:        "http://session-service:8080/services/session",
        DataCollectionURL: "http://dcm-service:8080/services/dcm",
    }

    client, err := interfaces.NewFreeze1Client(config)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    // Use unified API
    ctx := context.Background()
    if err := client.Metadata().GetEquipmentStructure(ctx); err != nil {
        log.Fatalf("Failed to get equipment structure: %v", err)
    }
}
```

### Automatic Version Detection

```go
package main

import (
    "context"
    "log"
    "interfaces"
)

func main() {
    config := interfaces.ClientConfig{
        MetadataURL:       "http://service:8080/metadata",
        SecurityURL:       "http://service:8080/security",
        SessionURL:        "http://service:8080/session",
        DataCollectionURL: "http://service:8080/dcm",
        DiscoveryURL:      "http://service:8080/discovery", // Used if Freeze2
    }

    client, version, err := interfaces.AutoDetectAndCreateClient(config)
    if err != nil {
        log.Fatalf("Failed to detect version: %v", err)
    }
    defer client.Close()

    log.Printf("Connected to %s", version)

    // Same code works for both versions!
    ctx := context.Background()
    client.Metadata().GetTypeDefinitions(ctx)
}
```

### Working with Different Versions

```go
package main

import (
    "context"
    "log"
    "interfaces"
)

func main() {
    freeze1Config := interfaces.ClientConfig{
        MetadataURL:       "http://freeze1-service:8080/metadata",
        SecurityURL:       "http://freeze1-service:8080/security",
        SessionURL:        "http://freeze1-service:8080/session",
        DataCollectionURL: "http://freeze1-service:8080/dcm",
    }

    freeze2Config := interfaces.ClientConfig{
        MetadataURL:       "http://freeze2-service:8080/metadata",
        SecurityURL:       "http://freeze2-service:8080/security",
        SessionURL:        "http://freeze2-service:8080/session",
        DataCollectionURL: "http://freeze2-service:8080/dcm",
        DiscoveryURL:      "http://freeze2-service:8080/discovery",
    }

    f1Client, _ := interfaces.NewFreeze1Client(freeze1Config)
    f2Client, _ := interfaces.NewFreeze2Client(freeze2Config)
    defer f1Client.Close()
    defer f2Client.Close()

    ctx := context.Background()

    // Same unified API works seamlessly
    if err := f1Client.Metadata().GetEquipmentStructure(ctx); err != nil {
        log.Printf("Freeze1 error: %v", err)
    }

    if err := f2Client.Metadata().GetEquipmentStructure(ctx); err != nil {
        log.Printf("Freeze2 error: %v", err)
    }
}
```

### Using Authentication

```go
config := interfaces.ClientConfig{
    MetadataURL:       "http://service:8080/metadata",
    SecurityURL:       "http://service:8080/security",
    SessionURL:        "http://service:8080/session",
    DataCollectionURL: "http://service:8080/dcm",
    Auth: &soap.BasicAuth{
        Login:    "username",
        Password: "password",
    },
}

client, _ := interfaces.NewFreeze2Client(config)
```

## Best Practices

### 1. Always Use Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
client.Metadata().GetTypeDefinitions(ctx)
```

### 2. Handle Errors Appropriately
```go
if err := client.DataCollection().ActivatePlan(ctx, planID); err != nil {
    log.Printf("Failed to activate plan: %v", err)
    // Handle version-specific errors
}
```

### 3. Version-Agnostic Code
```go
// This code works for both freeze1 and freeze2
func getEquipmentInfo(client interfaces.UnifiedServiceClient, ctx context.Context) error {
    return client.Metadata().GetEquipmentStructure(ctx)
}
```

### 4. Graceful Fallback
```go
client, version, _ := interfaces.AutoDetectAndCreateClient(config)

if version == interfaces.Freeze2 {
    // Use Freeze2-specific features
    services, _ := client.Discovery().DiscoverInterfaces(ctx)
}
```

## Files in This Package

- **apis.go** - Core interface definitions (6 interfaces)
- **freeze1_adapter.go** - Adapts Freeze1 WSDL clients to unified API 
- **freeze2_adapter.go** - Adapts Freeze2 WSDL clients to unified API
- **factory.go** - Client factory, version detection, utilities

## Generated Packages Used

### Freeze1 (7 packages)
- `f1_e125_1_v0305_eqpmetadatamgr` 
- `f1_e125_1_v0305_metadataclient`
- `f1_e132_1_v0305_securityadmin`
- `f1_e132_1_v0305_sessionclient`
- `f1_e132_1_v0305_sessionmanager`
- `f1_e134_1_v1105_client`
- `f1_e134_1_v1105_equipment`

### Freeze2 (8 packages)
- `f2_e125_1_v0710_eqsdclient`
- `f2_e125_1_v0710_eqsdequipment`
- `f2_e132_1_v0310_securityadmin`
- `f2_e132_1_v0310_sessionclient`
- `f2_e132_1_v0310_sessionmanager`
- `f2_e132_interfacediscovery`
- `f2_e134_1_v0710_dcmclient`
- `f2_e134_1_v0710_dcmequipment`

## Migration Path

To migrate from Freeze1 to Freeze2:

1. Replace Freeze1 config with Freeze2 config
2. Update service endpoints (URLs)
3. Code using the unified API should work unchanged!
4. Optionally use `Discovery().DiscoverInterfaces()` for new capabilities

```go
// Old Freeze1 code
f1Client, _ := interfaces.NewFreeze1Client(freeze1Config)

// New Freeze2 code - same interface!
f2Client, _ := interfaces.NewFreeze2Client(freeze2Config)

// Your code doesn't need to change
processEquipment(f1Client)  // Works
processEquipment(f2Client)  // Also works!
```

## Implementation Notes

- Adapters implement basic method signatures matching unified interfaces
- Require actual SOAP client instances for full functionality
- May need method parameter mapping for complex types
- Use placeholder implementations where WSDL specifics differ

## Dependencies

- `github.com/hooklift/gowsdl/soap` - SOAP client implementation
- Standard Go libraries (`context`, `encoding/xml`, `time`)

---

**Version**: 1.0
**Status**: Ready for integration
