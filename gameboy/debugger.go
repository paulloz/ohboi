// +build DEBUG

package gameboy

import (
	"fmt"
	"os"
	"time"

	"io/ioutil"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

type tDebugger struct {
	uifps       *widgets.Paragraph
	uiregisters *widgets.Paragraph
	uinext      *widgets.Paragraph
	uistack     *widgets.Paragraph
	uiio        *widgets.Paragraph

	stepByStep bool

	stepper chan int

	gb *GameBoy
}

var (
	debugger *tDebugger
)

func init() {
	ui.Init()

	debugger = &tDebugger{
		uifps:       widgets.NewParagraph(),
		uiregisters: widgets.NewParagraph(),
		uinext:      widgets.NewParagraph(),
		uistack:     widgets.NewParagraph(),
		uiio:        widgets.NewParagraph(),

		stepByStep: true,
		stepper:    make(chan int),
	}
}

func (debugger *tDebugger) start(gb *GameBoy) {
	debugger.gb = gb

	go func() {
		debugger.uifps.Title = "FPS"
		debugger.uifps.SetRect(0, 0, 10, 3)

		debugger.uiregisters.Title = "CPU Registers"
		debugger.uiregisters.SetRect(10, 0, 30, 12)

		debugger.uinext.SetRect(30, 0, 70, 12)

		debugger.uistack.Title = "Stack"
		debugger.uistack.SetRect(30, 12, 50, 36)

		debugger.uiio.SetRect(10, 12, 30, 36)

		ui.Render(debugger.uifps)
		ui.Render(debugger.uiregisters)
		ui.Render(debugger.uinext)
		ui.Render(debugger.uistack)
		ui.Render(debugger.uiio)

		ticker := time.NewTicker(time.Second / FPS).C
		uiEvents := ui.PollEvents()

		for {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>", "<Escape>":
					debugger.close()
					os.Exit(0)
				case "<Space>":
					if debugger.stepByStep {
						debugger.stepper <- 1
					}
				case "p":
					debugger.stepByStep = !debugger.stepByStep
					if !debugger.stepByStep {
						debugger.stepper <- 1
					}
				}
			case <-ticker:
				pc := gb.cpu.PC

				// uifps.Text = fmt.Sprintf("%02d", fps)
				debugger.uiregisters.Text = gb.cpu.Dump()

				name := ""
				if opCodeNames[gb.Memory.Read(pc)] != "" {
					name = opCodeNames[gb.Memory.Read(pc)]
				}
				next := fmt.Sprintf(" -> 0x%02x %s\n", gb.Memory.Read(pc), name)
				for i := uint16(1); i < 10; i++ {
					name = ""
					if opCodeNames[gb.Memory.Read(pc+i)] != "" {
						name = opCodeNames[gb.Memory.Read(pc+i)]
					}
					next += fmt.Sprintf("    0x%02x %s\n", gb.Memory.Read(pc+i), name)
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

				io := fmt.Sprintf("\nDIV : 0x%02x\n", gb.Memory.Read(0xff04))
				io += fmt.Sprintf("IE  : 0x%02x\n", gb.Memory.Read(0xffff))
				io += fmt.Sprintf("IF  : 0x%02x\n", gb.Memory.Read(0xff0f))
				io += fmt.Sprintf("TIMA: 0x%02x\n", gb.Memory.Read(0xff05))
				io += fmt.Sprintf("TAC : 0x%02x\n", gb.Memory.Read(0xff07))
				io += fmt.Sprintf("\nClock : %d\n", gb.GETCLOCK())
				debugger.uiio.Text = io

				ui.Render(debugger.uifps)
				ui.Render(debugger.uiregisters)
				ui.Render(debugger.uinext)
				ui.Render(debugger.uistack)
				ui.Render(debugger.uiio)
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
	if debugger.stepByStep {
		<-debugger.stepper
	}
}
