package cyw4343w

func atcmRamBaseAddress(chipId uint32) uint32 {
	switch chipId {
	case 0x4373:
		return 0x160000
	case 43909, 43907, 54907:
		return 0x1B0000
	case 55560:
		return 0x370000 + 0x2b4 + 0x800
	case 55900:
		return 0x3a0000 + 0x20
	case 55500:
		return 0x3a0000 + 0x2b4 + 0x800
	default:
		return 0
	}
}

func armCoreBaseAddress(chipId uint32) uint32 {
	switch chipId {
	case 0x4373, 55560, 55500:
		return 0x18002000 + wrapperRegisterOffset
	case 43012, 43022, 43430, 43439:
		return 0x18003000 + wrapperRegisterOffset
	case 43909, 43907, 54907:
		return 0x18011000 + wrapperRegisterOffset
	default:
		panic("unknown chip id")
	}
}

func socsramBaseAddress(chipId uint32, wrapper bool) uint32 {
	offset := uint32(0)
	if wrapper {
		offset = wrapperRegisterOffset
	}
	switch chipId {
	case 43012, 43022, 43430, 43439:
		return 0x18004000 + offset
	default:
		panic("unsupported chip id")
	}
}

func sdiodCoreBaseAddress(chipId uint32) uint32 {
	switch chipId {
	case 55560:
		return 0x18004000
	case 55500:
		return 0x18003000
	case 0x4373:
		return 0x18005000
	case 43012, 43022, 43430, 43439:
		return 0x18002000
	default:
		panic("unknown chip id")
	}
}

func chipRamSize(chipId uint32) uint32 {
	switch chipId {
	case 4334, 43340, 43342, 43430, 43439:
		return 512 * 1024
	case 43362, 4390:
		return 0x3C000
	case 43909, 43907, 54907:
		return 0x90000
	case 43012, 43022:
		return 0xA0000
	case 0x4373:
		return 0xE0000
	case 55560:
		return 0x150000 - 0x800 - 0x2b4
	case 55500:
		return 0xE0000 - 0x800 - 0x2b4
	default:
		return 0x80000
	}
}
