package cm4iofan

import (
	"errors"
	"github.com/go-daq/smbus"
)

const (
	Emc2301Bus          = 10   // cm4io.pdf: 2.9 Fan connector
	Emc2301Addr         = 0x2F // cm4io.pdf: 2.9 Fan connector
	Emc2301RegProductId = 0xFD // emc2301.pdf: TABLE 6-1
	Emc2301ProductId    = 0x37 // emc2301.pdf: TABLE 6-1
	Emc2301RegConfig    = 0x32 // emc2301.pdf: 4.1 Fan Control Modes of Operation, 6.14 Fan Configuration Registers
	Emc2301RegDutyCycle = 0x30 // emc2301.pdf: 4.1 Fan Control Modes of Operation, 6.12 Fan Drive Setting Register
)

type EMC2301 struct {
	conn *smbus.Conn
}

func New() (*EMC2301, error) {
	c, err := smbus.Open(Emc2301Bus, Emc2301Addr)
	if err != nil {
		return nil, err
	}
	v, err := c.ReadReg(Emc2301Addr, Emc2301RegProductId)
	if err != nil {
		return nil, err
	}
	if Emc2301ProductId != v {
		return nil, errors.New("unexpected ProductId")
	}
	// ensure FSC (Fan Speed Control) mode is DISABLED
	v, err = c.ReadReg(Emc2301Addr, Emc2301RegConfig)
	if err != nil {
		return nil, err
	}
	err = c.WriteReg(Emc2301Addr, Emc2301RegConfig, v|0<<7)
	if err != nil {
		return nil, err
	}
	return &EMC2301{conn: c}, nil
}

func (ctrl *EMC2301) SetDutyCycle(dc uint8) error {
	if dc < 0 || dc > 100 {
		return errors.New("expecting a value 0 <= x <= 100")
	}
	// emc2301.pdf: EQUATION 4-1
	v := 255 * (float32(dc) / 100)
	return ctrl.conn.WriteReg(Emc2301Addr, Emc2301RegDutyCycle, uint8(v))
}
