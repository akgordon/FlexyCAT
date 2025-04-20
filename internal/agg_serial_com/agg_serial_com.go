package agg_serial_com

import (
	"fmt"
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

func SerialDataSendAndReceive(portName string, baudRate int, dataToSend string) (string, error) {
	// Configure the serial port
	config := &serial.Config{Name: portName, Baud: baudRate, ReadTimeout: time.Second * 5}
	port, err := serial.OpenPort(config)
	if err != nil {
		return "", fmt.Errorf("failed to open port: %v", err)
	}
	defer port.Close()

	// Send data
	_, err = port.Write([]byte(dataToSend))
	if err != nil {
		return "", fmt.Errorf("failed to send data: %v", err)
	}

	time.Sleep(time.Millisecond * 250)

	// Wait for and read the response
	buffer := make([]byte, 1024) // Adjust buffer size as needed
	n, err := port.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Return the received response as a string
	return string(buffer[:n]), nil
}
