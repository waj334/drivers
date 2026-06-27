package display

import "pkg.si-go.dev/drivers/core/display"

func fillRGB888(surface display.Surface, color uint32) {
	fb := surface.Bytes()

	r := byte(color >> 16)
	g := byte(color >> 8)
	b := byte(color)

	for y := 0; y < surface.Height; y++ {
		row := fb[y*surface.Stride:]
		for x := 0; x < surface.Width; x++ {
			i := x * 3

			// RGB888 in little-endian 0x00RRGGBB memory order:
			// byte0 = B, byte1 = G, byte2 = R.
			row[i+0] = b
			row[i+1] = g
			row[i+2] = r
		}
	}
}

func fillRGB565(surface display.Surface, color uint16) {
	fb := surface.Bytes()

	lo := byte(color)
	hi := byte(color >> 8)

	for y := 0; y < surface.Height; y++ {
		row := fb[y*surface.Stride:]
		for x := 0; x < surface.Width; x++ {
			i := x * 2
			row[i+0] = lo
			row[i+1] = hi
		}
	}
}
