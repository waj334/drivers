#define EXTRA_MEMORY_REGIONS \
    SDRAM_BANK1 (rwx) : ORIGIN = 0xC0000000, LENGTH = 8M

#define HEAP_REGION SDRAM_BANK1
#define STACK_REGION DTCM

#include <linker_stm32h747xi.ld>

/* Add a section for things explicitly placed in SDRAM */
SECTIONS
{
    .sdram_bank1 (NOLOAD) : ALIGN(32)
    {
        __sdram_bank1_start = .;
        KEEP(*(.sdram_bank1))
        KEEP(*(.sdram_bank1.*))
        . = ALIGN(32);
        __sdram_bank1_end = .;
    } > SDRAM_BANK1
} INSERT AFTER .sram_backup;