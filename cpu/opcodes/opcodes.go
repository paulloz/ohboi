package opcodes

const (
	NOOP = 0x00

	// Interrupt
	DI = 0xf3

	// Load
	LD_NN_A     = 0xea
	LD_FF00_n_A = 0xe0
	LD_A_H      = 0x7c
	LD_A_L      = 0x7d
	LD_A_IMM    = 0x3e
	LD_SP_NN    = 0x31
	LD_HL_NN    = 0x21

	// Call
	CALL_NN = 0xcd

	// Jump
	JP_NN = 0xc3
	JR_N  = 0x18

	// ALU
	AND_A_E = 0xa3
	INC_A   = 0x3c
)
