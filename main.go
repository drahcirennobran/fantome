package main

import (
	"flag"
	"fmt"
	"math/rand"
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
	noir   uint32 = 0x000000
	rouge  uint32 = 0xff0000
	bleu   uint32 = 0x0000ff
	vert   uint32 = 0x00ff00
	blanc  uint32 = 0xffffff
	length int    = 13
	height int    = 7
)

func main() {

	gpioPin := flag.Int("gpio-pin", 21, "GPIO pin")
	brightness := flag.Int("brightness", 255, "Brightness (0-255)")
	ledCount := 92

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
	opt.Channels[0].LedCount = ledCount

	ws, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	err = ws.Init()
	checkError(err)

	zigzagFilter := computeZigzagFilter(13, 7)

	gyrophare(ws, zigzagFilter, *brightness)
	fire(ws, zigzagFilter, *brightness, 1000)

	ws.Fini()
}
func fire(ws *ws2811.WS2811, zigzagFilter []uint32, brightness int, nbIter int) {
	matrix := [][]uint32{
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir, noir},
		{rouge, noir, rouge, noir, rouge, noir, rouge, noir, rouge, noir, rouge, noir, noir}}

	for iter := 0; iter < nbIter; iter++ {
		for y := 0; y < height; y++ {
			for x := 1; x < length+1; x++ {
				pix1R := matrix[(y+1)%height][(x-1)%length] >> 16
				pix1G := (matrix[(y+1)%height][(x-1)%length] & 0x00ff00) >> 8
				pix1B := matrix[(y+1)%height][(x-1)%length] & 0x0000ff
				pix2R := matrix[(y+1)%height][(x-1)%length] >> 16
				pix2G := (matrix[(y+1)%height][(x-1)%length] & 0x00ff00) >> 8
				pix2B := matrix[(y+1)%height][(x-1)%length] & 0x0000ff
				pix3R := matrix[(y+1)%height][(x-1)%length] >> 16
				pix3G := (matrix[(y+1)%height][(x-1)%length] & 0x00ff00) >> 8
				pix3B := matrix[(y+1)%height][(x-1)%length] & 0x0000ff

				pixR := (pix1R + pix2R + pix3R) >> 2
				pixG := (pix1G + pix2G + pix3G) >> 2
				pixB := (pix1B + pix2B + pix3B) >> 2

				matrix[y%height][x%length] = pixR<<16 + pixG<<8 + pixB
			}
		}
		redLimit := rand.Intn(100)
		for x := 0; x < length; x++ {
			if rand.Intn(60) > redLimit {
				yellowLimit := rand.Intn(12)
				matrix[height-1][x] = rouge + 0x001100*(uint32)(yellowLimit)
			}
		}
		copy2Led(ws, applyFilter(matrixToLinear(matrix), zigzagFilter), brightness)
		time.Sleep(120 * time.Millisecond)
	}
	copy2Led(ws, blackLine(13*7), brightness)
}

func gyrophare(ws *ws2811.WS2811, zigzagFilter []uint32, brightness int) {
	matrix := [][]uint32{
		{noir, noir, noir, noir, noir, noir, bleu, noir, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, bleu, bleu, bleu, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, bleu, bleu, bleu, bleu, bleu, noir, noir, noir, noir},
		{noir, noir, noir, noir, bleu, bleu, bleu, bleu, bleu, noir, noir, noir, noir},
		{noir, noir, noir, noir, bleu, bleu, bleu, bleu, bleu, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, bleu, bleu, bleu, noir, noir, noir, noir, noir},
		{noir, noir, noir, noir, noir, noir, bleu, noir, noir, noir, noir, noir, noir}}
	line := matrixToLinear(matrix)
	rollXFilter := computeRollXFilter(13, 7)
	for tour := 0; tour < 20; tour++ {
		for i := 0; i < 13; i++ {
			line = applyFilter(line, rollXFilter)

			copy2Led(ws, applyFilter(line, zigzagFilter), brightness)
			time.Sleep(40 * time.Millisecond)
		}
	}
	copy2Led(ws, blackLine(13*7), brightness)

}
func blackLine(size int) []uint32 {
	var leds []uint32 = make([]uint32, size, size)
	for i := range leds {
		leds[i] = noir
	}
	return leds
}
func matrixToLinear(matrix [][]uint32) []uint32 {
	length := len(matrix[0])
	height := len(matrix)
	size := height * length
	var leds []uint32 = make([]uint32, size, size)
	//fmt.Printf("linear size : %v\n", size)
	for y, line := range matrix {
		//fmt.Printf("line %v : %v\n", y, line)
		for x, val := range line {
			//fmt.Printf("(%v,%v) led(%v)= %v\n", x, y, y*length+x, val)
			leds[y*length+x] = val
		}
	}
	return leds
}

func applyFilter(line []uint32, locationFilter []uint32) []uint32 {
	var result []uint32 = make([]uint32, len(line), len(line))
	for i := range result {
		result[i] = line[locationFilter[i]]
	}
	return result
}

func copy2Led(ws *ws2811.WS2811, bitmap []uint32, b int) {
	ws.SetBrightness(0, b)

	copy(ws.Leds(0), bitmap)
	ws.Render()
	ws.Wait()
}

func computeZigzagFilter(length, height int) []uint32 {
	size := height * length
	var leds []uint32 = make([]uint32, size, size)
	for y := 0; y < height; y++ {
		for x := 0; x < length; x += 2 {
			//fmt.Printf("%v*%v+%v-%v-1=%v (%v)| ", x, height, height, y, x*height+height-y-1, x+y*length)
			//leds[x+y*length] = (uint32)(x*height + height - y - 1)
			leds[x*height+height-y-1] = (uint32)(x + y*length)
			if x+1 < length {
				//fmt.Printf("(%v+1)*%v+%v=%v (%v)| ", x, height, y, (x+1)*height+y, x+1+y*length)
				//leds[x+1+y*length] = (uint32)((x+1)*height + y)
				leds[(x+1)*height+y] = (uint32)(x + 1 + y*length)
			}
		}
		//fmt.Println("")
	}
	return leds
}

func computeRevertFilter(size int) []uint32 {
	var revertFilter []uint32 = make([]uint32, size, size)
	for i := range revertFilter {
		revertFilter[i] = (uint32)(len(revertFilter) - i - 1)
	}
	return revertFilter
}

func computeSimpleScrollFilter(size int) []uint32 {
	var simpleScrollFilter []uint32 = make([]uint32, size, size)
	for i := range simpleScrollFilter {
		if i == 0 {
			simpleScrollFilter[0] = (uint32)(size - 1)
		} else {
			simpleScrollFilter[i] = (uint32)(i - 1)
		}
	}
	return simpleScrollFilter
}
func computeRollXFilter(length, height int) []uint32 {
	var filter []uint32 = make([]uint32, length*height, length*height)
	for y := 0; y < height; y++ {
		for x := 0; x < length; x++ {
			filter[y*length+x] = (uint32)(y*length + ((x + 1) % length))
		}
	}
	return filter
}

//var bitRB = []uint32{rouge, rouge, rouge, rouge, rouge, bleu, bleu, bleu, bleu, bleu}
//var bitBR = []uint32{bleu, bleu, bleu, bleu, bleu, rouge, rouge, rouge, rouge, rouge}
/*
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
*/
