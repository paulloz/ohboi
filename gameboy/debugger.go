package gameboy

import (
	"fmt"
	"os"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

const (
	debug = true
)

type tDebugger struct {
	uifps       *widgets.Paragraph
	uiregisters *widgets.Paragraph
	uinext      *widgets.Paragraph

	stepByStep bool

	stepper chan int
}

var (
	debugger *tDebugger
)

func init() {
	if debug {
		ui.Init()

		debugger = &tDebugger{
			uifps:       widgets.NewParagraph(),
			uiregisters: widgets.NewParagraph(),
			uinext:      widgets.NewParagraph(),

			stepByStep: true,
			stepper:    make(chan int),
		}
	}
}

func (debugger *tDebugger) start(gb *GameBoy) {
	go func() {
		debugger.uifps.Title = "FPS"
		debugger.uifps.SetRect(0, 0, 15, 3)

		debugger.uiregisters.Title = "CPU Registers"
		debugger.uiregisters.SetRect(15, 0, 35, 12)

		debugger.uinext.SetRect(35, 0, 55, 12)

		ui.Render(debugger.uifps)
		ui.Render(debugger.uiregisters)
		ui.Render(debugger.uinext)

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

				ui.Render(debugger.uifps)
				ui.Render(debugger.uiregisters)
				ui.Render(debugger.uinext)
			}
		}
	}()
}

func (debugger *tDebugger) close() {
	ui.Close()
}
