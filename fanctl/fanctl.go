package main

import (
	"fmt"
	"github.com/tmsmr/cm4iofan"
	"os"
	"strconv"
)

func usage() {
	usage := `Set PWM duty cycle in %
	$ fanctl set VALUE
Get PWM duty cycle in %
	$ fanctl get
Get current (guessed) RPM:
	$ fanctl rpm`
	fmt.Println(usage)
	os.Exit(1)
}

func set(ctrl *cm4iofan.EMC2301, dc int) {
	err := ctrl.SetDutyCycle(dc)
	if err != nil {
		panic(err)
	}
}

func get(ctrl *cm4iofan.EMC2301) {
	dc, err := ctrl.GetDutyCycle()
	if err != nil {
		panic(err)
	}
	fmt.Println(dc)
}

func rpm(ctrl *cm4iofan.EMC2301) {
	rpm, err := ctrl.GetRPM()
	if err != nil {
		panic(err)
	}
	fmt.Println(rpm)
}

func main() {
	ctrl, err := cm4iofan.New()
	if err != nil {
		panic(err)
	}
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "set":
		if len(os.Args) < 3 {
			usage()
		}
		dc, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		set(ctrl, dc)
	case "get":
		get(ctrl)
	case "rpm":
		rpm(ctrl)
	default:
		usage()
	}
}
