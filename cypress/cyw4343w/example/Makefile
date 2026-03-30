SIGOC ?= sigoc

all:
	$(SIGOC) build \
		--dump-ir \
		--cpu stm32h747xi \
		-g \
		-O0 \
		--float hardfp \
		--stack-size=8192 \
		-o./bin/firmware.elf \
		.

clean:
	rm -rf ./bin