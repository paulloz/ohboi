package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/test"
)

func TestRom01Special(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/01-special.gb", t)
}

func TestRom02Interrupts(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/02-interrupts.gb", t)
}

func TestRom03OPSPHL(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/03-op sp,hl.gb", t)
}

func TestRom04OPRIMM(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/04-op r,imm.gb", t)
}

func TestRom05OPRP(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/05-op rp.gb", t)
}

func TestRom06LDRR(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/06-ld r,r.gb", t)
}

func TestRom07JRJPCALLRETRST(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/07-jr,jp,call,ret,rst.gb", t)
}

func TestRom08MISC(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/08-misc instrs.gb", t)
}

func TestRom09OPRR(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/09-op r,r.gb", t)
}

func TestRom10BIT(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/10-bit ops.gb", t)
}

func TestRom11OPA(t *testing.T) {
	test.ExecuteROMTest("cpu_instrs/individual/11-op a,(hl).gb", t)
}
