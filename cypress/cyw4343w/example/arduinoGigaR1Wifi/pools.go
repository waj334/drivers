package main

import (
	"cache"
)

const (
	// SDPCM+BCD+ethernet maxframe, rounded to 32-byte cache lines.
	txSlotSize  = 8192
	txSlotCount = 8

	// IOCTL command + iovar name + payload, conservative.
	ctrlSlotSize  = 512
	ctrlSlotCount = 4

	// RX frames are the same size class as TX.
	rxSlotSize  = 8192
	rxSlotCount = 4
)

type Cache = cache.NoCache

//sigo:section txBacking .axi_sram
//sigo:attribute retain
//sigo:align 32
var txBacking [txSlotCount * txSlotSize]byte

//sigo:section rxBacking .axi_sram
//sigo:attribute retain
//sigo:align 32
var rxBacking [rxSlotCount * rxSlotSize]byte

//sigo:section ctrlBacking .axi_sram
//sigo:attribute retain
//sigo:align 32
var ctrlBacking [ctrlSlotCount * ctrlSlotSize]byte
