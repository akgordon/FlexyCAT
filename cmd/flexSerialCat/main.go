package main

import (
	"FlexyCAT/internal/agg_serial_com"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type savedData struct {
	savedData map[string]string
}

func main() {

	args := os.Args
	if len(args) < 2 {
		fmt.Print("Welcome to Flexy CAT - Serial CAT \n")
		fmt.Print("    by the Alan Gordon Group\n")
		fmt.Print("           v1.0\n")
		fmt.Print("\n")
		fmt.Print("Edit the config.ini file for port and baud rate\n")
		fmt.Print("  or preface each command with 'CONFIG{PORT=COM8,BAUD=9600}'\n")
		fmt.Print("  Example:\n")
		fmt.Print("     CONFIG{PORT=COM8,BAUD=9600}:CAT:ZZMD03;ZZFA00014054350;\n")
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("Commands:\n")
		fmt.Print("\n")
		fmt.Print("  CAT:ZZxx  send CAT command to radio. Can be string together using semi-colon separator\n")
		fmt.Print("     Note use of colon and semi-colon in command string. \n")
		fmt.Print("  Example:\n")
		fmt.Print("     CAT:ZZMD03;ZZFA00014054350;\n")
		fmt.Print("\n")
		fmt.Print("  GET   to get response from radio and save for later send back to radio\n")
		fmt.Print("  SET   send stored data to radio\n")
		fmt.Print("\n")
		fmt.Print("     Responses are saved by id name for later to SET command\n")
		fmt.Print("     cmd = A FLEX CAT command. e.g ZZFA;  Can string commands and save under one ID. \n")
		fmt.Print("     Note use of colon and semi-colon in command string. \n")
		fmt.Print("\n")
		fmt.Print("Example to GET current VFO-A settings\n")
		fmt.Print("   GET:VFOA:ZZFA;ZZFI;ZZGT;ZZMD;ZZMG;ZZRG;ZZRT;ZZXG;\n")
		fmt.Print("\n")
		fmt.Print("Example to send saved data to radio\n")
		fmt.Print("   SET:VFOA\n")
		fmt.Print("\n")
		os.Exit(0)
	}

	flds := strings.FieldsFunc(args[1], splitColon)

	var theSavedData = &savedData{}
	theSavedData.savedData = make(map[string]string)

	//var savedData = make(map[string]string)

	// First check to see if port and baud in command line
	var commport = ""
	var baudrate = 0
	for _, v := range flds {
		if strings.Index(v, "CONFIG") == 0 {
			// CONFIG{PORT=COM8,BAUD=9600}
			v = v[7 : len(v)-1]
			cfgFlds := strings.FieldsFunc(v, splitComma)
			for _, c := range cfgFlds {
				if strings.Index(c, "PORT=") == 0 {
					commport = c[5:]
					commport = strings.TrimSpace(commport)
				}
				if strings.Index(c, "BAUD=") == 0 {
					bs := c[5:]
					bs = strings.TrimSpace(bs)
					baudrate, _ = strconv.Atoi(bs)
				}
			}
		}
	}
	if baudrate == 0 {
		// Get port info
		commport, baudrate = getConfig("config.ini")
	}
	if baudrate == 0 {
		fmt.Print("ERROR: Unable to get port and baud rate.")
		os.Exit(1)
	}
	fmt.Printf("Using port:%s  with baud rate:%d\n", commport, baudrate)

	//******************
	// Process command
	//******************
	var err error

	for idx, cmd := range flds {
		if (cmd == "GET") || (cmd == "SET") {
			// First get ANY saved info
			// Get any saved data
			err = theSavedData.getSavedData("saved.txt")
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
		}

		if cmd == "CAT" {
			// Run command
			if len(flds) > (idx + 1) {
				DoCAT(commport, baudrate, flds[idx+1])
			}
		}

		if cmd == "GET" {
			// Run command
			if len(flds) > (idx + 2) {
				theSavedData.DoGET(commport, baudrate, flds[idx+1], flds[idx+2])
			}
		}

		if cmd == "SET" {
			// Run command
			if len(flds) > (idx + 1) {
				theSavedData.DoSET(commport, baudrate, flds[idx+1])
			}
		}

		if (cmd == "GET") || (cmd == "SET") {
			// First get ANY saved info
			// Get any saved data
			err = theSavedData.setSavedData("saved.txt")
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
		}
	}
}

func getConfig(filename string) (commport string, baud int) {
	// Open file
	var filePtr *os.File
	var err error
	filePtr, err = os.Open(filename)
	if err != nil {
		return "", 0
	}
	defer filePtr.Close()

	// Read in structure
	var cport string = ""
	var baudrate int = 0

	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		ln := scanner.Text()
		if strings.Index(ln, "port:") == 0 {
			cport = ln[5:]
			cport = strings.TrimSpace(cport)
		}

		if strings.Index(ln, "baud:") == 0 {
			var err error
			bs := ln[5:]
			bs = strings.TrimSpace(bs)
			baudrate, err = strconv.Atoi(bs)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
	}

	return cport, baudrate
}

func (sd *savedData) getSavedData(filename string) error {

	// Clear any current map data
	for key := range sd.savedData {
		delete(sd.savedData, key)
	}

	// See if file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}

	// Open file
	var filePtr *os.File
	filePtr, err = os.Open(filename)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	// Read in data
	// Format:  xxx:<data>
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		ln := scanner.Text()
		flds := strings.FieldsFunc(ln, splitColon)
		if len(flds) == 2 {
			sd.savedData[strings.TrimSpace(flds[0])] = strings.TrimSpace(flds[1])
		}
	}

	return nil
}

func (sd *savedData) setSavedData(filename string) error {

	// Open file for writing
	var filePtr *os.File
	var err error
	filePtr, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	defer filePtr.Close()

	for key, val := range sd.savedData {
		filePtr.WriteString(key + ":" + val + "\n")
	}

	return nil
}

func splitColon(r rune) bool {
	return r == ':'
}
func splitSemiColon(r rune) bool {
	return r == ';'
}
func splitEqual(r rune) bool {
	return r == '='
}
func splitComma(r rune) bool {
	return r == ','
}

func DoCAT(commport string, baudrate int, cmd string) {
	flds := strings.FieldsFunc(cmd, splitSemiColon)
	for _, v := range flds {
		v = v + ";"
		fmt.Printf("%s\n", v)
		agg_serial_com.SerialComm_Write(commport, baudrate, []byte(v))
	}
}

func (sd *savedData) DoGET(commport string, baudrate int, key string, cmd string) {
	var responses string = ""
	flds := strings.FieldsFunc(cmd, splitSemiColon)
	for _, v := range flds {
		v = v + ";"
		resp, err := agg_serial_com.SerialDataSendAndReceive(commport, baudrate, v)

		if err == nil {
			responses += resp
		} else {
			fmt.Print(err.Error())
		}
	}
	sd.savedData[key] = responses
}

func (sd *savedData) DoSET(commport string, baudrate int, id string) {
	cmd, ok := sd.savedData[id]
	if ok {
		DoCAT(commport, baudrate, cmd)
	} else {
		fmt.Printf("No saved command named:%s", id)
	}
}
