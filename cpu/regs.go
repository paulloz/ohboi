package cpu

type Setter interface {
	Set(*CPU, uint8)
}

type Getter interface {
	Get(*CPU) uint8
}

type registerA struct{}

func (r registerA) Get(cpu *CPU) uint8 {
	return cpu.A.Get()
}

func (r registerA) Set(cpu *CPU, v uint8) {
	cpu.A.Set(v)
}

var RegisterA = registerA{}

type registerB struct{}

func (r registerB) Get(cpu *CPU) uint8 {
	return cpu.B.Get()
}

func (r registerB) Set(cpu *CPU, v uint8) {
	cpu.B.Set(v)
}

var RegisterB = registerB{}

type registerC struct{}

func (r registerC) Get(cpu *CPU) uint8 {
	return cpu.C.Get()
}

func (r registerC) Set(cpu *CPU, v uint8) {
	cpu.C.Set(v)
}

var RegisterC = registerC{}

type registerD struct{}

func (r registerD) Get(cpu *CPU) uint8 {
	return cpu.D.Get()
}

func (r registerD) Set(cpu *CPU, v uint8) {
	cpu.D.Set(v)
}

var RegisterD = registerD{}

type registerE struct{}

func (r registerE) Get(cpu *CPU) uint8 {
	return cpu.E.Get()
}

func (r registerE) Set(cpu *CPU, v uint8) {
	cpu.E.Set(v)
}

var RegisterE = registerE{}

type registerF struct{}

func (r registerF) Get(cpu *CPU) uint8 {
	return cpu.F.Get()
}

func (r registerF) Set(cpu *CPU, v uint8) {
	cpu.F.Set(v)
}

var RegisterF = registerF{}

type registerH struct{}

func (r registerH) Get(cpu *CPU) uint8 {
	return cpu.H.Get()
}

func (r registerH) Set(cpu *CPU, v uint8) {
	cpu.H.Set(v)
}

var RegisterH = registerH{}

type registerL struct{}

func (r registerL) Get(cpu *CPU) uint8 {
	return cpu.L.Get()
}

func (r registerL) Set(cpu *CPU, v uint8) {
	cpu.L.Set(v)
}

var RegisterL = registerL{}

type immediate struct{}

func (i immediate) Get(cpu *CPU) uint8 {
	return cpu.FetchByte()
}

var Immediate = immediate{}

type addressHL struct{}

func (a addressHL) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(cpu.HL.hilo)
}

func (a addressHL) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(cpu.HL.hilo, v)
}

var AddressHL = addressHL{}
