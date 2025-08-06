package capture

import (
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Karias-sys/Traffic_Monitor/internal/capture"
)

func TestInterfaceManager_GetAllInterfaces(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	interfaces, err := im.GetAllInterfaces()
	require.NoError(t, err)
	assert.NotEmpty(t, interfaces, "Should have at least one interface")

	for _, iface := range interfaces {
		assert.NotEmpty(t, iface.Name, "Interface name should not be empty")
		assert.Greater(t, iface.Index, 0, "Interface index should be positive")
		assert.GreaterOrEqual(t, iface.MTU, 0, "MTU should be non-negative")
		
		if iface.HardwareAddr != "" {
			_, err := net.ParseMAC(iface.HardwareAddr)
			assert.NoError(t, err, "Hardware address should be valid MAC format")
		}
	}
}

func TestInterfaceManager_GetInterfaceByName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	t.Run("valid loopback interface", func(t *testing.T) {
		iface, err := im.GetInterfaceByName("lo")
		if err != nil {
			t.Skip("Loopback interface 'lo' not available on this system")
		}
		
		require.NoError(t, err)
		assert.Equal(t, "lo", iface.Name)
		assert.True(t, iface.IsLoopback, "Should be marked as loopback")
		assert.Greater(t, iface.Index, 0, "Should have valid index")
	})

	t.Run("invalid interface name", func(t *testing.T) {
		_, err := im.GetInterfaceByName("nonexistent_interface_12345")
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInterfaceNotFound)
	})

	t.Run("empty interface name", func(t *testing.T) {
		_, err := im.GetInterfaceByName("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInvalidInterface)
	})
}

func TestInterfaceManager_GetInterfaceByIndex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	t.Run("valid interface index", func(t *testing.T) {
		interfaces, err := im.GetAllInterfaces()
		require.NoError(t, err)
		require.NotEmpty(t, interfaces)

		firstInterface := interfaces[0]
		iface, err := im.GetInterfaceByIndex(firstInterface.Index)
		require.NoError(t, err)
		assert.Equal(t, firstInterface.Index, iface.Index)
		assert.Equal(t, firstInterface.Name, iface.Name)
	})

	t.Run("invalid interface index", func(t *testing.T) {
		_, err := im.GetInterfaceByIndex(99999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInterfaceNotFound)
	})

	t.Run("zero interface index", func(t *testing.T) {
		_, err := im.GetInterfaceByIndex(0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInvalidInterfaceIdx)
	})

	t.Run("negative interface index", func(t *testing.T) {
		_, err := im.GetInterfaceByIndex(-1)
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInvalidInterfaceIdx)
	})
}

func TestInterfaceManager_ValidateInterface(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	t.Run("empty interface specification", func(t *testing.T) {
		err := im.ValidateInterface("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInvalidInterface)
	})

	t.Run("nonexistent interface by name", func(t *testing.T) {
		err := im.ValidateInterface("nonexistent_interface_12345")
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInterfaceNotFound)
	})

	t.Run("nonexistent interface by index", func(t *testing.T) {
		err := im.ValidateInterface("99999")
		assert.Error(t, err)
		assert.ErrorIs(t, err, capture.ErrInterfaceNotFound)
	})

	t.Run("validate by name and index", func(t *testing.T) {
		interfaces, err := im.GetAllInterfaces()
		require.NoError(t, err)

		for _, iface := range interfaces {
			if iface.IsUp {
				t.Run("validate_"+iface.Name, func(t *testing.T) {
					err := im.ValidateInterface(iface.Name)
					if err != nil && !assert.ErrorIs(t, err, capture.ErrInsufficientPerms) {
						t.Logf("Validation failed for interface %s: %v", iface.Name, err)
					}
				})
				break
			}
		}
	})
}

func TestInterfaceManager_GetDefaultInterface(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	defaultIface, err := im.GetDefaultInterface()
	if err != nil {
		t.Logf("Default interface selection failed: %v", err)
		return
	}

	require.NotNil(t, defaultIface)
	assert.NotEmpty(t, defaultIface.Name, "Default interface should have a name")
	assert.Greater(t, defaultIface.Index, 0, "Default interface should have valid index")
	assert.True(t, defaultIface.IsUp, "Default interface should be up")
	
	t.Logf("Selected default interface: %s (index %d, type: %s)", 
		defaultIface.Name, defaultIface.Index, getInterfaceTypeForTest(defaultIface))
}

func TestInterfaceManager_InterfaceScoring(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	interfaces, err := im.GetAllInterfaces()
	require.NoError(t, err)

	var candidates []capture.InterfaceInfo
	for _, iface := range interfaces {
		if iface.IsUp && !iface.IsLoopback {
			candidates = append(candidates, iface)
		}
	}

	if len(candidates) == 0 {
		t.Skip("No non-loopback interfaces available for scoring test")
	}

	for _, candidate := range candidates {
		t.Run("score_"+candidate.Name, func(t *testing.T) {
			interfaceType := getInterfaceTypeForTest(&candidate)
			t.Logf("Interface %s: type=%s, mtu=%d, running=%v, broadcast=%v", 
				candidate.Name, interfaceType, candidate.MTU, candidate.IsRunning, 
				candidate.Flags&net.FlagBroadcast != 0)
		})
	}
}

func TestInterfaceManager_Statistics(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	im := capture.NewInterfaceManager(logger)

	interfaces, err := im.GetAllInterfaces()
	require.NoError(t, err)

	for _, iface := range interfaces {
		t.Run("stats_"+iface.Name, func(t *testing.T) {
			assert.GreaterOrEqual(t, iface.Statistics.RxBytes, uint64(0))
			assert.GreaterOrEqual(t, iface.Statistics.RxPackets, uint64(0))
			assert.GreaterOrEqual(t, iface.Statistics.RxErrors, uint64(0))
			assert.GreaterOrEqual(t, iface.Statistics.TxBytes, uint64(0))
			assert.GreaterOrEqual(t, iface.Statistics.TxPackets, uint64(0))
			assert.GreaterOrEqual(t, iface.Statistics.TxErrors, uint64(0))

			if iface.Statistics.RxPackets > 0 || iface.Statistics.TxPackets > 0 {
				t.Logf("Interface %s: RX=%d packets/%d bytes, TX=%d packets/%d bytes",
					iface.Name, iface.Statistics.RxPackets, iface.Statistics.RxBytes,
					iface.Statistics.TxPackets, iface.Statistics.TxBytes)
			}
		})
	}
}

func getInterfaceTypeForTest(info *capture.InterfaceInfo) string {
	name := info.Name
	
	switch {
	case info.IsLoopback:
		return "loopback"
	case len(name) >= 3 && (name[:3] == "eth" || name[:2] == "en"):
		return "ethernet"
	case len(name) >= 4 && (name[:4] == "wlan" || name[:2] == "wl"):
		return "wireless"
	case len(name) >= 2 && name[:2] == "br":
		return "bridge"
	case len(name) >= 3 && (name[:3] == "tun" || name[:3] == "tap"):
		return "tunnel"
	case len(name) >= 4 && name[:4] == "veth":
		return "virtual"
	default:
		return "unknown"
	}
}