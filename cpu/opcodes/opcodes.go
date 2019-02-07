package opcodes

const (
	NOOP = 0x00

	// Interrupt
	DI = 0xf3

	// LD n, nn
	LD_B_N = 0x06
	LD_C_N = 0x0e
	LD_D_N = 0x16
	LD_E_N = 0x1e
	LD_H_N = 0x26
	LD_L_N = 0x2e

	// LD r1, r2
	LD_A_A  = 0x7f
	LD_A_B  = 0x78
	LD_A_C  = 0x79
	LD_A_D  = 0x7a
	LD_A_E  = 0x7b
	LD_A_H  = 0x7c
	LD_A_L  = 0x7d
	LD_A_BC = 0x0a
	LD_A_DE = 0x1a
	LD_A_HL = 0x7e
	LD_A_NN = 0xfa
	LD_A_N  = 0x3e

	LD_B_B  = 0x40
	LD_B_C  = 0x41
	LD_B_D  = 0x42
	LD_B_E  = 0x43
	LD_B_H  = 0x44
	LD_B_L  = 0x45
	LD_B_HL = 0x46

	LD_C_B  = 0x48
	LD_C_C  = 0x49
	LD_C_D  = 0x4a
	LD_C_E  = 0x4b
	LD_C_H  = 0x4c
	LD_C_L  = 0x4d
	LD_C_HL = 0x4e

	LD_D_B  = 0x50
	LD_D_C  = 0x51
	LD_D_D  = 0x52
	LD_D_E  = 0x53
	LD_D_H  = 0x54
	LD_D_L  = 0x55
	LD_D_HL = 0x56

	LD_E_B  = 0x58
	LD_E_C  = 0x59
	LD_E_D  = 0x5a
	LD_E_E  = 0x5b
	LD_E_H  = 0x5c
	LD_E_L  = 0x5d
	LD_E_HL = 0x5e

	LD_H_B  = 0x60
	LD_H_C  = 0x61
	LD_H_D  = 0x62
	LD_H_E  = 0x63
	LD_H_H  = 0x64
	LD_H_L  = 0x65
	LD_H_HL = 0x66

	LD_L_B  = 0x68
	LD_L_C  = 0x69
	LD_L_D  = 0x6a
	LD_L_E  = 0x6b
	LD_L_H  = 0x6c
	LD_L_L  = 0x6d
	LD_L_HL = 0x6e

	LD_HL_B = 0x70
	LD_HL_C = 0x71
	LD_HL_D = 0x72
	LD_HL_E = 0x73
	LD_HL_H = 0x74
	LD_HL_L = 0x75
	LD_HL_N = 0x36

	LD_B_A  = 0x47
	LD_C_A  = 0x4f
	LD_D_A  = 0x57
	LD_E_A  = 0x5f
	LD_H_A  = 0x67
	LD_L_A  = 0x6f
	LD_BC_A = 0x02
	LD_DE_A = 0x12
	LD_HL_A = 0x77
	LD_NN_A = 0xea

	LD_A_CADDR = 0xf2
	LD_CADDR_A = 0xe2

	LD_A_HLD = 0x3a
	LD_HLD_A = 0x32

	LD_A_HLI = 0x2a
	LD_HLI_A = 0x22

	LDH_FF00N_A = 0xe0
	LDH_A_FF00N = 0xf0

	LD_BC_NN = 0x01
	LD_DE_NN = 0x11
	LD_HL_NN = 0x21
	LD_SP_NN = 0x31

	LD_SP_HL   = 0xf9
	LD_HL_SP_N = 0xf8
	LD_NN_SP   = 0x08

	// PUSH/POP
	PUSH_AF = 0xf5
	PUSH_BC = 0xc5
	PUSH_DE = 0xd5
	PUSH_HL = 0xe5
	POP_AF  = 0xf1
	POP_BC  = 0xc1
	POP_DE  = 0xd1
	POP_HL  = 0xe1

	// ALU
	ADD_A_A  = 0x87
	ADD_A_B  = 0x80
	ADD_A_C  = 0x81
	ADD_A_D  = 0x82
	ADD_A_E  = 0x83
	ADD_A_H  = 0x84
	ADD_A_L  = 0x85
	ADD_A_HL = 0x86
	ADD_A_N  = 0xc6

	ADC_A_A  = 0x8f
	ADC_A_B  = 0x88
	ADC_A_C  = 0x89
	ADC_A_D  = 0x8a
	ADC_A_E  = 0x8b
	ADC_A_H  = 0x8c
	ADC_A_L  = 0x8d
	ADC_A_HL = 0x8e
	ADC_A_N  = 0xc8

	AND_A_E = 0xa3

	CP_A  = 0xbf
	CP_B  = 0xb8
	CP_C  = 0xb9
	CP_D  = 0xba
	CP_E  = 0xbb
	CP_H  = 0xbc
	CP_L  = 0xbd
	CP_HL = 0xbe
	CP_N  = 0xfe

	DEC_A  = 0x3d
	DEC_B  = 0x05
	DEC_C  = 0x0d
	DEC_D  = 0x15
	DEC_E  = 0x1d
	DEC_H  = 0x25
	DEC_L  = 0x2d
	DEC_HL = 0x35

	INC_A  = 0x3c
	INC_B  = 0x04
	INC_C  = 0x0c
	INC_D  = 0x14
	INC_E  = 0x1c
	INC_H  = 0x24
	INC_L  = 0x2c
	INC_HL = 0x34

	SUB_A  = 0x97
	SUB_B  = 0x90
	SUB_C  = 0x91
	SUB_D  = 0x92
	SUB_E  = 0x93
	SUB_H  = 0x94
	SUB_L  = 0x95
	SUB_HL = 0x96
	SUB_N  = 0xd6

	XOR_A  = 0xaf
	XOR_B  = 0xa8
	XOR_C  = 0xa9
	XOR_D  = 0xaa
	XOR_E  = 0xab
	XOR_H  = 0xac
	XOR_L  = 0xad
	XOR_HL = 0xae
	XOR_N  = 0xee

	// 16-Bit ALU
	INC16_BC = 0x03
	INC16_DE = 0x13
	INC16_HL = 0x23
	INC16_SP = 0x33

	// Jumps
	JP_NN    = 0xc3
	JP_C_NN  = 0xda
	JP_HL    = 0xe9
	JP_NC_NN = 0xd2
	JP_NZ_NN = 0xc2
	JP_Z_NN  = 0xca
	JR_N     = 0x18 //!\ n is a signed byte value
	JR_C_N   = 0x38 //!\ n is a signed byte value
	JR_NC_N  = 0x30 //!\ n is a signed byte value
	JR_NZ_N  = 0x20 //!\ n is a signed byte value
	JR_Z_N   = 0x28 //!\ n is a signed byte value

	// Calls
	CALL_NN    = 0xcd
	CALL_C_NN  = 0xdc
	CALL_NC_NN = 0xd4
	CALL_NZ_NN = 0xc4
	CALL_Z_NN  = 0xcc

	// Restarts
	RST_00H = 0xC7
	RST_08H = 0xCF
	RST_10H = 0xD7
	RST_18H = 0xDF
	RST_20H = 0xE7
	RST_28H = 0xEF
	RST_30H = 0xF7
	RST_38H = 0xFF

	// Returns
	RET    = 0xc9
	RETI   = 0xd9
	RET_C  = 0xd8
	RET_NC = 0xd0
	RET_NZ = 0xc0
	RET_Z  = 0xc8

	// Rotates & Shifts
	RLCA = 0x07
	RLA  = 0x17

	RL_A  = 0x17
	RL_B  = 0x10
	RL_C  = 0x11
	RL_D  = 0x12
	RL_E  = 0x13
	RL_H  = 0x14
	RL_L  = 0x15
	RL_HL = 0x16

	// 16-bit OpCodes
	CB = 0xcb

	// Bit
	BIT_0_B  = 0x40
	BIT_0_C  = 0x41
	BIT_0_D  = 0x42
	BIT_0_E  = 0x43
	BIT_0_H  = 0x44
	BIT_0_L  = 0x45
	BIT_0_HL = 0x46
	BIT_0_A  = 0x47
	BIT_1_B  = 0x48
	BIT_1_C  = 0x49
	BIT_1_D  = 0x4a
	BIT_1_E  = 0x4b
	BIT_1_H  = 0x4c
	BIT_1_L  = 0x4d
	BIT_1_HL = 0x4e
	BIT_1_A  = 0x4f
	BIT_2_B  = 0x50
	BIT_2_C  = 0x51
	BIT_2_D  = 0x52
	BIT_2_E  = 0x53
	BIT_2_H  = 0x54
	BIT_2_L  = 0x55
	BIT_2_HL = 0x56
	BIT_2_A  = 0x57
	BIT_3_B  = 0x58
	BIT_3_C  = 0x59
	BIT_3_D  = 0x5a
	BIT_3_E  = 0x5b
	BIT_3_H  = 0x5c
	BIT_3_L  = 0x5d
	BIT_3_HL = 0x5e
	BIT_3_A  = 0x5f
	BIT_4_B  = 0x60
	BIT_4_C  = 0x61
	BIT_4_D  = 0x62
	BIT_4_E  = 0x63
	BIT_4_H  = 0x64
	BIT_4_L  = 0x65
	BIT_4_HL = 0x66
	BIT_4_A  = 0x67
	BIT_5_B  = 0x68
	BIT_5_C  = 0x69
	BIT_5_D  = 0x6a
	BIT_5_E  = 0x6b
	BIT_5_H  = 0x6c
	BIT_5_L  = 0x6d
	BIT_5_HL = 0x6e
	BIT_5_A  = 0x6f
	BIT_6_B  = 0x70
	BIT_6_C  = 0x71
	BIT_6_D  = 0x72
	BIT_6_E  = 0x73
	BIT_6_H  = 0x74
	BIT_6_L  = 0x75
	BIT_6_HL = 0x76
	BIT_6_A  = 0x77
	BIT_7_B  = 0x78
	BIT_7_C  = 0x79
	BIT_7_D  = 0x7a
	BIT_7_E  = 0x7b
	BIT_7_H  = 0x7c
	BIT_7_L  = 0x7d
	BIT_7_HL = 0x7e
	BIT_7_A  = 0x7f

	RES_0_B  = 0x80
	RES_0_C  = 0x81
	RES_0_D  = 0x82
	RES_0_E  = 0x83
	RES_0_H  = 0x84
	RES_0_L  = 0x85
	RES_0_HL = 0x86
	RES_0_A  = 0x87
	RES_1_B  = 0x88
	RES_1_C  = 0x89
	RES_1_D  = 0x8a
	RES_1_E  = 0x8b
	RES_1_H  = 0x8c
	RES_1_L  = 0x8d
	RES_1_HL = 0x8e
	RES_1_A  = 0x8f
	RES_2_B  = 0x90
	RES_2_C  = 0x91
	RES_2_D  = 0x92
	RES_2_E  = 0x93
	RES_2_H  = 0x94
	RES_2_L  = 0x95
	RES_2_HL = 0x96
	RES_2_A  = 0x97
	RES_3_B  = 0x98
	RES_3_C  = 0x99
	RES_3_D  = 0x9a
	RES_3_E  = 0x9b
	RES_3_H  = 0x9c
	RES_3_L  = 0x9d
	RES_3_HL = 0x9e
	RES_3_A  = 0x9f
	RES_4_B  = 0xa0
	RES_4_C  = 0xa1
	RES_4_D  = 0xa2
	RES_4_E  = 0xa3
	RES_4_H  = 0xa4
	RES_4_L  = 0xa5
	RES_4_HL = 0xa6
	RES_4_A  = 0xa7
	RES_5_B  = 0xa8
	RES_5_C  = 0xa9
	RES_5_D  = 0xaa
	RES_5_E  = 0xab
	RES_5_H  = 0xac
	RES_5_L  = 0xad
	RES_5_HL = 0xae
	RES_5_A  = 0xaf
	RES_6_B  = 0xb0
	RES_6_C  = 0xb1
	RES_6_D  = 0xb2
	RES_6_E  = 0xb3
	RES_6_H  = 0xb4
	RES_6_L  = 0xb5
	RES_6_HL = 0xb6
	RES_6_A  = 0xb7
	RES_7_B  = 0xb8
	RES_7_C  = 0xb9
	RES_7_D  = 0xba
	RES_7_E  = 0xbb
	RES_7_H  = 0xbc
	RES_7_L  = 0xbd
	RES_7_HL = 0xbe
	RES_7_A  = 0xbf

	SET_0_B  = 0xc0
	SET_0_C  = 0xc1
	SET_0_D  = 0xc2
	SET_0_E  = 0xc3
	SET_0_H  = 0xc4
	SET_0_L  = 0xc5
	SET_0_HL = 0xc6
	SET_0_A  = 0xc7
	SET_1_B  = 0xc8
	SET_1_C  = 0xc9
	SET_1_D  = 0xca
	SET_1_E  = 0xcb
	SET_1_H  = 0xcc
	SET_1_L  = 0xcd
	SET_1_HL = 0xce
	SET_1_A  = 0xcf
	SET_2_B  = 0xd0
	SET_2_C  = 0xd1
	SET_2_D  = 0xd2
	SET_2_E  = 0xd3
	SET_2_H  = 0xd4
	SET_2_L  = 0xd5
	SET_2_HL = 0xd6
	SET_2_A  = 0xd7
	SET_3_B  = 0xd8
	SET_3_C  = 0xd9
	SET_3_D  = 0xda
	SET_3_E  = 0xdb
	SET_3_H  = 0xdc
	SET_3_L  = 0xdd
	SET_3_HL = 0xde
	SET_3_A  = 0xdf
	SET_4_B  = 0xe0
	SET_4_C  = 0xe1
	SET_4_D  = 0xe2
	SET_4_E  = 0xe3
	SET_4_H  = 0xe4
	SET_4_L  = 0xe5
	SET_4_HL = 0xe6
	SET_4_A  = 0xe7
	SET_5_B  = 0xe8
	SET_5_C  = 0xe9
	SET_5_D  = 0xea
	SET_5_E  = 0xeb
	SET_5_H  = 0xec
	SET_5_L  = 0xed
	SET_5_HL = 0xee
	SET_5_A  = 0xef
	SET_6_B  = 0xf0
	SET_6_C  = 0xf1
	SET_6_D  = 0xf2
	SET_6_E  = 0xf3
	SET_6_H  = 0xf4
	SET_6_L  = 0xf5
	SET_6_HL = 0xf6
	SET_6_A  = 0xf7
	SET_7_B  = 0xf8
	SET_7_C  = 0xf9
	SET_7_D  = 0xfa
	SET_7_E  = 0xfb
	SET_7_H  = 0xfc
	SET_7_L  = 0xfd
	SET_7_HL = 0xfe
	SET_7_A  = 0xff
)
