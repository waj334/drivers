package main

//sigo:export wake runtime.wake
func wake(t uint64)

func alarm(t uint64) {
	wake(t * timescale)
}

//sigo:export nanotime runtime.nanotime
func nanotime() uint64 {
	// The timer resolution is 1uS per tick.
	return TIM2.Tick() * timescale
}

//sigo:export addsleep runtime.addsleep
func addsleep(deadline uint64) {
	TIM2.SetAlarm(deadline/timescale, alarm)
}
