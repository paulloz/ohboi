// +build DEBUG

package gameboy

import (
	"fmt"
	"os"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

type tDebugger struct {
	uifps       *widgets.Paragraph
	uiregisters *widgets.Paragraph
	uinext      *widgets.Paragraph
	uistack     *widgets.Paragraph

	stepByStep bool

	stepper chan int
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

		stepByStep: true,
		stepper:    make(chan int),
	}
}

func (debugger *tDebugger) start(gb *GameBoy) {
	go func() {
		debugger.uifps.Title = "FPS"
		debugger.uifps.SetRect(0, 0, 10, 3)

		debugger.uiregisters.Title = "CPU Registers"
		debugger.uiregisters.SetRect(10, 0, 30, 12)

		debugger.uinext.SetRect(30, 0, 50, 12)

		debugger.uistack.Title = "Stack"
		debugger.uistack.SetRect(30, 12, 50, 36)

		ui.Render(debugger.uifps)
		ui.Render(debugger.uiregisters)
		ui.Render(debugger.uinext)
		ui.Render(debugger.uistack)

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

				next := fmt.Sprintf(" -> 0x%02x\n", gb.memory.Read(pc))
				for i := uint16(1); i < 10; i++ {
					next += fmt.Sprintf("    0x%02x\n", gb.memory.Read(pc+i))
				}
				debugger.uinext.Text = next

				stack := fmt.Sprintf("\nSP: 0x%04x\n\n", gb.cpu.SP.Get())
				for i := gb.cpu.SP.Get(); i <= 0xfffe; i++ {
					prefix := "   "
					if i == gb.cpu.SP.Get() {
						prefix = " =>"
					}
					stack += fmt.Sprintf("%s (0x%04x) 0x%02x\n", prefix, i, gb.memory.Read(i))
				}
				debugger.uistack.Text = stack

				ui.Render(debugger.uifps)
				ui.Render(debugger.uiregisters)
				ui.Render(debugger.uinext)
				ui.Render(debugger.uistack)
			}
		}
	}()
}

func (debugger *tDebugger) close() {
	ui.Close()
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
