package cm4iofan

import (
	"errors"
	"math"

	"github.com/go-daq/smbus"
)

const (
	/* CM4IO */

	// addresses (cm4io.pdf: 2.9 Fan connector)
	Emc2301Bus  = 10
	Emc2301Addr = 0x2F

	/* EMC2301 */

	// general (emc2301.pdf: TABLE 6-1)
	Emc2301ProductIdReg = 0xFD
	Emc2301ProductIdVal = 0x37

	// configuration (emc2301.pdf: 4.1 Fan Control Modes of Operation, 6.14 Fan Configuration Registers)
	Emc2301ConfigReg = 0x32

	// fan drive setting (emc2301.pdf: 6.12 Fan Drive Setting Register)
	Emc2301DutyCycleReg = 0x30

	// TACH reading (emc2301.pdf: 6.23 TACH Reading Registers)
	Emc2301TachHighReg = 0x3E
	Emc2301TachLowReg  = 0x3F

	// TACH conversion (emc2301.pdf: EQUATION 4-3: SIMPLIFIED TACH CONVERSION)
	Emc2301Tach2RPM = 3932160
)

type EMC2301 struct {
	conn *smbus.Conn
}

func New() (*EMC2301, error) {
	c, err := smbus.Open(Emc2301Bus, Emc2301Addr)
	if err != nil {
		return nil, err
	}
	id, err := c.ReadReg(Emc2301Addr, Emc2301ProductIdReg)
	if err != nil {
		return nil, err
	}
	if Emc2301ProductIdVal != id {
		return nil, errors.New("unexpected Product ID")
	}
	conf, err := c.ReadReg(Emc2301Addr, Emc2301ConfigReg)
	if err != nil {
		return nil, err
	}
	// set RNG1[1:0] to 500 RPM
	conf &= ^(uint8(0b11) << 5)
	err = c.WriteReg(Emc2301Addr, Emc2301ConfigReg, conf)
	if err != nil {
		return nil, err
	}
	return &EMC2301{conn: c}, nil
}

func (ctrl *EMC2301) GetDutyCycle() (int, error) {
	v, err := ctrl.conn.ReadReg(Emc2301Addr, Emc2301DutyCycleReg)
	if err != nil {
		return -1, err
	}
	// emc2301.pdf: EQUATION 4-1: REGISTER VALUE TO DRIVE
	return int(math.Round((float64(v) / 255) * 100)), nil
}

func (ctrl *EMC2301) SetDutyCycle(percent int) error {
	if percent < 0 || percent > 100 {
		return errors.New("expecting a value of 0 <= x <= 100")
	}
	// emc2301.pdf: EQUATION 4-1: REGISTER VALUE TO DRIVE
	v := math.Round(255 * (float64(percent) / 100))
	return ctrl.conn.WriteReg(Emc2301Addr, Emc2301DutyCycleReg, uint8(v))
}

func (ctrl *EMC2301) GetRPM() (int, error) {
	h, err := ctrl.conn.ReadReg(Emc2301Addr, Emc2301TachHighReg)
	if err != nil {
		return -1, err
	}
	l, err := ctrl.conn.ReadReg(Emc2301Addr, Emc2301TachLowReg)
	if err != nil {
		return -1, err
	}
	// HIGH BYTE - Bit 7: 2048 ... Bit 0: 32
	// LOW BYTE - Bit 7: 16 ... Bit 3: 1, Bit 0-2: ignored
	tach := uint16(h) << 5
	tach |= uint16(l) >> 3
	// emc2301.pdf: EQUATION 4-3: SIMPLIFIED TACH CONVERSION
	return int(Emc2301Tach2RPM / uint32(tach)), nil
}
