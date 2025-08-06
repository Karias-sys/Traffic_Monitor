//go:build darwin

package capture

import (
	"fmt"
	"log/slog"
	"syscall"
)

func (im *InterfaceManager) checkCapturePermissionsPlatform() error {
	testSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		if err == syscall.EPERM {
			im.logger.Error("insufficient permissions for packet capture - root privileges required")
			return fmt.Errorf("%w: root privileges required for packet capture", ErrInsufficientPerms)
		}
		im.logger.Debug("failed to create test socket",
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to test capture permissions: %w", err)
	}

	syscall.Close(testSocket)
	return nil
}
