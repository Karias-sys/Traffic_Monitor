//go:build linux

package capture

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (im *InterfaceManager) getInterfaceStatsPlatform(name string) (InterfaceStats, error) {
	var stats InterfaceStats

	rxBytesPath := fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", name)
	rxPacketsPath := fmt.Sprintf("/sys/class/net/%s/statistics/rx_packets", name)
	rxErrorsPath := fmt.Sprintf("/sys/class/net/%s/statistics/rx_errors", name)
	txBytesPath := fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", name)
	txPacketsPath := fmt.Sprintf("/sys/class/net/%s/statistics/tx_packets", name)
	txErrorsPath := fmt.Sprintf("/sys/class/net/%s/statistics/tx_errors", name)

	if rxBytes, err := readInterfaceStat(rxBytesPath); err == nil {
		stats.RxBytes = rxBytes
	}
	if rxPackets, err := readInterfaceStat(rxPacketsPath); err == nil {
		stats.RxPackets = rxPackets
	}
	if rxErrors, err := readInterfaceStat(rxErrorsPath); err == nil {
		stats.RxErrors = rxErrors
	}
	if txBytes, err := readInterfaceStat(txBytesPath); err == nil {
		stats.TxBytes = txBytes
	}
	if txPackets, err := readInterfaceStat(txPacketsPath); err == nil {
		stats.TxPackets = txPackets
	}
	if txErrors, err := readInterfaceStat(txErrorsPath); err == nil {
		stats.TxErrors = txErrors
	}

	return stats, nil
}

func readInterfaceStat(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read interface stat from %s: %w", path, err)
	}

	valueStr := strings.TrimSpace(string(data))
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse stat value '%s' from %s: %w", valueStr, path, err)
	}

	return value, nil
}