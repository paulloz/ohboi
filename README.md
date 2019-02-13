A simple GB emulator written for fun in Go.

## Current Blargg's Gameboy hardware test ROMs status

* [x] cpu_instrs
  * [x] 01-special
  * [x] 02-interrupts
  * [x] 03-op sp,hl
  * [x] 04-op r,imm
  * [x] 05-op rp
  * [x] 06-ld r,r
  * [x] 07-jr,jp,call,ret,rst
  * [x] 08-misc instrs
  * [x] 09-op r,r
  * [x] 10-bit ops
  * [x] 11-op a,(hl)
* instr_timing
  * Failed #255
* interrupt_time
  * Failed
* mem_timing
  * 01-read_timing
    * Seems like it's getting stuck in a loop
  * 02-write_timing
    * Seems like it's getting stuck in a loop
  * 03-modify_timing
    * Seems like it's getting stuck in a loop

## Resources

* [http://gbdev.gg8.se/wiki/articles/Pan_Docs](http://gbdev.gg8.se/wiki/articles/Pan_Docs)
* [http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html](http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
* [https://rednex.github.io/rgbds/gbz80.7.html](https://rednex.github.io/rgbds/gbz80.7.html)
