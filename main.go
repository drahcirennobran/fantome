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

const (
	noir  uint32 = 0x000000
	rouge uint32 = 0xff0000
	bleu  uint32 = 0x0000ff
	vert  uint32 = 0x00ff00
)

func main() {

	gpioPin := flag.Int("gpio-pin", 18, "GPIO pin")
	brightness := flag.Int("brightness", 128, "Brightness (0-255)")

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

	//var bitRB = []uint32{rouge, rouge, rouge, rouge, rouge, bleu, bleu, bleu, bleu, bleu}
	//var bitBR = []uint32{bleu, bleu, bleu, bleu, bleu, rouge, rouge, rouge, rouge, rouge}
	var bit10N = []uint32{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir}
	var bit5B5N = []uint32{bleu, bleu, bleu, bleu, bleu, noir, noir, noir, noir, noir}
	var bit10B = []uint32{bleu, bleu, bleu, bleu, bleu, bleu, bleu, bleu, bleu, bleu}
	var bit5N5R = []uint32{noir, noir, noir, noir, noir, rouge, rouge, rouge, rouge, rouge}

	for k := 0; k < 5; k++ {
		for j := 0; j < 5; j++ {
			coincoin(ws, bit5B5N, *brightness)
			time.Sleep(20 * time.Millisecond)
			coincoin(ws, bit10N, *brightness)
			time.Sleep(20 * time.Millisecond)
		}
		for j := 0; j < 5; j++ {
			coincoin(ws, bit5N5R, *brightness)
			time.Sleep(20 * time.Millisecond)
			coincoin(ws, bit10N, *brightness)
			time.Sleep(20 * time.Millisecond)
		}
	}

	for k := 0; k < 5; k++ {
		for j := 0; j < 255; j++ {
			coincoin(ws, bit10B, j)
		}
	}
	bitmapScroll1 := [][]uint32{
		{bleu, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, bleu, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, bleu, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, bleu, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, bleu, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, bleu, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, bleu, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, bleu, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, bleu, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, bleu}}

	bitmapScroll2 := [][]uint32{
		{bleu, bleu, bleu, bleu, bleu, bleu, bleu, bleu, bleu, bleu},
		{bleu, bleu, bleu, bleu, rouge, bleu, bleu, bleu, bleu, bleu},
		{bleu, bleu, bleu, bleu, rouge, rouge, bleu, bleu, bleu, bleu},
		{bleu, bleu, bleu, rouge, rouge, rouge, bleu, bleu, bleu, bleu},
		{bleu, bleu, rouge, rouge, rouge, rouge, bleu, bleu, bleu, bleu},
		{bleu, bleu, rouge, rouge, rouge, rouge, rouge, bleu, bleu, bleu},
		{bleu, rouge, rouge, rouge, rouge, rouge, rouge, bleu, bleu, bleu},
		{bleu, rouge, rouge, rouge, rouge, rouge, rouge, rouge, bleu, bleu},
		{rouge, rouge, rouge, rouge, rouge, rouge, rouge, rouge, rouge, rouge},
		{rouge, rouge, rouge, rouge, rouge, rouge, rouge, rouge, rouge, rouge}}
	bitmapPong(ws, bitmapScroll1, *brightness, 5)
	bitmapPong(ws, bitmapScroll2, *brightness, 5)
	coincoin(ws, bit10N, 0)
	ws.Fini()
}
func bitmapPong(ws *ws2811.WS2811, bitmap [][]uint32, brightness int, repeat int) {
	for k := 0; k < 10; k++ {
		for j := 0; j < 10; j++ {
			coincoin(ws, bitmap[j], brightness)
			time.Sleep(20 * time.Millisecond)
		}
		for j := 8; j >= 0; j-- {
			coincoin(ws, bitmap[j], brightness)
			time.Sleep(20 * time.Millisecond)
		}
	}

}
func coincoin(ws *ws2811.WS2811, bitmap []uint32, b int) {
	ws.SetBrightness(0, b)

	copy(ws.Leds(0), bitmap)
	ws.Render()
	ws.Wait()
}
