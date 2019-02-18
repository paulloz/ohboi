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

func ExecuteROMTest(rom string, t *testing.T) []string {
	output := ""
	totalTime := 0
	ticker := time.NewTicker(time.Second).C

	config.Get().Emulation.SkipBoot = true
	gb := gameboy.NewSerialTextGameBoy(func(v uint8) {
		output += fmt.Sprintf("%c", v)
		if strings.Contains(output, "Failed") || strings.Contains(output, "Passed") {
			totalTime = 1000
		}
	})
	gb.InsertCartridgeFromFile(fmt.Sprintf("../gb-test-roms/%s", rom))

	var wg sync.WaitGroup
	wg.Add(1)
	stop := make(chan int)
	go func() {
		gb.PowerOn(stop)
		wg.Done()
	}()

	for totalTime < 60 {
		select {
		case <-ticker:
			totalTime += 1
		}
	}

	stop <- 1
	wg.Wait()

	lines := strings.Split(output, "\n")

	romPath := strings.Split(rom, "/")
	romName := strings.Split(romPath[len(romPath)-1], ".")[0]
	if romName != lines[0] {
		t.Errorf("Expected rom name to be `%s`, but got `%s`", romName, lines[0])
	}

	for i := 2; i < len(lines); i++ {
		if strings.Contains(lines[i], "Failed") {
			t.Errorf(lines[i])
		}
	}

	return lines
}
