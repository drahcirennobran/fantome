package main

import (
	"flag"
	"fmt"
	"os/user"
	"runtime"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	gpioPin := flag.Int("gpio-pin", 18, "GPIO pin")
	brightness := flag.Int("brightness", 64, "Brightness (0-255)")

	flag.Parse()

	user, err := user.Current()
	checkError(err)

	if runtime.GOARCH == "arm" && user.Uid != "0" {
		fmt.Println("This test requires root privilege")
		fmt.Println("Please try \"sudo go test -v\"")
		checkError(err)
	}

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = *brightness
	opt.Channels[0].GpioPin = *gpioPin

	ws, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	err = ws.Init()
	checkError(err)

	for k := 0; k < 5; k++ {
		for j := 0; j < 5; j++ {
			coincoin(ws, 0x0000ff, 0x000000)
			time.Sleep(20 * time.Millisecond)
			coincoin(ws, 0x000000, 0x0000)
			time.Sleep(20 * time.Millisecond)
		}
		for j := 0; j < 5; j++ {
			coincoin(ws, 0x000000, 0xff0000)
			time.Sleep(20 * time.Millisecond)
			coincoin(ws, 0x000000, 0x0000)
			time.Sleep(20 * time.Millisecond)
		}
	}
	coincoin(ws, 0x000000ff, 0x0000ff)
	for j := 0; j < 255; j++ {
		ws.SetBrightness(0, j)
	}
	for j := 255; j >= 0; j-- {
		ws.SetBrightness(0, j)
	}
	coincoin(ws, 0x000000, 0x0000)

	ws.Fini()
}
func coincoin(ws *ws2811.WS2811, color1, color2 uint32) {
	bitmap := make([]uint32, 10)
	for i := 0; i < 5; i++ {
		bitmap[i] = color1
	}
	for i := 5; i < 10; i++ {
		bitmap[i] = color2
	}
	copy(ws.Leds(0), bitmap)
	ws.Render()
	ws.Wait()
}
