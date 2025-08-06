//go:build windows

package capture

import (
	"fmt"
)

func (im *InterfaceManager) checkCapturePermissionsPlatform() error {
	im.logger.Warn("permission checking not implemented for Windows - assuming sufficient permissions")
	return nil
}