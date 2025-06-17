package cyw4343w

import (
	"time"

	"pkg.si-go.dev/chip/core/hal"
)

type coreType uint32

const (
	wlanArmCore coreType = 0
	socramCore  coreType = 1
	sdiodCore   coreType = 2
)

type wlanCoreFlagType uint32

const (
	wlanCoreFlagNone    wlanCoreFlagType = 0
	wlanCoreFlagCpuHalt wlanCoreFlagType = 1
)

func (c *Cyw4343w[SDIO]) coreAddress(coreId coreType) (uint32, error) {
	switch coreId {
	case wlanArmCore:
		return armCoreBaseAddress(c.chipId), nil
	case socramCore:
		return socsramBaseAddress(c.chipId, true), nil
	case sdiodCore:
		return sdiodCoreBaseAddress(c.chipId), nil
	default:
		return 0, hal.ErrInvalidParameter
	}
}

func (c *Cyw4343w[SDIO]) disableDeviceCore(coreId coreType, coreFlag wlanCoreFlagType) error {
	base, err := c.coreAddress(coreId)
	if err != nil {
		return err
	}

	// Read the reset control.
	_, err = c.readBackplaneValue(base+aiResetctrlOffset, 1)
	if err != nil {
		return err
	}

	// Read the reset control and check if it is already in reset.
	regdata, err := c.readBackplaneValue(base+aiResetctrlOffset, 1)
	if err != nil {
		return err
	}

	if regdata&aircReset != 0 {
		// The device is already in reset.
		return nil
	}

	// Write 0 to the IO control and read it back.
	var ioCtrlValue uint32
	if coreFlag == wlanCoreFlagCpuHalt {
		ioCtrlValue = sicfCpuhalt
	}

	err = c.writeBackplaneValue(base+aiIoctrlOffset, 1, ioCtrlValue)
	if err != nil {
		return err
	}

	_, err = c.readBackplaneValue(base+aiIoctrlOffset, 1)
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	err = c.writeBackplaneValue(base+aiResetctrlOffset, 1, aircReset)
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond)

	return nil
}

func (c *Cyw4343w[SDIO]) resetDeviceCore(coreId coreType, coreFlag wlanCoreFlagType) error {
	var temp uint32

	base, err := c.coreAddress(coreId)
	if err != nil {
		return err
	}

	err = c.disableDeviceCore(coreId, coreFlag)
	if err != nil {
		return err
	}

	var ioCtrlValue uint32
	if coreFlag == wlanCoreFlagCpuHalt {
		ioCtrlValue = sicfCpuhalt
	}

	err = c.writeBackplaneValue(base+aiIoctrlOffset, 1, sicfFgc|sicfClockEn|ioCtrlValue)
	if err != nil {
		return err
	}

	temp, err = c.readBackplaneValue(base+aiIoctrlOffset, 1)
	if err != nil {
		return err
	}

	if temp != (sicfFgc | sicfClockEn | ioCtrlValue) {
		return errInvalidCommand
	}

	err = c.writeBackplaneValue(base+aiResetctrlOffset, 1, 0)
	if err != nil {
		return err
	}

	temp, err = c.readBackplaneValue(base+aiResetctrlOffset, 1)
	if err != nil {
		return err
	}

	if temp != 0 {
		return errInvalidCommand
	}

	time.Sleep(time.Millisecond)

	err = c.writeBackplaneValue(base+aiIoctrlOffset, 1, sicfClockEn|ioCtrlValue)
	if err != nil {
		return err
	}

	temp, err = c.readBackplaneValue(base+aiIoctrlOffset, 1)
	if err != nil {
		return err
	}

	if temp != (sicfClockEn | ioCtrlValue) {
		return errInvalidCommand
	}

	time.Sleep(time.Millisecond)

	return nil
}

func (c *Cyw4343w[SDIO]) deviceCoreIsUp(coreId coreType) (bool, error) {
	base, err := c.coreAddress(coreId)
	if err != nil {
		return false, err
	}

	regdata, err := c.readBackplaneValue(base+aiIoctrlOffset, 1)
	if err != nil {
		return false, err
	}

	// Verify that the clock is enabled and something else is not on.
	if regdata&(sicfFgc|sicfClockEn) != sicfClockEn {
		return false, nil
	}

	// Read the reset control and verify that it is not in reset.
	regdata, err = c.readBackplaneValue(base+aiResetctrlOffset, 1)
	if err != nil {
		return false, err
	}

	if regdata&aircReset != 0 {
		return false, nil
	}

	return true, nil
}
