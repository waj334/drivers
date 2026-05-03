SIGOC ?= sigoc
stack_size ?= 8192
all:
	$(SIGOC) build --dump-ir --cpu stm32h747xi -g -O0 --float hardfp --stack-size=$(stack_size) -o./bin/firmware.elf .

debug-passes:
	$(SIGOC) build --debug-lowering --dump-ir --cpu stm32h747xi -g -O0 --float hardfp --stack-size=$(stack_size) -o./bin/firmware.elf . > out.txt 2>&1

clean:
	rm -rf ./bin