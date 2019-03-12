// +build DEBUG

package gameboy

import (
	"fmt"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"github.com/Knetic/govaluate"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/paulloz/ohboi/cartridge"
	"github.com/paulloz/ohboi/consts"
)

type Breakpoint struct {
	Condition *govaluate.EvaluableExpression
}

type tDebugger struct {
	uifps       *widgets.Paragraph
	uiregisters *widgets.Paragraph
	uinext      *widgets.Paragraph
	uistack     *widgets.Paragraph
	uiio        *widgets.Paragraph
	uimem       *widgets.Paragraph
	uiinput     *widgets.Paragraph
	uicartridge *widgets.Paragraph
	uitabpane   *widgets.TabPane
	breakpoints []Breakpoint

	memoryPos  int
	stepByStep bool
	stepPC     int

	stepper chan int

	gb *GameBoy
}

var (
	debugger *tDebugger
)

func init() {
	if err := ui.Init(); err != nil {
		panic(err)
	}

	debugger = &tDebugger{
		uifps:       widgets.NewParagraph(),
		uiregisters: widgets.NewParagraph(),
		uinext:      widgets.NewParagraph(),
		uistack:     widgets.NewParagraph(),
		uiio:        widgets.NewParagraph(),
		uimem:       widgets.NewParagraph(),
		uiinput:     widgets.NewParagraph(),
		uicartridge: widgets.NewParagraph(),
		uitabpane:   widgets.NewTabPane(),

		stepByStep: true,
		stepper:    make(chan int),
	}
}

func (debugger *tDebugger) start(gb *GameBoy) {
	debugger.gb = gb

	instructionSize := func(opcodeName string) int {
		if strings.Contains(opcodeName, ",nn") || strings.Contains(opcodeName, " nn") {
			return 3
		} else if strings.Contains(opcodeName, ",n") || strings.Contains(opcodeName, "ff00n") {
			return 2
		}
		return 1
	}

	go func() {
		debugger.uitabpane = widgets.NewTabPane("Inspect", "Shell")
		debugger.uitabpane.SetRect(0, 0, 50, 3)
		debugger.uitabpane.Border = true

		debugger.uifps.Title = "FPS"
		debugger.uifps.SetRect(0, 3, 10, 6)

		debugger.uiregisters.Title = "CPU Registers"
		debugger.uiregisters.SetRect(10, 3, 30, 12)

		debugger.uinext.Title = "Next instructions"
		debugger.uinext.SetRect(30, 3, 68, 12)

		debugger.uicartridge.Title = "Cartridge"
		debugger.uicartridge.SetRect(10, 12, 68, 17)

		debugger.uistack.Title = "Stack"
		debugger.uistack.SetRect(23, 17, 43, 36)

		debugger.uiio.Title = "IO"
		debugger.uiio.SetRect(10, 17, 23, 36)

		debugger.uimem.Title = "Memory"
		debugger.uimem.SetRect(43, 17, 68, 36)

		debugger.uiinput.Title = "Command"
		debugger.uiinput.SetRect(0, 36, 70, 39)

		renderTab := func() {
			switch debugger.uitabpane.ActiveTabIndex {
			case 0:
				ui.Render(debugger.uifps)
				ui.Render(debugger.uiregisters)
				ui.Render(debugger.uinext)
				ui.Render(debugger.uistack)
				ui.Render(debugger.uiio)
				ui.Render(debugger.uimem)
				ui.Render(debugger.uicartridge)
			case 1:
				ui.Render(debugger.uiinput)
			}
		}

		ticker := time.NewTicker(time.Second / consts.FPS).C
		uiEvents := ui.PollEvents()

		for {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>", "<Escape>":
					debugger.close()
					os.Exit(0)
				case "<Tab>":
					debugger.uitabpane.ActiveTabIndex = (debugger.uitabpane.ActiveTabIndex + 1) % 2
					ui.Clear()
					ui.Render(debugger.uitabpane)
					renderTab()
				case "<Space>":
					if debugger.uitabpane.ActiveTabIndex == 1 {
						if debugger.stepByStep {
							debugger.stepper <- 1
						}
					}
				case "<PageUp>":
					debugger.memoryPos -= 0xa0
					if debugger.memoryPos < 0 {
						debugger.memoryPos = 65536 - 0xa0
					}
				case "<PageDown>":
					debugger.memoryPos += 0xa0
					if debugger.memoryPos > 65535 {
						debugger.memoryPos = 0
					}
				case "p":
					debugger.stepByStep = !debugger.stepByStep
					if !debugger.stepByStep {
						debugger.stepper <- 1
					}
					debugger.stepPC = -1
				case "s":
					debugger.stepByStep = true
					debugger.stepPC = -1
					debugger.stepper <- 1
				case "n":
					opcode := debugger.gb.Memory.Read(debugger.gb.cpu.PC)
					opcodeName := opCodeNames[opcode]
					debugger.stepPC = int(debugger.gb.cpu.PC) + instructionSize(opcodeName)
					debugger.stepByStep = false
					debugger.stepper <- 1
				}
			case <-ticker:
				pc := gb.cpu.PC

				// uifps.Text = fmt.Sprintf("%02d", fps)
				debugger.uiregisters.Text = gb.cpu.Dump()

				name := ""
				if opCodeNames[gb.Memory.Read(pc)] != "" {
					name = opCodeNames[gb.Memory.Read(pc)]
				}

				formatOpcode := func(name string) (int, string) {
					size := instructionSize(name)
					value := ""
					if size == 2 {
						value = fmt.Sprintf(" (%02x)", uint16(gb.Memory.Read(pc+1)))
					} else if size == 3 {
						value = fmt.Sprintf(" (%02x)", uint16(gb.Memory.ReadWord(pc+1)))
					}
					return size, fmt.Sprintf("0x%02x %s%s", gb.Memory.Read(pc), name, value)
				}

				size, s := formatOpcode(name)
				next := fmt.Sprintf(" -> %s\n", s)
				pc += uint16(size)
				for i := 1; i < 10; i++ {
					name = opCodeNames[gb.Memory.Read(pc)]
					size, s := formatOpcode(name)
					next += fmt.Sprintf("    %s\n", s)
					pc += uint16(size)
				}
				debugger.uinext.Text = next

				stack := fmt.Sprintf("\nSP: 0x%04x\n\n", gb.cpu.SP.Get())
				for i := gb.cpu.SP.Get(); i <= 0xfffe; i++ {
					prefix := "   "
					if i == gb.cpu.SP.Get() {
						prefix = " =>"
					}
					stack += fmt.Sprintf("%s (0x%04x) 0x%02x\n", prefix, i, gb.Memory.Read(i))
				}
				debugger.uistack.Text = stack

				mem := ""
				addr := debugger.memoryPos
				for i := 0; i < 22 && addr < 65556; i++ {
					mem += fmt.Sprintf(" %04x:", addr)
					for j := 0; j < 8 && addr+j < 65536; j++ {
						mem += fmt.Sprintf("%02x", debugger.gb.Memory.Read(uint16(addr+j)))
					}
					mem += "\n"
					addr += 8
				}
				debugger.uimem.Text = mem

				cr := debugger.gb.Memory.Cartridge()
				if mbc1, ok := cr.MBC.(*cartridge.MBC1); ok {
					mode := mbc1.GetMode()
					ramEnabled, ram, ramBank, activeRAMBankStart := mbc1.GetRAMState()
					text := fmt.Sprintf(" RAM banking: %v, RAM %d\n RAM enabled %+v, RAM bank %d, RAM address %X\n", mode, ram, ramEnabled, ramBank, activeRAMBankStart)
					romBank, activeROMBankStart := mbc1.GetROMState()
					text += fmt.Sprintf(" ROM bank: %d, ROM address: %X", romBank, activeROMBankStart)
					debugger.uicartridge.Text = text
				}

				io := fmt.Sprintf("\nDIV : 0x%02x\n", gb.Memory.Read(0xff04))
				io += fmt.Sprintf("IE  : 0x%02x\n", gb.Memory.Read(0xffff))
				io += fmt.Sprintf("IF  : 0x%02x\n", gb.Memory.Read(0xff0f))
				io += fmt.Sprintf("TIMA: 0x%02x\n", gb.Memory.Read(0xff05))
				io += fmt.Sprintf("TAC : 0x%02x\n", gb.Memory.Read(0xff07))
				io += fmt.Sprintf("\nClock : %d\n", gb.GETCLOCK())
				debugger.uiio.Text = io

				ui.Render(debugger.uitabpane)
				renderTab()
			}
		}
	}()
}

func (debugger *tDebugger) close() {
	ui.Close()
}

func (debugger *tDebugger) hexDump(data []uint8, addrPrefix uint16) {
	// fmt.Println("NOW")
	str := time.Now().Format("15:04:05")
	for a, v := range data {
		if a%0x10 == 0 {
			str += fmt.Sprintf("\n%08x   ", addrPrefix+(uint16(a)))
		}
		str += fmt.Sprintf(" %02x", v)
	}
	ioutil.WriteFile("vram.txt", []byte(str), 0644)
}

func debuggerStart(gb *GameBoy) {
	debugger.start(gb)
}

func debuggerStop() {
	debugger.close()
}

func debuggerStep() {
	cpu := debugger.gb.cpu
	if debugger.stepByStep {
		<-debugger.stepper
	} else if debugger.stepPC == int(cpu.PC) {
		<-debugger.stepper
	} else {
		parameters := make(map[string]interface{})
		if len(debugger.breakpoints) > 0 {
			parameters["PC"] = cpu.PC
		}

		for _, breakpoint := range debugger.breakpoints {
			result, err := breakpoint.Condition.Evaluate(parameters)
			if err != nil {
				fmt.Printf("Error while executing condition %s", err)
				continue
			}

			if result == true {
				debugger.stepByStep = true
				<-debugger.stepper
				return
			}
		}
	}
}

func newBreakpoint(condition string) (Breakpoint, error) {
	expression, err := govaluate.NewEvaluableExpression(condition)
	if err != nil {
		return Breakpoint{}, err
	}

	return Breakpoint{expression}, nil
}

func AddBreakpoint(condition string) error {
	breakpoint, err := newBreakpoint(condition)
	if err != nil {
		return err
	}

	debugger.breakpoints = append(debugger.breakpoints, breakpoint)
	return nil
}

func StepByStep(state bool) {
	debugger.stepByStep = state
}
