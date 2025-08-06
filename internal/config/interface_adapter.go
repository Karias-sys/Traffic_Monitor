package config

// InterfaceManagerAdapter provides a bridge between the capture package's InterfaceManager
// and the config package's validation requirements. It implements the InterfaceValidator interface.
type InterfaceManagerAdapter struct {
	manager InterfaceValidator
}

// NewInterfaceManagerAdapter creates a new adapter with proper type safety
func NewInterfaceManagerAdapter(manager InterfaceValidator) *InterfaceManagerAdapter {
	return &InterfaceManagerAdapter{
		manager: manager,
	}
}

// ValidateInterface delegates to the underlying interface manager
func (ima *InterfaceManagerAdapter) ValidateInterface(nameOrIndex string) error {
	return ima.manager.ValidateInterface(nameOrIndex)
}

// GetDefaultInterfaceForConfig delegates to the underlying interface manager
func (ima *InterfaceManagerAdapter) GetDefaultInterfaceForConfig() (*InterfaceInfo, error) {
	return ima.manager.GetDefaultInterfaceForConfig()
}