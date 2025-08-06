package capture

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrInterfaceNotFound   = errors.New("network interface not found")
	ErrInvalidInterface    = errors.New("invalid interface specification")
	ErrInsufficientPerms   = errors.New("insufficient permissions for packet capture")
	ErrInterfaceDown       = errors.New("network interface is down")
	ErrNoCaptureSupport    = errors.New("interface does not support packet capture")
	ErrNoInterfacesFound   = errors.New("no network interfaces found")
	ErrInvalidInterfaceIdx = errors.New("invalid interface index")
)

type InterfaceInfo struct {
	Name         string
	Index        int
	MTU          int
	Flags        net.Flags
	HardwareAddr string
	IsUp         bool
	IsRunning    bool
	IsLoopback   bool
	Statistics   InterfaceStats
}

type InterfaceStats struct {
	RxBytes   uint64
	RxPackets uint64
	RxErrors  uint64
	TxBytes   uint64
	TxPackets uint64
	TxErrors  uint64
}

type InterfaceManager struct {
	logger *slog.Logger
}

func NewInterfaceManager(logger *slog.Logger) *InterfaceManager {
	return &InterfaceManager{
		logger: logger,
	}
}

func (im *InterfaceManager) GetAllInterfaces() ([]InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		im.logger.Error("failed to enumerate network interfaces", 
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to enumerate network interfaces: %w", err)
	}

	if len(interfaces) == 0 {
		im.logger.Warn("no network interfaces found")
		return nil, ErrNoInterfacesFound
	}

	var interfaceList []InterfaceInfo
	for _, iface := range interfaces {
		info := InterfaceInfo{
			Name:       iface.Name,
			Index:      iface.Index,
			MTU:        iface.MTU,
			Flags:      iface.Flags,
			IsUp:       iface.Flags&net.FlagUp != 0,
			IsRunning:  iface.Flags&net.FlagRunning != 0,
			IsLoopback: iface.Flags&net.FlagLoopback != 0,
		}

		if iface.HardwareAddr != nil {
			info.HardwareAddr = iface.HardwareAddr.String()
		}

		stats, err := im.getInterfaceStats(iface.Name)
		if err != nil {
			im.logger.Debug("failed to get interface statistics",
				slog.String("interface", iface.Name),
				slog.String("error", err.Error()))
		} else {
			info.Statistics = stats
		}

		interfaceList = append(interfaceList, info)
	}

	sort.Slice(interfaceList, func(i, j int) bool {
		return interfaceList[i].Index < interfaceList[j].Index
	})

	im.logger.Info("enumerated network interfaces", 
		slog.Int("count", len(interfaceList)))

	return interfaceList, nil
}

func (im *InterfaceManager) GetInterfaceByName(name string) (*InterfaceInfo, error) {
	if name == "" {
		return nil, ErrInvalidInterface
	}

	iface, err := net.InterfaceByName(name)
	if err != nil {
		im.logger.Debug("interface not found by name",
			slog.String("interface", name),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, name)
	}

	return im.buildInterfaceInfo(iface)
}

func (im *InterfaceManager) GetInterfaceByIndex(index int) (*InterfaceInfo, error) {
	if index <= 0 {
		return nil, ErrInvalidInterfaceIdx
	}

	iface, err := net.InterfaceByIndex(index)
	if err != nil {
		im.logger.Debug("interface not found by index",
			slog.Int("index", index),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: index %d", ErrInterfaceNotFound, index)
	}

	return im.buildInterfaceInfo(iface)
}

func (im *InterfaceManager) ValidateInterface(nameOrIndex string) error {
	if nameOrIndex == "" {
		return ErrInvalidInterface
	}

	var info *InterfaceInfo
	var err error

	if index, parseErr := strconv.Atoi(nameOrIndex); parseErr == nil {
		info, err = im.GetInterfaceByIndex(index)
	} else {
		info, err = im.GetInterfaceByName(nameOrIndex)
	}

	if err != nil {
		return err
	}

	if !info.IsUp {
		im.logger.Debug("interface validation failed: interface is down",
			slog.String("interface", nameOrIndex))
		return fmt.Errorf("%w: interface '%s' is not up", ErrInterfaceDown, nameOrIndex)
	}

	if err := im.checkCapturePermissions(); err != nil {
		return err
	}

	if !im.supportsCaptureMode(info) {
		im.logger.Debug("interface validation failed: no capture support",
			slog.String("interface", nameOrIndex),
			slog.Bool("loopback", info.IsLoopback),
			slog.Bool("broadcast", info.Flags&net.FlagBroadcast != 0),
			slog.Bool("point_to_point", info.Flags&net.FlagPointToPoint != 0))
		return fmt.Errorf("%w: interface '%s' does not support packet capture (flags: %v)", ErrNoCaptureSupport, nameOrIndex, info.Flags)
	}

	return nil
}

func (im *InterfaceManager) checkCapturePermissions() error {
	return im.checkCapturePermissionsPlatform()
}

func (im *InterfaceManager) supportsCaptureMode(info *InterfaceInfo) bool {
	if info.IsLoopback {
		return true
	}

	return info.Flags&net.FlagBroadcast != 0 || info.Flags&net.FlagPointToPoint != 0
}

func (im *InterfaceManager) buildInterfaceInfo(iface *net.Interface) (*InterfaceInfo, error) {
	info := &InterfaceInfo{
		Name:       iface.Name,
		Index:      iface.Index,
		MTU:        iface.MTU,
		Flags:      iface.Flags,
		IsUp:       iface.Flags&net.FlagUp != 0,
		IsRunning:  iface.Flags&net.FlagRunning != 0,
		IsLoopback: iface.Flags&net.FlagLoopback != 0,
	}

	if iface.HardwareAddr != nil {
		info.HardwareAddr = iface.HardwareAddr.String()
	}

	stats, err := im.getInterfaceStats(iface.Name)
	if err != nil {
		im.logger.Debug("failed to get interface statistics",
			slog.String("interface", iface.Name),
			slog.String("error", err.Error()))
	} else {
		info.Statistics = stats
	}

	return info, nil
}

func (im *InterfaceManager) getInterfaceStats(name string) (InterfaceStats, error) {
	return im.getInterfaceStatsPlatform(name)
}

func (im *InterfaceManager) GetDefaultInterface() (*InterfaceInfo, error) {
	interfaces, err := im.GetAllInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces for default selection: %w", err)
	}

	if len(interfaces) == 0 {
		return nil, ErrNoInterfacesFound
	}

	var candidates []InterfaceInfo
	for _, iface := range interfaces {
		if !iface.IsUp || !iface.IsRunning || iface.IsLoopback {
			continue
		}

		if im.supportsCaptureMode(&iface) {
			candidates = append(candidates, iface)
		}
	}

	if len(candidates) == 0 {
		im.logger.Warn("no suitable interfaces found for packet capture, falling back to any up interface")
		for _, iface := range interfaces {
			if iface.IsUp && !iface.IsLoopback {
				im.logger.Info("selected fallback interface", 
					slog.String("interface", iface.Name),
					slog.Int("index", iface.Index))
				return &iface, nil
			}
		}
		return nil, fmt.Errorf("no suitable interfaces available for packet capture")
	}

	selected := im.selectBestInterface(candidates)
	im.logger.Info("selected default interface", 
		slog.String("interface", selected.Name),
		slog.Int("index", selected.Index),
		slog.String("type", im.getInterfaceType(selected)))

	return selected, nil
}

func (im *InterfaceManager) selectBestInterface(candidates []InterfaceInfo) *InterfaceInfo {
	sort.Slice(candidates, func(i, j int) bool {
		a, b := &candidates[i], &candidates[j]
		
		scoreA := im.calculateInterfaceScore(a)
		scoreB := im.calculateInterfaceScore(b)
		
		if scoreA != scoreB {
			return scoreA > scoreB
		}

		if a.Statistics.RxPackets != b.Statistics.RxPackets {
			return a.Statistics.RxPackets > b.Statistics.RxPackets
		}

		return a.Index < b.Index
	})

	return &candidates[0]
}

func (im *InterfaceManager) calculateInterfaceScore(info *InterfaceInfo) int {
	score := 0

	interfaceType := im.getInterfaceType(info)
	switch interfaceType {
	case "ethernet":
		score += 100
	case "wireless":
		score += 80
	case "bridge":
		score += 60
	case "tunnel":
		score += 40
	case "virtual":
		score += 20
	default:
		score += 10
	}

	if info.IsRunning {
		score += 50
	}

	if info.Flags&net.FlagBroadcast != 0 {
		score += 30
	}

	if info.MTU >= 1500 {
		score += 20
	}

	if info.Statistics.RxPackets > 0 {
		score += 10
	}

	return score
}

func (im *InterfaceManager) getInterfaceType(info *InterfaceInfo) string {
	name := strings.ToLower(info.Name)
	
	switch {
	case strings.HasPrefix(name, "eth") || strings.HasPrefix(name, "en"):
		return "ethernet"
	case strings.HasPrefix(name, "wlan") || strings.HasPrefix(name, "wl") || 
		 strings.HasPrefix(name, "wifi") || strings.HasPrefix(name, "ath"):
		return "wireless"
	case strings.HasPrefix(name, "br") || strings.HasPrefix(name, "bridge"):
		return "bridge"
	case strings.HasPrefix(name, "tun") || strings.HasPrefix(name, "tap") ||
		 strings.HasPrefix(name, "vpn"):
		return "tunnel"
	case strings.HasPrefix(name, "veth") || strings.HasPrefix(name, "docker") ||
		 strings.HasPrefix(name, "lxc") || strings.HasPrefix(name, "vir"):
		return "virtual"
	case info.IsLoopback:
		return "loopback"
	default:
		return "unknown"
	}
}

func (im *InterfaceManager) GetDefaultInterfaceForConfig() (*InterfaceInfo, error) {
	info, err := im.GetDefaultInterface()
	if err != nil {
		return nil, fmt.Errorf("failed to get default interface for config: %w", err)
	}
	
	return &InterfaceInfo{
		Name:  info.Name,
		Index: info.Index,
		MTU:   info.MTU,
		Flags: info.Flags,
	}, nil
}

