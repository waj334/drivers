package cyw4343w

import (
	"context"
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

// sendIoctlU32 sends an IOCTL command with a single uint32 payload value.
func (c *Cyw4343w[HostT, CacheT]) sendIoctlU32(cmd uint32, cmdType sdpcmCommandType, value uint32) error {
	buf := make([]byte, ioctlCommandHeaderLength+4)
	binary.LittleEndian.PutUint32(buf[ioctlCommandHeaderLength:], value)
	handle, err := c.sendIoctl(ioctlCommandType{
		cmd:     cmd,
		cmdType: cmdType,
		data:    buf,
	})

	if err != nil {
		return err
	}

	handle.Close()
	return nil
}

// MACAddress returns the 6-byte hardware MAC address of the WiFi interface.
func (c *Cyw4343w[HostT, CacheT]) MACAddress() ([6]byte, error) {
	handle, err := c.Iovar(IovarStrCurEtheraddr, 6)
	if err != nil {
		return [6]byte{}, err
	}
	var mac [6]byte
	copy(mac[:], handle.Data[:6])
	handle.Close()
	return mac, nil
}

func (c *Cyw4343w[HostT, CacheT]) Up() error {
	buf := make([]byte, ioctlCommandHeaderLength)
	handle, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcUp,
		cmdType: cdcSet,
		data:    buf,
	})

	if err != nil {
		return err
	}

	handle.Close()
	return nil
}

// JoinWPA2 connects to a WPA2-PSK protected WiFi network.
func (c *Cyw4343w[HostT, CacheT]) JoinWPA2(ctx context.Context, ssid string, passphrase string) error {
	if len(ssid) > 32 {
		return errSSIDTooLong
	}
	if len(passphrase) > 64 {
		return errPassphraseTooLong
	}

	waiter := c.events.Watch(
		wlcESetSsid,
		wlcEAuth,
		wlcEAssoc,
		wlcELink,
		wlcEPskSup,
		wlcEDisassoc,
		wlcEDeauth,
	)
	if waiter == nil {
		return errPoolExhausted
	}

	// Register for join-related events before starting association so no early
	// event can be dropped as "unsubscribed".
	if err := c.setEventMask(
		wlcESetSsid,
		wlcEAuth,
		wlcEAssoc,
		wlcELink,
		wlcEPskSup,
		wlcEDisassoc,
		wlcEDeauth,
	); err != nil {
		c.events.Unwatch(waiter)
		return err
	}

	// Set infrastructure mode.
	if err := c.sendIoctlU32(wlcSetInfra, cdcSet, 1); err != nil {
		c.events.Unwatch(waiter)
		return err
	}

	// Set auth mode to open system.
	if err := c.sendIoctlU32(wlcSetAuth, cdcSet, 0); err != nil {
		c.events.Unwatch(waiter)
		return err
	}

	// Set wireless security to AES.
	if err := c.sendIoctlU32(wlcSetWsec, cdcSet, aesEnabled); err != nil {
		c.events.Unwatch(waiter)
		return err
	}

	// Set WPA2-PSK auth.
	if err := c.sendIoctlU32(wlcSetWpaAuth, cdcSet, wpa2AuthPsk); err != nil {
		c.events.Unwatch(waiter)
		return err
	}

	// Enable the internal WPA supplicant.
	var supWpa [4]byte
	binary.LittleEndian.PutUint32(supWpa[:], 1)

	handle, err := c.SetIovar(IovarStrSupWpa, supWpa[:])
	if err != nil {
		c.events.Unwatch(waiter)
		return err
	}
	handle.Close()

	// Set the passphrase using wsec_pmk_t.
	{
		pmkBuf, pmk := ioctlPacket(wsecPmkType{})
		pmk.keyLen = uint16(len(passphrase))
		pmk.flags = wsecPassphrase
		copy(pmk.key[:], passphrase)

		handle, err := c.sendIoctl(ioctlCommandType{
			cmd:     wlcSetWsecPmk,
			cmdType: cdcSet,
			data:    pmkBuf,
		})
		if err != nil {
			c.events.Unwatch(waiter)
			return err
		}
		handle.Close()
	}

	// Set the SSID. This triggers association.
	{
		ssidBuf, s := ioctlPacket(ssidType{})
		s.length = uint32(len(ssid))
		copy(s.value[:], ssid)

		handle, err := c.sendIoctl(ioctlCommandType{
			cmd:     wlcSetSsid,
			cmdType: cdcSet,
			data:    ssidBuf,
		})
		if err != nil {
			c.events.Unwatch(waiter)
			return err
		}
		handle.Close()
	}

	deadline, _ := ctx.Deadline()
	for {
		ev, ok := waiter.Pop()
		if !ok {
			if !deadline.IsZero() && time.Now().After(deadline) {
				c.events.Unwatch(waiter)
				return errTimeout
			}

			time.Sleep(time.Millisecond)
			continue
		}

		switch ev.Type {
		case wlcESetSsid:
			status := ev.Status
			ev.Close()
			c.events.Unwatch(waiter)

			if status == wlcEStatusSuccess {
				return nil
			}

			return errJoinFailed

		case wlcEPskSup:
			status := ev.Status
			ev.Close()
			switch status {
			case wlcSupKeyed:
				// Supplicant completed. SetSsid will follow with success.
				// Do nothing — wait for SetSsid event to fire return nil.
			case wlcSupTimeout:
				c.events.Unwatch(waiter)
				return errJoinFailed
			default:
				// Intermediate states (CONNECTING, AUTHENTICATING, KEYXCHANGE, etc.)
				// and offset-shifted error codes >= LAST_BASIC_STATE + offset.
				if status >= wlcSupLastBasicState+wlcSupStatusOffset {
					// Real error from supplicant.
					c.events.Unwatch(waiter)
					return errJoinFailed
				}
				// Otherwise it's a progress event; keep waiting.
			}
		case wlcEDisassoc, wlcEDeauth:
			ev.Close()
			c.events.Unwatch(waiter)
			return errJoinFailed

		default:
			// Auth, assoc, link, and any other watched join-progress events are
			// consumed so their RX handles do not leak.
			ev.Close()
		}
	}
}

// Disconnect disassociates from the current WiFi network.
func (c *Cyw4343w[HostT, CacheT]) Disconnect() error {
	buf := make([]byte, ioctlCommandHeaderLength)
	handle, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcDisassoc,
		cmdType: cdcSet,
		data:    buf,
	})

	if err != nil {
		return err
	}

	handle.Close()
	return nil
}

// setEventMask sends "bsscfg:event_msgs" to tell the firmware which async events
// to deliver. Without this the firmware sends no events at all.
// eventNums are the wlcE* event type constants to enable.
func (c *Cyw4343w[HostT, CacheT]) setEventMask(eventNums ...uint32) error {
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
	handle, err := c.SetIovar("bsscfg:"+IovarStrEventMsgs, buf[:])
	if err != nil {
		return err
	}

	handle.Close()
	return nil
}

func (c *Cyw4343w[HostT, CacheT]) ScanWifiNetworks() ([]string, error) {
	waiter := c.events.Watch(wlcEEscanResult)
	if waiter == nil {
		return nil, errPoolExhausted
	}

	// Register for escan result events before starting the scan.
	// The waiter must exist first so early events are claimed instead of dropped.
	if err := c.setEventMask(wlcEEscanResult); err != nil {
		c.events.Unwatch(waiter)
		return nil, err
	}

	// Allocate buffer for scan parameters.
	buffer, payload, err := c.iovarBuffer(IovarStrEscan, int(paramsSize))
	if err != nil {
		c.events.Unwatch(waiter)
		return nil, err
	}

	// Set the scan parameters.
	scanParams := (*escanParamsType)(unsafe.Pointer(&payload.Data[0]))
	scanParams.version = 1
	scanParams.action = 1
	scanParams.params.scanType = uint8(whdScanTypeActive)
	scanParams.params.bssType = uint8(whdBssTypeAny)

	// Wildcard BSSID: scan for any AP.
	scanParams.params.bssid = macType{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	// Use firmware defaults.
	scanParams.params.nprobes = 0xFFFFFFFF
	scanParams.params.activeTime = 0xFFFFFFFF
	scanParams.params.passiveTime = 0xFFFFFFFF
	scanParams.params.homeTime = 0xFFFFFFFF

	// Send the escan command.
	ioctlHandle, err := c.sendIoctl(ioctlCommandType{
		cmd:     wlcSetVar,
		cmdType: cdcSet,
		data:    buffer.Data,
	})

	// Release the command buffer back to the control pool.
	buffer.Close()

	if err != nil {
		c.events.Unwatch(waiter)
		return nil, err
	}

	ioctlHandle.Close()

	var networks []string
	networkMap := make(map[string]bool)

	deadline := time.Now().Add(10 * time.Second)

	for {
		ev, ok := waiter.Pop()
		if !ok {
			if time.Now().After(deadline) {
				c.events.Unwatch(waiter)
				return networks, nil
			}

			time.Sleep(time.Millisecond)
			continue
		}

		// From here on, this function owns ev.Handle and must close ev.
		if ev.Status != wlcEStatusPartial {
			ev.Close()
			c.events.Unwatch(waiter)
			return networks, nil
		}

		if ev.PayloadOffset >= len(ev.Data) {
			ev.Close()
			continue
		}

		scanData := ev.Data[ev.PayloadOffset:]

		if len(scanData) < int(unsafe.Sizeof(escanResultType{})) {
			ev.Close()
			continue
		}

		escanResult := *(*escanResultType)(unsafe.Pointer(&scanData[0]))
		if escanResult.bssCount == 0 {
			ev.Close()
			continue
		}

		// bss_info[0] follows immediately after the fixed escan result header.
		bssOffset := int(unsafe.Sizeof(escanResultType{}))
		if bssOffset+int(unsafe.Sizeof(bssInfoType{})) > len(scanData) {
			ev.Close()
			continue
		}

		bss := *(*bssInfoType)(unsafe.Pointer(&scanData[bssOffset]))
		if bss.ssidLen > 0 && bss.ssidLen <= 32 {
			ssid := string(bss.ssid[:bss.ssidLen])
			if !networkMap[ssid] {
				networks = append(networks, ssid)
				networkMap[ssid] = true
			}
		}

		ev.Close()
	}
}
