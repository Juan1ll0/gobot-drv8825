package main

import (
	"drv8825/driver"
	"gobot.io/x/gobot/v2/platforms/firmata"

	"gobot.io/x/gobot/v2"
)

func main() {
	r := firmata.NewAdaptor("/dev/ttyUSB0")
	// bus := 1
	// mcpDevice := i2c.NewMCP23017Driver(r, i2c.WithBus(bus), i2c.WithAddress(0x21))
	// pcaDevice := i2c.NewPCA9685Driver(r, i2c.WithBus(bus), i2c.WithAddress(0x40))

	// mcpDigitalWriter := mcpDigitalWriter.NewMCPDigitalWriter(mcpDevice)
	drv8825 := driver.NewDRV8825Driver(r, "24", "25",
		driver.WithAngle(1.8),
		driver.WithSignalDelay(0),
		//driver.WithEnablePin("40"),
	)
	// drv8825.SetEnable(true)

	work := func() {
		// Set direction
		drv8825.SetDirection(driver.Forward)
		drv8825.Move(10000, 500, 2000, 3000)
		drv8825.SetDirection(driver.Backward)
		drv8825.Move(10000, 20000, 0, 0)

	}

	robot := gobot.NewRobot("stepperBot",
		[]gobot.Connection{r},
		// []gobot.Device{mcpDevice, drv8825},
		[]gobot.Device{drv8825},
		work,
	)

	robot.Start()
}
