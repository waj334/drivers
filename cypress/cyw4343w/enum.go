package cyw4343w

type resourceType uint8

const (
	wlanFirmware resourceType = iota
	wlanNvram
	wlanClm
)
