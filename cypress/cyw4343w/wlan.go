package cyw4343w

import "unsafe"

const paramsSize = unsafe.Sizeof(escanParamsType{})

type ssidType struct {
	length uint8
	value  [32]byte
}

func (s ssidType) String() string {
	return unsafe.String(&s.value[0], s.length)
}

type scanResultType struct {
	ssid           ssidType
	bssid          macType
	signalStrength uint16
	maxDataRate    uint32
	bssType        bssType
	securityType   securityType
	channel        uint8
	band           band80211Type
	ccode          [2]uint8
	flags          uint8
	next           *scanResultType
	iePtr          *uint8
	ieLen          uint8
}

type scanParamsType struct {
	ssid        ssidType
	bssid       macType
	bssType     uint8
	scanType    uint8
	nprobes     uint32
	activeTime  uint32
	passiveTime uint32
	homeTime    uint32
	channelNum  uint32
	channelList [1]uint16
}

type escanParamsType struct {
	version uint32
	action  uint16
	syncId  uint16
	params  scanParamsType
}

func (c *Cyw4343w[SDIO]) ScanWifiNetworks() ([]string, error) {
	// Allocate buffer for scan parameters.
	buffer, payload := iovarBuffer(IovarStrEscan, int(paramsSize))

	// Set the scan parameters.
	scanParams := (*escanParamsType)(unsafe.Pointer(&payload[0]))
	scanParams.version = 1
	scanParams.action = 1
	scanParams.params.scanType = uint8(whdScanTypeActive)
	scanParams.params.bssType = uint8(whdBssTypeAny)

	// Send the command and wait for the response.
	_, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcSetVar,
		cmdType: cdcSet,
		data:    buffer,
	})

	if err != nil {
		return nil, err
	}

	// TODO: Wait for the async event that should contain the struct holding the list of scanned wireless networks.

	return nil, nil
}
