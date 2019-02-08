package cartridge

// SGB flags
const (
	NoSGBSupport = 0 // Normal Gameboy or CGB only game
	SGBSupport   = 3 // Game supports SGB functions
)

// Ram sizes
const (
	RAMNone = 0
	RAM2k   = 1
	RAM8k   = 2
	RAM32k  = 3
	RAM128k = 4
	RAM64k  = 5
)

// Destination codes
const (
	Japanese    = 0
	NonJapanese = 1
)

// Cartridge types
const (
	ROMOnlyType                    = 0x00
	MBC1Type                       = 0x01
	MBC1RAMType                    = 0x02
	MBC1RAMBatteryType             = 0x03
	MBC2Type                       = 0x05
	MBC2BatteryType                = 0x06
	ROMRAMType                     = 0x08
	ROMRAMBatteryType              = 0x09
	MMM01Type                      = 0x0b
	MMM01RAMType                   = 0x0c
	MMM01RAMBatteryType            = 0x0d
	MBC3TimerBatteryType           = 0x0f
	MBC3TimerRAMBatteryType        = 0x10
	MBC3Type                       = 0x11
	MBC3RAMType                    = 0x12
	MBC3RamBatteryType             = 0x13
	MBC5Type                       = 0x19
	MBC5RAMType                    = 0x1a
	MBC5RAMBatteryType             = 0x1b
	MBC5RumbleType                 = 0x1c
	MBC5RumbleRAMType              = 0x1d
	MBC5RumbleRAMBatteryType       = 0x1e
	MBC6Type                       = 0x20
	MBC7SendorRumbleRamBatteryType = 0x22
	PocketCameraType               = 0xfc
	BandaiTama5Type                = 0xfd
	HuC3Type                       = 0xfe
	HuC1RAMBatteryType             = 0xff
)

// LicenseeCodes maps a licensee code to its name
var LicenseeCodes map[uint8]string

func init() {
	LicenseeCodes = map[uint8]string{
		0:   "none",
		1:   "Nintendo R&D1",
		8:   "Capcom",
		13:  "Electronic Arts",
		18:  "Hudson Soft",
		19:  "b-ai",
		20:  "kss",
		22:  "pow",
		24:  "PCM Complete",
		25:  "san-x",
		28:  "Kemco Japan",
		29:  "seta",
		30:  "Viacom",
		31:  "Nintendo",
		32:  "Bandai",
		33:  "Ocean/Acclaim",
		34:  "Konami",
		35:  "Hector",
		37:  "Taito",
		38:  "Hudson",
		39:  "Banpresto",
		41:  "Ubi Soft",
		42:  "Atlus",
		44:  "Malibu",
		46:  "angel",
		47:  "Bullet-Proof",
		49:  "irem",
		50:  "Absolute",
		51:  "Acclaim",
		52:  "Activision",
		53:  "American sammy",
		54:  "Konami",
		55:  "Hi tech entertainment",
		56:  "LJN",
		57:  "Matchbox",
		58:  "Mattel",
		59:  "Milton Bradley",
		60:  "Titus",
		61:  "Virgin",
		64:  "LucasArts",
		67:  "Ocean",
		69:  "Electronic Arts",
		70:  "Infogrames",
		71:  "Interplay",
		72:  "Broderbund",
		73:  "sculptured",
		75:  "sci",
		78:  "THQ",
		79:  "Accolade",
		80:  "misawa",
		83:  "lozc",
		86:  "tokuma shoten i*",
		87:  "tsukuda ori*",
		91:  "Chunsoft",
		92:  "Video system",
		93:  "Ocean/Acclaim",
		95:  "Varie",
		96:  "Yonezawa/s'pal",
		97:  "Kaneko",
		99:  "Pack in soft",
		154: "Konami (Yu-Gi-Oh!)",
	}
}
