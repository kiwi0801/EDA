package interfaces

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/hooklift/gowsdl/soap"

	// Import freeze1 packages
	f1metadata "github.com/eda/generated/freeze1/e125_v0305_metadataclient"
	f1security "github.com/eda/generated/freeze1/e132_v0305_securityadmin"
	f1session "github.com/eda/generated/freeze1/e132_v0305_sessionclient"
	f1dcm "github.com/eda/generated/freeze1/e134_v1105_equipment"

	// Import freeze2 packages
	f2metadata "github.com/eda/generated/freeze2/e125_v0710_eqsdclient"
	f2discovery "github.com/eda/generated/freeze2/e132_interfacediscovery"
	f2security "github.com/eda/generated/freeze2/e132_v0310_securityadmin"
	f2session "github.com/eda/generated/freeze2/e132_v0310_sessionclient"
	f2dcm "github.com/eda/generated/freeze2/e134_v0710_dcmequipment"
)

// ClientConfig holds configuration for creating unified service clients
type ClientConfig struct {
	// SOAP service endpoint URLs
	MetadataURL       string
	SecurityURL       string
	SessionURL        string
	DataCollectionURL string
	DiscoveryURL      string // Freeze2 only

	// Optional authentication
	Username string
	Password string

	// Optional TLS configuration
	TLSConfig *ClientTLSConfig
}

// ClientTLSConfig holds TLS configuration
type ClientTLSConfig struct {
	InsecureSkipVerify bool
	CertFile           string
	KeyFile            string
}

// NewFreeze1Client creates a unified service client for freeze1
func NewFreeze1Client(config ClientConfig) (UnifiedServiceClient, error) {
	// Build SOAP client options
	var opts []soap.Option
	if config.Username != "" && config.Password != "" {
		opts = append(opts, soap.WithBasicAuth(config.Username, config.Password))
	}
	if config.TLSConfig != nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.TLSConfig.InsecureSkipVerify,
		}
		opts = append(opts, soap.WithTLS(tlsConfig))
	}

	// Create SOAP clients for each service
	metadataSoapClient := soap.NewClient(config.MetadataURL, opts...)
	metadataClient := f1metadata.NewMetadataClient(metadataSoapClient)

	securitySoapClient := soap.NewClient(config.SecurityURL, opts...)
	securityClient := f1security.NewSecurityAdmin(securitySoapClient)

	sessionSoapClient := soap.NewClient(config.SessionURL, opts...)
	sessionClient := f1session.NewSessionClient(sessionSoapClient)

	dcmSoapClient := soap.NewClient(config.DataCollectionURL, opts...)
	dcmClient := f1dcm.NewDataCollectionManagement(dcmSoapClient)

	return NewFreeze1ServiceClient(metadataClient, securityClient, sessionClient, dcmClient), nil
}

// NewFreeze2Client creates a unified service client for freeze2
func NewFreeze2Client(config ClientConfig) (UnifiedServiceClient, error) {
	// Build SOAP client options
	var opts []soap.Option
	if config.Username != "" && config.Password != "" {
		opts = append(opts, soap.WithBasicAuth(config.Username, config.Password))
	}
	if config.TLSConfig != nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.TLSConfig.InsecureSkipVerify,
		}
		opts = append(opts, soap.WithTLS(tlsConfig))
	}

	// Create SOAP clients for each service
	metadataSoapClient := soap.NewClient(config.MetadataURL, opts...)
	metadataClient := f2metadata.NewMetadataClient(metadataSoapClient)

	securitySoapClient := soap.NewClient(config.SecurityURL, opts...)
	securityClient := f2security.NewSecurityAdmin(securitySoapClient)

	sessionSoapClient := soap.NewClient(config.SessionURL, opts...)
	sessionClient := f2session.NewSessionClient(sessionSoapClient)

	dcmSoapClient := soap.NewClient(config.DataCollectionURL, opts...)
	dcmClient := f2dcm.NewDataCollectionManager(dcmSoapClient)

	discoverySoapClient := soap.NewClient(config.DiscoveryURL, opts...)
	discoveryClient := f2discovery.NewInterfaceDiscovery(discoverySoapClient)

	return NewFreeze2ServiceClient(metadataClient, securityClient, sessionClient, dcmClient, discoveryClient), nil
}

// AutoDetectAndCreateClient automatically detects the service version and creates appropriate client
func AutoDetectAndCreateClient(config ClientConfig) (UnifiedServiceClient, Version, error) {
	// Try freeze2 first (more features)
	client, err := NewFreeze2Client(config)
	if err == nil {
		return client, Freeze2, nil
	}

	// Fall back to freeze1
	client, err = NewFreeze1Client(config)
	if err == nil {
		return client, Freeze1, nil
	}

	return nil, "", fmt.Errorf("could not create client for either freeze1 or freeze2")
}

// Version represents the service version
type Version string

const (
	Freeze1 Version = "freeze1"
	Freeze2 Version = "freeze2"
)

// ValidateEndpoint validates that a SOAP endpoint is reachable
func ValidateEndpoint(endpoint string) error {
	_, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}
	return nil
}

// MigrationHelper provides utilities for migrating from freeze1 to freeze2
type MigrationHelper struct {
	freeze1Client UnifiedServiceClient
	freeze2Client UnifiedServiceClient
}

// NewMigrationHelper creates a helper for migrating between versions
func NewMigrationHelper(freeze1Client, freeze2Client UnifiedServiceClient) *MigrationHelper {
	return &MigrationHelper{
		freeze1Client: freeze1Client,
		freeze2Client: freeze2Client,
	}
}

// CanMigrate checks if both clients are available
func (m *MigrationHelper) CanMigrate() bool {
	return m.freeze1Client != nil && m.freeze2Client != nil
}

// Close closes both clients
func (m *MigrationHelper) Close() error {
	if m.freeze1Client != nil {
		if err := m.freeze1Client.Close(); err != nil {
			return err
		}
	}
	if m.freeze2Client != nil {
		if err := m.freeze2Client.Close(); err != nil {
			return err
		}
	}
	return nil
}
