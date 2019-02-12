A simple GB emulator written for fun in Go.

## Current Blargg's Gameboy hardware test ROMs status

* cpu_instrs
  * [ ] 01-special
    * Needs DAA
  * [x] 02-interrupts
  * [x] 03-op sp,hl
  * [x] 04-op r,imm
  * [x] 05-op rp
  * [x] 06-ld r,r
  * [ ] 07-jr,jp,call,ret,rst
    * Seems like it gets stuck in an infinite loop
  * [x] 08-misc instrs
  * [x] 09-op r,r
  * [x] 10-bit ops
  * [ ] 11-op a,(hl)
    * Needs DAA

## Resources

* [http://gbdev.gg8.se/wiki/articles/Pan_Docs](http://gbdev.gg8.se/wiki/articles/Pan_Docs)
* [http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html](http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
* [https://rednex.github.io/rgbds/gbz80.7.html](https://rednex.github.io/rgbds/gbz80.7.html)