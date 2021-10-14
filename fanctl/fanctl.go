package main

import (
	"github.com/tmsmr/cm4iofan"
	"os"
	"strconv"
)

func main() {
	dc, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	ctrl, err := cm4iofan.New()
	if err != nil {
		panic(err)
	}
	err = ctrl.SetDutyCycle(uint8(dc))
	if err != nil {
		panic(err)
	}
}
