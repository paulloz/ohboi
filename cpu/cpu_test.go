package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/test"
)

func TestGBRom01Special(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/01-special.gb", t)
}

func TestGBRom02Interrupts(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/02-interrupts.gb", t)
}

func TestGBRom03OPSPHL(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/03-op sp,hl.gb", t)
}

func TestGBRom04OPRIMM(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/04-op r,imm.gb", t)
}

func TestGBRom05OPRP(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/05-op rp.gb", t)
}

func TestGBRom06LDRR(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/06-ld r,r.gb", t)
}

func TestGBRom07JRJPCALLRETRST(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb", t)
}

func TestGBRom08MISC(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/08-misc instrs.gb", t)
}

func TestGBRom09OPRR(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/09-op r,r.gb", t)
}

func TestGBRom10BIT(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/10-bit ops.gb", t)
}

func TestGBRom11OPA(t *testing.T) {
	test.ExecuteGBROMTest("gb-test-roms/cpu_instrs/individual/11-op a,(hl).gb", t)
}

func TestGBRomBitsMemOAM(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/bits/mem_oam.gb", t)
}

func TestGBRomBitsRegF(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/bits/reg_f.gb", t)
}

func TestGBRomBitsUnusedHwioGS(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/bits/unused_hwio-GS.gb", t)
}

/*
// Only works with CGB IO registers uncommented
func TestRomBitsUnusedHwioC(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/misc/bits/unused_hwio-C.gb", t)
}
*/

func TestMooneyeRomInstrDaa(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/instr/daa.gb", t)
}

func TestMooneyeRomDivTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/div_timing.gb", t)
}

func TestMooneyeRomJpTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/jp_timing.gb", t)
}

func TestMooneyeRomJpCcTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/jp_cc_timing.gb", t)
}

func TestMooneyeRomRetTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/ret_timing.gb", t)
}

func TestMooneyeRomRetCcTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/ret_cc_timing.gb", t)
}

func TestMooneyeRomRetiTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/reti_timing.gb", t)
}

func TestMooneyeRomCallTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/call_timing.gb", t)
}

func TestMooneyeRomCallCcTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/call_cc_timing.gb", t)
}

func TestMooneyeRomCallTiming2(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/call_timing2.gb", t)
}

func TestMooneyeRomCallCcTiming2(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/call_cc_timing2.gb", t)
}

func TestMooneyeRomIntrTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/intr_timing.gb", t)
}

func TestMooneyeRomAddSpETiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/add_sp_e_timing.gb", t)
}

func TestMooneyeRomOAMDMABasic(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/oam_dma/basic.gb", t)
}

func TestMooneyeRomOAMDMATiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/oam_dma_timing.gb", t)
}

func TestMooneyeRomOAMDMARegRead(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/oam_dma/reg_read.gb", t)
}

func TestMooneyeRomOAMDMASources(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/oam_dma/sources-dmgABCmgbS.gb", t)
}

func TestMooneyeRomOAMDMARestart(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/oam_dma_restart.gb", t)
}

func TestMooneyeRomTimerDivWrite(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/timer/div_write.gb", t)
}

func TestMooneyeRomEiSequence(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/ei_sequence.gb", t)
}

func TestMooneyeRomEiTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/ei_timing.gb", t)
}

func TestMooneyeRomHalfIme0Ei(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/halt_ime0_ei.gb", t)
}

func TestMooneyeRomIfIeRegisters(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/if_ie_registers.gb", t)
}

func TestMooneyeRomPushTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/push_timing.gb", t)
}

func TestMooneyeRomPopTiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/pop_timing.gb", t)
}

func TestMooneyeRomLdHlSpETiming(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/ld_hl_sp_e_timing.gb", t)
}

func TestMooneyeRomBitsRamEn(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/bits_ram_en.gb", t)
}

func TestMooneyeRomRam256Kb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/ram_256Kb.gb", t)
}

func TestMooneyeRomRam64Kb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/ram_64Kb.gb", t)
}

/*
func TestMooneyeRomRom512Kb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/rom_512Kb.gb", t)
}

func TestMooneyeRomRom1Mb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/rom_1Mb.gb", t)
}

func TestMooneyeRomRom2Mb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/rom_2Mb.gb", t)
}
*/

func TestMooneyeRomRom4Mb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/rom_4Mb.gb", t)
}

/*
func TestMooneyeRomRom8Mb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/rom_8Mb.gb", t)
}

func TestMooneyeRomRom16Mb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/rom_16Mb.gb", t)
}

func TestMooneyeRomRMulticartRom8Mb(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/emulator-only/mbc1/multicart_rom_8Mb.gb", t)
}
*/

func TestMooneyeRomTim00(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/timer/tim00.gb", t)
}

func TestMooneyeRomTim10(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/timer/tim10.gb", t)
}

func TestMooneyeRomTim11(t *testing.T) {
	test.ExecuteMooneyeROMTest("mooneye/acceptance/timer/tim11.gb", t)
}
