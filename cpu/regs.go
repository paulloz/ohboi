package cpu

type GetterSetter interface {
	Get(*CPU) uint8
	Set(*CPU, uint8)
}

type GetterSetter16 interface {
	Get(*CPU) uint16
	Set(*CPU, uint16)
}

type Setter interface {
	Set(*CPU, uint8)
}

type Getter interface {
	Get(*CPU) uint8
}

type Setter16 interface {
	Set(*CPU, uint16)
}

type Getter16 interface {
	Get(*CPU) uint16
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

type registerS struct{}

func (r registerS) Get(cpu *CPU) uint8 {
	return cpu.SP.Hi()
}

func (r registerS) Set(cpu *CPU, v uint8) {
	cpu.SP.SetHi(v)
}

var RegisterS = registerS{}

type registerP struct{}

func (r registerP) Get(cpu *CPU) uint8 {
	return cpu.SP.Lo()
}

func (r registerP) Set(cpu *CPU, v uint8) {
	cpu.SP.SetLo(v)
}

var RegisterP = registerP{}

type immediate struct{}

func (i immediate) Get(cpu *CPU) uint8 {
	return cpu.FetchByte()
}

var Immediate = immediate{}

type immediateOperand struct {
	v uint8
}

func (i immediateOperand) Get(cpu *CPU) uint8 {
	return i.v
}

var ImmediateOperand = &immediateOperand{}

type immediate16 struct{}

func (i immediate16) Get(cpu *CPU) uint16 {
	lo := cpu.FetchByte()
	hi := cpu.FetchByte()
	return uint16(hi)<<8 | uint16(lo)
}

var Immediate16 = immediate16{}

type addressC struct{}

func (a addressC) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(0xff00 + uint16(cpu.C.Get()))
}

func (a addressC) Set(cpu *CPU, v uint8) {
	cpu.mem.Write((0xff00 + uint16(cpu.C.Get())), v)
}

var AddressC = addressC{}

type registerAF struct{}

func (r registerAF) Get(cpu *CPU) uint16 {
	return cpu.AF.Get()
}

func (r registerAF) Set(cpu *CPU, v uint16) {
	v = (cpu.AF.Get() & 0x000f) | (v & 0xfff0)
	cpu.AF.Set(v)
}

var RegisterAF = registerAF{}

type registerBC struct{}

func (r registerBC) Get(cpu *CPU) uint16 {
	return cpu.BC.Get()
}

func (r registerBC) Set(cpu *CPU, v uint16) {
	cpu.BC.Set(v)
}

var RegisterBC = registerBC{}

type addressBC struct{}

func (a addressBC) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(cpu.BC.hilo)
}

func (a addressBC) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(cpu.BC.hilo, v)
}

var AddressBC = addressBC{}

type registerDE struct{}

func (r registerDE) Get(cpu *CPU) uint16 {
	return cpu.DE.Get()
}

func (r registerDE) Set(cpu *CPU, v uint16) {
	cpu.DE.Set(v)
}

var RegisterDE = registerDE{}

type addressDE struct{}

func (a addressDE) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(cpu.DE.hilo)
}

func (a addressDE) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(cpu.DE.hilo, v)
}

var AddressDE = addressDE{}

type registerHL struct{}

func (r registerHL) Get(cpu *CPU) uint16 {
	return cpu.HL.Get()
}

func (r registerHL) Set(cpu *CPU, v uint16) {
	cpu.HL.Set(v)
}

var RegisterHL = registerHL{}

type addressHL struct{}

func (a addressHL) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(cpu.HL.hilo)
}

func (a addressHL) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(cpu.HL.hilo, v)
}

var AddressHL = addressHL{}

type addressHLDec struct{}

func (a addressHLDec) Get(cpu *CPU) uint8 {
	value := cpu.mem.Read(cpu.HL.hilo)
	cpu.HL.hilo--
	return value
}

func (a addressHLDec) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(cpu.HL.hilo, v)
	cpu.HL.hilo--
}

var AddressHLDec = addressHLDec{}

type addressHLInc struct{}

func (a addressHLInc) Get(cpu *CPU) uint8 {
	value := cpu.mem.Read(cpu.HL.hilo)
	cpu.HL.hilo++
	return value
}

func (a addressHLInc) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(cpu.HL.hilo, v)
	cpu.HL.hilo++
}

var AddressHLInc = addressHLInc{}

type addressImmediate struct{}

func (i addressImmediate) Get(cpu *CPU) uint8 {
	lo := cpu.FetchByte()
	hi := cpu.FetchByte()
	return cpu.mem.Read(uint16(hi)<<8 | uint16(lo))
}

func (i addressImmediate) Set(cpu *CPU, v uint8) {
	lo := cpu.FetchByte()
	hi := cpu.FetchByte()
	cpu.mem.Write(uint16(hi)<<8|uint16(lo), v)
}

var AddressImmediate = addressImmediate{}

type addressImmediateOperand struct {
	lo, hi uint8
}

func (o *addressImmediateOperand) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(uint16(o.hi)<<8 | uint16(o.lo))
}

func (o *addressImmediateOperand) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(uint16(o.hi)<<8|uint16(o.lo), v)
}

func (o *addressImmediateOperand) Lo() GetterSetter {
	return o
}

func (o *addressImmediateOperand) Hi() GetterSetter {
	return &addressHighImmediateOperand{
		hi: &o.hi,
		lo: &o.lo,
	}
}

var AddressImmediateOperand = &addressImmediateOperand{}

type addressHighImmediateOperand struct {
	hi, lo *uint8
}

func (o *addressHighImmediateOperand) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(uint16(*o.hi)<<8 | uint16(*o.lo) + 1)
}

func (o *addressHighImmediateOperand) Set(cpu *CPU, v uint8) {
	cpu.mem.Write(uint16(*o.hi)<<8|uint16(*o.lo)+1, v)
}

type addressImmediate16 struct{}

func (i addressImmediate16) Get(cpu *CPU) uint16 {
	return cpu.mem.ReadWord(cpu.FetchWord())
}

func (i addressImmediate16) Set(cpu *CPU, v uint16) {
	lo := cpu.FetchByte()
	hi := cpu.FetchByte()
	cpu.mem.WriteWord(uint16(hi)<<8|uint16(lo), v)
}

var AddressImmediate16 = addressImmediate16{}

type addressFF00N struct{}

func (a addressFF00N) Get(cpu *CPU) uint8 {
	return cpu.mem.Read(0xff00 + uint16(cpu.FetchByte()))
}

func (a addressFF00N) Set(cpu *CPU, v uint8) {
	cpu.mem.Write((0xff00 + uint16(cpu.FetchByte())), v)
}

var AddressFF00N = addressFF00N{}

type addressFF00NOperand struct {
	v uint8
}

func (a addressFF00NOperand) Get(cpu *CPU) uint8 {
	return cpu.mem.Read((0xff00 + uint16(a.v)))
}

func (a addressFF00NOperand) Set(cpu *CPU, v uint8) {
	cpu.mem.Write((0xff00 + uint16(a.v)), v)
}

var AddressFF00NOperand = &addressFF00NOperand{}

type registerSP struct{}

func (r registerSP) Get(cpu *CPU) uint16 {
	return cpu.SP.Get()
}

func (r registerSP) Set(cpu *CPU, v uint16) {
	cpu.SP.Set(v)
}

var RegisterSP = registerSP{}

type registerPC struct{}

func (r registerPC) Get(cpu *CPU) uint16 {
	return cpu.PC
}

func (r registerPC) Set(cpu *CPU, v uint16) {
	cpu.PC = v
}

var RegisterPC = registerPC{}
