package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tmsmr/cm4iofan"
)

func usage() {
	usage := `Set PWM duty cycle in %
	$ fanctl set VALUE

Get PWM duty cycle in %
	$ fanctl get

Get current RPM
	$ fanctl rpm
	
	0: PWM duty cycle is 0%
	-1: unable to determine RPM
	1-n: measured/calculated RPM`

	fmt.Println(usage)
	os.Exit(1)
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
		err = ctrl.SetDutyCycle(dc)
		if err != nil {
			panic(err)
		}

	case "get":
		dc, err := ctrl.GetDutyCycle()
		if err != nil {
			panic(err)
		}
		fmt.Println(dc)

	case "rpm":
		rpm, err := ctrl.GetRPM()
		if err != nil {
			panic(err)
		}
		if rpm.Stopped {
			fmt.Println(0)
			return
		}
		if rpm.Undef {
			fmt.Println(-1)
			return
		}
		fmt.Println(rpm.Rpm)

	default:
		usage()
	}
}
