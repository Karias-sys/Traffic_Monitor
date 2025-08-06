//go:build darwin

package capture

import (
	"log/slog"
)

func (im *InterfaceManager) getInterfaceStatsPlatform(name string) (InterfaceStats, error) {
	// Darwin doesn't expose network interface statistics through easily accessible files
	// Statistics would typically be retrieved through system calls or specific APIs
	// For now, return empty statistics with debug logging
	im.logger.Debug("interface statistics collection not implemented for Darwin",
		slog.String("interface", name))
	
	return InterfaceStats{}, nil
}