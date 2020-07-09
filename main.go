package main

import (
	"flag"
	"fmt"
	"os/user"
	"runtime"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	//pixelColor = 255 << 16 // red
	pixelColor = 0x0000ff
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	gpioPin := flag.Int("gpio-pin", 18, "GPIO pin")
	width := flag.Int("width", 10, "LED matrix width")
	height := flag.Int("height", 1, "LED matrix height")
	brightness := flag.Int("brightness", 64, "Brightness (0-255)")

	flag.Parse()

	user, err := user.Current()
	checkError(err)

	if runtime.GOARCH == "arm" && user.Uid != "0" {
		fmt.Println("This test requires root privilege")
		fmt.Println("Please try \"sudo go test -v\"")
		checkError(err)
	}

	size := *width * *height
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = *brightness
	opt.Channels[0].LedCount = size
	opt.Channels[0].GpioPin = *gpioPin

	ws, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	err = ws.Init()
	checkError(err)

	bitmap := make([]uint32, size)

	for i := 0; i < size; i++ {
		if i > 0 {
			bitmap[i-1] = 0
		}

		bitmap[i] = pixelColor
		copy(ws.Leds(0), bitmap)
		ws.Render()
		ws.Wait()
	}

	for i := 0; i < len(ws.Leds(0)); i++ {
		ws.Leds(0)[i] = 0
	}

	ws.Fini()
}
