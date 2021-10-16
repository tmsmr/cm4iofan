# cm4iofan
*Simple utility to get/set the PWM duty cycle and to measure the RPM for a fan connected to the 4-pin header on the CM4IO.*

## Requirements
- Enabled `i2c_vc` overlay (`dtparam=i2c_vc=on` in `/boot/config.txt`)
- Loaded `i2c-dev` module (e.g. `modprobe i2c-dev` or using `/etc/modules`)
- A user in the `i2c` group (or `root`)

## Using the Go Module

```go
import "github.com/tmsmr/cm4iofan"
 
...

ctrl, err := cm4iofan.New()
if err != nil {
	panic(err)
}
err = ctrl.SetDutyCycle(50)
if err != nil {
	panic(err)
}
```

## Using `fanctl`
- Build with `cd fanctl && go build`.
- For every tagged version of the Go Module, the `fanctl` utility is built (Available on the Releases page).

```shell
$ fanctl set 50
$ fanctl get
50
$ fanctl rpm
2623
```

## Comments

### Direct Setting mode vs. Fan Speed Control mode
The EMC2301 has a built-in closed loop Fan Speed Control algorithm (Besides many awesome features), that allows the user to select a target RPM which is controlled by a PID.
I opted for the Direct Setting mode, because it feels more natural to me (Yes, i know that duty cycle/RPM is far from linear...).

### Should i use this for my 24/7-running project?
If you want a reliable/automatic cooling solution, you should probably go with something like https://github.com/neg2led/cm4io-fan.

There are scenarios where this Module might be handy (Well, that's why i built it...). But keep in mind, all you do with it, is setting the fan's PWM signal:
- Some fans won't stop with 0% duty cycle.
- Most fans will stall with a low duty cycle.
- ...

### TACH to RPM conversion
Unfortunately it is not possible to detect RPM's lower than 500 using the EMC2301. That's why the RPM measurement returns a `RPMResult`:

```go
type RPMResult struct {
    Rpm int
    Stopped bool
    Undef bool
}
```

- When the PWM duty cycle is set to 0%, `Stopped` is `true`. This is only true for fans that are able to stop completely!
- When the calculated RPM is below 500 (smallest value could be 480), it is not possible to determine the real RPM. In this case `Undef` is `true`.

### On poles, edges and ranges...
The EMC2301 may be used to control a wide variety of PWM fans. TACH measurement works different based on the fan's design type (emc2301.pdf: 4.4 Tachometer Measurement).

Noctua's PWM specification states that (regardless of the actual motor design) the tachometer always pulses two times per revolution (Emulating a fan with 2 poles). I **think** this is true for other regular, modern PWM fans. That's why i settled with the following values for the RPM calculation:
- *poles = 2*
- *n = 5*

Since i wanted to measure RPM's as low as possible, *RANGE[1:0]* is set to 500 RPM. This means the TACH multiplier is *m = 1*.

Using these values, a simplified formula for TACH conversion can be used (emc2301.pdf: EQUATION 4-3).
