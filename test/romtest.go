package test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/gameboy"
)

func ExecuteROMTest(gb *gameboy.GameBoy, rom string, t *testing.T, check func() int) []string {
	output := ""
	totalTime := 0
	ticker := time.NewTicker(time.Second).C

	config.Get().Emulation.SkipBoot = true
	gb.InsertCartridgeFromPath(fmt.Sprintf("../%s", rom))

	var wg sync.WaitGroup
	wg.Add(1)
	stop := make(chan int)
	go func() {
		gb.PowerOn(stop)
		wg.Done()
	}()

LOOP:
	for totalTime < 60 {
		select {
		case <-ticker:
			if c := check(); c == 1 {
				break LOOP
			} else if c == -1 {
				t.Error("Failed")
				break LOOP
			}
			totalTime += 1
		}
	}

	if totalTime == 60 {
		t.Errorf("Failed (too many tries)")
	}

	stop <- 1
	wg.Wait()

	lines := strings.Split(output, "\n")

	for i := 2; i < len(lines); i++ {
		if strings.Contains(lines[i], "Failed") {
			t.Errorf(lines[i])
		}
	}

	return lines
}

func ExecuteGBROMTest(rom string, t *testing.T) []string {
	output := ""

	config.Get().Emulation.SkipBoot = true
	gb, err := gameboy.NewSerialTextGameBoy(func(v uint8) {
		output += fmt.Sprintf("%c", v)
	})
	if err != nil {
		t.Fatal(err)
	}

	ExecuteROMTest(gb, rom, t, func() int {
		if strings.Contains(output, "Failed") {
			return -1
		} else if strings.Contains(output, "Passed") {
			return 1
		}
		return 0
	})

	lines := strings.Split(output, "\n")

	for i := 2; i < len(lines); i++ {
		if strings.Contains(lines[i], "Failed") {
			t.Errorf(lines[i])
		}
	}

	return lines
}

func ExecuteMooneyeROMTest(rom string, t *testing.T) {
	config.Get().Emulation.SkipBoot = true
	gb, err := gameboy.NewGameBoy()
	if err != nil {
		t.Fatal(err)
	}
	cpu := gb.GetCPU()

	ExecuteROMTest(gb, rom, t, func() int {
		if cpu.B.Get() == 3 && cpu.C.Get() == 5 && cpu.D.Get() == 8 &&
			cpu.E.Get() == 13 && cpu.H.Get() == 21 && cpu.L.Get() == 34 {
			return 1
		} else if cpu.D.Get() == 0x42 {
			return -1
		}
		return 0
	})
}
