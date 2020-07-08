package main

import (
	"fmt"

	"tinygo.org/x/drivers/ws2812"
)

func main() {

	fmt.Println("coucou")
	Device ruban = ws2812.New(37)
}
