//go:build linux

package capture

import (
	"fmt"
	"log/slog"
	"syscall"
)

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func (im *InterfaceManager) checkCapturePermissionsPlatform() error {
	testSocket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		if err == syscall.EPERM {
			im.logger.Error("insufficient permissions for packet capture - CAP_NET_RAW required")
			return fmt.Errorf("%w: CAP_NET_RAW capability required for packet capture", ErrInsufficientPerms)
		}
		im.logger.Debug("failed to create test socket",
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to test capture permissions: %w", err)
	}

	syscall.Close(testSocket)
	return nil
}
