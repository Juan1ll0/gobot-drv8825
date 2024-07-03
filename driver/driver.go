package driver

import (
	"fmt"
	"math"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// const RPS = 0.0167

// Path: drv8825_driver.go

type DRV8825Direction uint8

const (
	Forward DRV8825Direction = iota
	Backward
)

type DRV8825Driver struct {
	adaptor            gpio.DigitalWriter
	name               string
	stepPin            string
	dirPin             string
	enablePin          string
	angle              float64
	isEnabled          bool
	signalDelay        time.Duration
	direction          DRV8825Direction
	stepsPerRevolution int
}

type Option func(*DRV8825Driver)

func WithSignalDelay(delay int64) Option {
	return func(d *DRV8825Driver) {
		d.signalDelay = time.Duration(delay) * time.Microsecond
	}
}

func WithAngle(angle float64) Option {
	return func(d *DRV8825Driver) {
		d.angle = angle
	}
}

func WithEnablePin(pin string) Option {
	return func(d *DRV8825Driver) {
		d.enablePin = pin
	}
}

func NewDRV8825Driver(adaptor gpio.DigitalWriter, stepPin string, dirPin string, options ...Option) *DRV8825Driver {
	// Initialize with default values
	drv8825 := &DRV8825Driver{
		adaptor:            adaptor,
		stepPin:            stepPin,
		dirPin:             dirPin,
		enablePin:          "",
		angle:              1.8,
		isEnabled:          true,
		signalDelay:        2 * time.Microsecond,
		stepsPerRevolution: int(360 / 1.8),
	}

	for _, option := range options {
		option(drv8825)
	}
	fmt.Println("Enable pin", drv8825)
	// Enable driver by default
	// drv8825.SetEnable(true)

	return drv8825
}

func (d *DRV8825Driver) StepPin() string {
	return d.stepPin
}

func (d *DRV8825Driver) EnablePin() string {
	return d.enablePin
}

func (d *DRV8825Driver) SignalDelay() time.Duration {
	return d.signalDelay
}

func (d *DRV8825Driver) DirPin() string {
	return d.dirPin
}

func (d *DRV8825Driver) Angle() float64 {
	return d.angle
}

func (d *DRV8825Driver) Name() string {
	return d.name
}

func (d *DRV8825Driver) SetName(name string) {
	d.name = name
}

func (d *DRV8825Driver) Direction() DRV8825Direction {
	return d.direction
}

func (d *DRV8825Driver) SetDirection(direction DRV8825Direction) {
	d.direction = direction
	d.adaptor.DigitalWrite(d.dirPin, byte(d.direction))

}

func (d *DRV8825Driver) SetEnable(enable bool) {
	//if d.enablePin != "" {
	//	if enable {
	//		if err := d.adaptor.DigitalWrite(d.enablePin, 0); err != nil {
	//			fmt.Println("failed to enable driver")
	//			return
	//		}
	//	} else {
	//		if err := d.adaptor.DigitalWrite(d.enablePin, 1); err != nil {
	//			fmt.Println("failed to disable driver")
	//			return
	//		}
	//	}
	//}
	d.isEnabled = enable
}

func (d *DRV8825Driver) SetSpeed(rpm int) {
	microseconds := int64((30000000 / (rpm * d.stepsPerRevolution)))
	fmt.Println(microseconds)
}

func calculateSlope(x1, y1, x2, y2 int) (float64, error) {
	x := x2 - x1
	if x == 0 {
		return 0.0, fmt.Errorf("division by zero")
	}

	value := float64(y2-y1) / float64(x)
	return value, nil
}

func calculateSpeed(step, rpm, acceleration int, slope float64) float64 {
	speed := float64(rpm) - math.Abs(slope)*(float64(acceleration)-float64(step))
	if speed > float64(rpm) {
		return float64(rpm)
	}
	return speed
}

func (d *DRV8825Driver) checkCanMove() bool {
	if !d.isEnabled {
		return false
	}
	return true
}

func (d *DRV8825Driver) Move(steps int, rpm int, accelSteps int, decelSteps int) error {
	if d.checkCanMove() {
		var speed float64
		accelSlope, _ := calculateSlope(0, 0, accelSteps, rpm)
		decelSlope, _ := calculateSlope(steps-decelSteps, rpm, steps, 0)

		fmt.Print("Delay", d.signalDelay)
		for i := 1; i < steps; i++ {
			if err := d.adaptor.DigitalWrite(d.stepPin, 1); err != nil {
				fmt.Print("failed on")
			}
			// Signal needs to be high for +-2 microseconds in DRV8825
			time.Sleep(d.signalDelay)
			if err := d.adaptor.DigitalWrite(d.stepPin, 0); err != nil {
				fmt.Print("failed off")
			}
			switch {
			case i <= accelSteps:
				speed = calculateSpeed(i, rpm, accelSteps, accelSlope)
			case i >= steps-decelSteps:
				speed = calculateSpeed(steps-i, rpm, decelSteps, decelSlope)
			default:
				speed = float64(rpm)
			}

			// Speed delay
			microseconds := int64(30000000 / (speed * float64(d.stepsPerRevolution)))
			time.Sleep(time.Duration(microseconds) * time.Microsecond)
		}
		return nil
	}
	return fmt.Errorf("driver is disabled")
}

func (d *DRV8825Driver) Connection() gobot.Connection {
	return d.adaptor.(gobot.Connection)
}

func (d *DRV8825Driver) Finalize() error {
	return nil
}

func (d *DRV8825Driver) Device() gobot.Device {
	return nil
}

func (d *DRV8825Driver) Start() error {
	return nil
}

// Halt implements the Driver interface
func (d *DRV8825Driver) Halt() error {
	return nil
}
