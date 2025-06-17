package cyw4343w

import (
	"time"

	"pkg.si-go.dev/chip/core/hal/sdio"
)

type Config[Host sdio.Host] struct {
	Host Host
}

type Cyw4343w[Host sdio.Host] struct {
	host Host
}

func (c *Cyw4343w[SDIO]) Configure(config Config[SDIO]) error {
	c.host = config.Host
	return nil
}

func (c *Cyw4343w[SDIO]) Initialize() error {
	var resp sdio.Response
	var err error
	ready := false

	// Send CMD0 to reset the card
	if _, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD0, Argument: 0}); err != nil {
		return err
	}

	for retry := 0; retry < 1000; retry++ {
		// Send CMD5 to get the ready status of the card.
		if resp, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD5, Argument: 0x00FF8000}); err != nil {
			time.Sleep(time.Millisecond * 5)
			continue
		} else {
			if resp[0]>>31 == 0 {
				// The card is not ready. Try again.
				time.Sleep(time.Millisecond * 5)
				continue
			}
			ready = true
		}
		break
	}

	if !ready {
		return sdio.ErrNotReady
	}

	// Send CMD3 to get the address of the card.
	resp, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD3, Argument: 0})
	if err != nil {
		return err
	}

	// Send CMD7 with the returned RCA to select the card.
	_, err = c.host.SendCommand(sdio.Command{Index: sdio.CMD7, Argument: resp[0] & 0xFFFF0000})
	if err != nil {
		return err
	}

	return nil
}
