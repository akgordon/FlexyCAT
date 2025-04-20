package agg_serial_com

import (
	"log"
	"time"

	"github.com/tarm/serial"
)

var SerialComm serial.Config

func SerialComm_Read(commport string, baud int) (data []byte, len int, err error) {

	// Configure the serial port
	config := &serial.Config{
		Name: commport, //e.g. COM5
		Baud: baud,     // Set your baud rate
		// Optional: add ReadTimeout if you want a timeout
		ReadTimeout: time.Second * 2,
	}

	// Open the port
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Read data from the port
	buf := make([]byte, 128)
	n, err := port.Read(buf)

	return buf, n, err
}

func SerialComm_Write(commport string, baud int, data []byte) (len int, err error) {

	// Configure the serial port
	config := &serial.Config{
		Name: commport, //e.g. COM5
		Baud: baud,     // Set your baud rate
		// Optional: add ReadTimeout if you want a timeout
		ReadTimeout: time.Second * 2,
	}

	// Open the port
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Write data to the port
	var n int
	n, err = port.Write(data)

	return n, err
}
