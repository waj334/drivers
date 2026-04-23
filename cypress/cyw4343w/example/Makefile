SIGOC ?= sigoc

all:
	CGO_ENABLED=1 $(SIGOC) build --dump-ir --cpu stm32h747xi -g -O0 --float hardfp --stack-size=8192 -o./bin/firmware.elf .

debug-passes:
	CGO_ENABLED=1 $(SIGOC) build --debug-lowering --dump-ir --cpu stm32h747xi -g -O0 --float hardfp --stack-size=8192 -o./bin/firmware.elf . > out.txt 2>&1

clean:
	rm -rf ./bin