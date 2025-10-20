# PMS7003 Go Driver

A Go package for interfacing with the PMS7003 particulate matter sensor. The PMS7003 is a digital particle concentration sensor that measures suspended particles in the air and outputs data via a serial interface.

This package provides a simple, typed interface to connect to the sensor and read air quality measurements. It works with any device that supports serial communication, such as Raspberry Pi, ESP32, or standard computers with USB-to-serial adapters.

## Features

- 📊 Read PM1.0, PM2.5, and PM10 concentrations (both CF=1 standard and atmospheric environment)
- 🔢 Access particle counts for different size thresholds (>0.3μm, >0.5μm, >1.0μm, >2.5μm, >5.0μm, >10μm)
- 🔄 Support for both Active and Passive reading modes
- 💤 Sleep/Wake functionality to control sensor power consumption
- ✅ Built-in checksum validation for data integrity
- 🔌 Simple API with minimal dependencies

## Requirements

- Go 1.19 or higher
- Serial port access to the PMS7003 sensor (e.g., `/dev/ttyAMA0` on Raspberry Pi, `/dev/ttyUSB0` on Linux with USB adapter)
- Appropriate permissions to access the serial device

## Installation

```bash
go get github.com/shivasaxena/PMS7003
```

## Quick Start

Here's a simple example to get started:

```go
package main

import (
    "fmt"
    "time"
    
    PMS7003 "github.com/shivasaxena/PMS7003"
)

func main() {
    // Open connection to the sensor in Active mode
    device, err := PMS7003.Open("/dev/ttyAMA0", PMS7003.ActiveMode)
    if err != nil {
        panic(err)
    }
    defer device.Close()

    // Wait for sensor to stabilize
    time.Sleep(10 * time.Second)

    // Read sensor data
    value, err := device.Read()
    if err != nil {
        panic(err)
    }

    fmt.Printf("PM2.5: %d μg/m³\n", value.PM25Atmospheric)
    fmt.Printf("PM10: %d μg/m³\n", value.PM10Atmospheric)
}
```

For a complete example, see [cmd/main.go](cmd/main.go).

## Usage

### Opening a Connection

The [`Open`](PMS7003.go) function establishes a serial connection to the sensor:

```go
device, err := PMS7003.Open("/dev/ttyAMA0", PMS7003.ActiveMode)
```

**Modes:**
- [`ActiveMode`](PMS7003.go): Sensor continuously transmits data (default behavior)
- [`PassiveMode`](PMS7003.go): Sensor only sends data when requested

### Reading Sensor Data

Use the [`Read`](PMS7003.go) method to get current measurements:

```go
value, err := device.Read()
if err != nil {
    // handle error
}

// Access measurements
fmt.Println(value.PM25Atmospheric)    // PM2.5 concentration
fmt.Println(value.PM10Atmospheric)    // PM10 concentration
fmt.Println(value.ParticlesGT25)      // Particle count > 2.5μm
```

### Power Management

Control the sensor's power state to reduce energy consumption:

```go
// Put sensor to sleep
err := device.Sleep()

// Wake up sensor
err := device.WakeUp()
```

**Note:** After waking up, wait ~30 seconds for the fan to stabilize before taking readings.

### Closing the Connection

Always close the connection when done:

```go
defer device.Close()
```

## API Reference

### Types

The types are modeled from the references of device,packet and communication protocol as specided in the [manual](./docs/static_files/PMS7003_series_data_manual_English_V2.5.pdf).

- [`PMS7003Device`](PMS7003.go) - Main device handle with methods to interact with the sensor
- [`PMS7003SensorValue`](PMS7003Model.go) - Structure containing all sensor measurements
- [`Mode`](PMS7003.go) - Enum for Active/Passive modes

### Key Functions

- [`Open(serialDevice string, mode Mode) (PMS7003Device, error)`](PMS7003.go) - Opens serial connection and initializes the sensor
- [`Read() (PMS7003SensorValue, error)`](PMS7003.go) - Reads and parses sensor data
- [`WakeUp() error`](PMS7003.go) - Wakes sensor from sleep mode
- [`Sleep() error`](PMS7003.go) - Puts sensor into sleep mode
- [`Close()`](PMS7003.go) - Closes the serial connection

### Sensor Data Fields

[`PMS7003SensorValue`](PMS7003Model.go) contains:

| Field | Description |
|-------|-------------|
| `PM10CF10Standard`, `PM25CF10Standard`, `PM100CF10Standard` | Standard particle concentrations (CF=1) |
| `PM10Atmospheric`, `PM25Atmospheric`, `PM100Atmospheric` | Atmospheric environment concentrations |
| `ParticlesGT03`, `ParticlesGT05`, `ParticlesGT10` | Particle counts >0.3μm, >0.5μm, >1.0μm |
| `ParticlesGT25`, `ParticlesGT50`, `ParticlesGT100` | Particle counts >2.5μm, >5.0μm, >10μm |

## Hardware Setup

### Raspberry Pi

1. Enable serial port in `raspi-config`
2. Connect sensor to GPIO pins (TX, RX, 5V, GND)
3. Use device path `/dev/ttyAMA0` or `/dev/ttyS0`
4. Ensure user has permission: `sudo usermod -a -G dialout $USER`

### Linux with USB Adapter

1. Connect PMS7003 via USB-to-serial adapter
2. Identify device: `ls /dev/ttyUSB*`
3. Set permissions: `sudo chmod 666 /dev/ttyUSB0`

## Troubleshooting

**Permission Denied:**
```bash
sudo usermod -a -G dialout $USER
# Log out and back in
```

**No Data Received:**
- Verify correct serial device path
- Check physical connections
- Ensure sensor is powered (5V)
- Wait 30+ seconds after `WakeUp()` call

**Checksum Errors:**
- Check for loose connections
- Verify power supply stability
- Try a different USB cable/adapter

## Air Quality Index (AQI)

Use the PM2.5 readings to calculate AQI according to your regional standards (EPA, WHO, etc.). The atmospheric environment values (`PM25Atmospheric`) are typically used for AQI calculations.

## Building

```bash
# Build the example
go build ./cmd

# Run directly
go run cmd/main.go

# Run tests (if available)
go test ./...
```

## Dependencies

- [go.bug.st/serial](https://github.com/bugst/go-serial) - Cross-platform serial port library

See [go.mod](go.mod) for complete dependency list.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## References

- [PMS7003 Datasheet](http://download.kamami.pl/p564008-PMS7003%20series%20data%20manua_English_V2.5.pdf)
- [GitHub Repository](https://github.com/shivasaxena/PMS7003)

## Acknowledgments

Thanks to all contributors and the open-source community for making this project possible.