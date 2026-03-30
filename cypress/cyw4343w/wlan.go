package cyw4343w

import (
	"encoding/binary"
	"time"
	"unsafe"
)

const paramsSize = unsafe.Sizeof(escanParamsType{})

type ssidType struct {
	length uint32 // wlc_ssid_t.SSID_len is uint32_t
	value  [32]byte
}

func (s ssidType) String() string {
	return unsafe.String(&s.value[0], s.length)
}

// escanResultType mirrors wl_escan_result_t (fixed header portion).
type escanResultType struct {
	buflen   uint32
	version  uint32
	syncId   uint16
	bssCount uint16
	// wl_bss_info_t entries follow immediately
}

// bssInfoType mirrors the fixed fields of wl_bss_info_t up through the SSID.
type bssInfoType struct {
	version      uint32
	length       uint32
	bssid        macType
	beaconPeriod uint16
	capability   uint16
	ssidLen      uint8
	ssid         [32]byte
	_            [1]byte // reserved1 — padding to keep rateset.count 4-byte aligned
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

func (c *Cyw4343w[SDIO]) Up() error {
	buf := make([]byte, ioctlCommandHeaderLength)
	_, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcUp,
		cmdType: cdcSet,
		data:    buf,
	})
	return err
}

// setEventMask sends "bsscfg:event_msgs" to tell the firmware which async events
// to deliver. Without this the firmware sends no events at all.
// eventNums are the wlcE* event type constants to enable.
func (c *Cyw4343w[SDIO]) setEventMask(eventNums ...uint32) error {
	const maskLen = 16 // WL_EVENTING_MASK_LEN (128 bits)
	var buf [4 + maskLen]byte
	// First 4 bytes are the bsscfgidx (0 = primary interface).
	binary.LittleEndian.PutUint32(buf[:4], 0)
	mask := buf[4:]
	for _, e := range eventNums {
		if e < uint32(maskLen)*8 {
			mask[e/8] |= 1 << (e % 8)
		}
	}
	_, err := c.SetIovar("bsscfg:"+IovarStrEventMsgs, buf[:])
	return err
}

func (c *Cyw4343w[SDIO]) ScanWifiNetworks() ([]string, error) {
	// Register for escan result events before starting the scan.
	if err := c.setEventMask(wlcEEscanResult); err != nil {
		return nil, err
	}

	// Allocate buffer for scan parameters.
	buffer, payload := iovarBuffer(IovarStrEscan, int(paramsSize))

	// Set the scan parameters.
	scanParams := (*escanParamsType)(unsafe.Pointer(&payload[0]))
	scanParams.version = 1
	scanParams.action = 1
	scanParams.params.scanType = uint8(whdScanTypeActive)
	scanParams.params.bssType = uint8(whdBssTypeAny)

	// Set BSSID to broadcast (ff:ff:ff:ff:ff:ff) — wildcard meaning "any AP".
	// Leaving it as all zeros means "match only BSSID 00:00:00:00:00:00" which matches nothing.
	scanParams.params.bssid = macType{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	// Set timing params to -1 (0xFFFFFFFF) to use firmware defaults.
	scanParams.params.nprobes = 0xFFFFFFFF
	scanParams.params.activeTime = 0xFFFFFFFF
	scanParams.params.passiveTime = 0xFFFFFFFF
	scanParams.params.homeTime = 0xFFFFFFFF

	// Send the command and wait for the response.
	_, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcSetVar,
		cmdType: cdcSet,
		data:    buffer,
	})

	if err != nil {
		return nil, err
	}

	// Wait for the async escan result events.
	// Events are processed by a background polling goroutine.
	var networks []string
	networkMap := make(map[string]bool)

	deadline := time.Now().Add(10 * time.Second)
	for {
		// Try to dequeue an escan result event.
		eventData, ok := c.eventQueue.Dequeue(wlcEEscanResult)
		if !ok {
			if time.Now().After(deadline) {
				break
			}
			// Wait for the background polling goroutine to deliver events.
			time.Sleep(time.Millisecond)
			continue
		}

		// Parse the BCD header to locate the event structure.
		header := (*bcdHeaderType)(unsafe.Pointer(&eventData[0]))
		eventOffset := int(unsafe.Sizeof(bcdHeaderType{})) + int(header.dataOffset)*4
		if eventOffset+int(unsafe.Sizeof(eventType{})) > len(eventData) {
			continue
		}
		event := (*eventType)(unsafe.Pointer(&eventData[eventOffset]))

		// Check if this is the last event (status indicates completion).
		if event.event.status != wlcEStatusPartial {
			// Scan complete or error - return results.
			break
		}

		if event.event.dataLength == 0 {
			continue
		}

		// The scan result data (wl_escan_result_t) follows the event structure.
		scanDataOffset := eventOffset + int(unsafe.Sizeof(eventType{}))
		if scanDataOffset+int(unsafe.Sizeof(escanResultType{})) > len(eventData) {
			continue
		}
		scanData := eventData[scanDataOffset:]

		escanResult := (*escanResultType)(unsafe.Pointer(&scanData[0]))
		if escanResult.bssCount == 0 {
			continue
		}

		// bss_info[0] follows immediately after the fixed escan result header.
		bssOffset := int(unsafe.Sizeof(escanResultType{}))
		if bssOffset+int(unsafe.Sizeof(bssInfoType{})) > len(scanData) {
			continue
		}
		bss := (*bssInfoType)(unsafe.Pointer(&scanData[bssOffset]))
		if bss.ssidLen > 0 && bss.ssidLen <= 32 {
			ssid := string(bss.ssid[:bss.ssidLen])
			if !networkMap[ssid] {
				networks = append(networks, ssid)
				networkMap[ssid] = true
			}
		}
	}

	return networks, nil
}
