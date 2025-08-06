//go:build windows

package capture

import (
	"log/slog"
)

func (im *InterfaceManager) getInterfaceStatsPlatform(name string) (InterfaceStats, error) {
	// Windows interface statistics would typically be retrieved through WMI or Windows APIs
	// For now, return empty statistics with debug logging
	im.logger.Debug("interface statistics collection not implemented for Windows",
		slog.String("interface", name))
	
	return InterfaceStats{}, nil
}