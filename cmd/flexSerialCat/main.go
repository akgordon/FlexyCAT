package main

import (
	"FlexyCAT/internal/agg_serial_com"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	args := os.Args
	if len(args) < 2 {
		fmt.Print("Welcome to Flexy CAT - Serial CAT \n")
		fmt.Print("    by the Alan Gordon Group\n")
		fmt.Print("           v1.0\n")
		fmt.Print("\n")
		fmt.Print("Edit the config.ini file for port and baud rate\n")
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("Commands:\n")
		fmt.Print("\n")
		fmt.Print("  CAT;ZZxx  send CAT command to radio. Can be string together using semi-colon separator\n")
		fmt.Print("  Example:\n")
		fmt.Print("     CAT;ZZMD03;ZZFA00014054350;\n")
		fmt.Print("\n")
		fmt.Print("  GET   to get response from radio and save for later send back to radio\n")
		fmt.Print("  SET   send stored data to radio\n")
		fmt.Print("\n")
		fmt.Print("     Responses are saved by id number later to send command\n")
		fmt.Print("     xx = id number used for later SET\n")
		fmt.Print("     cmd = CAT command. e.g ZZFA;  Can string commands and save under one ID. \n")
		fmt.Print("  GET:xx:cmd1:cmd2:....\n")
		fmt.Print("\n")
		fmt.Print("Example to get current VFO-A settings\n")
		fmt.Print("   GET;01;ZZFA;ZZFI;ZZGT;ZZMD;ZZMG;ZZRG;ZZRT;ZZXG;\n")
		fmt.Print("\n")
		fmt.Print("Example to send saved data to radio\n")
		fmt.Print("   SET;01\n")
		fmt.Print("\n")
		os.Exit(0)
	}

	// Get port info
	commport, baudrate := getConfig("config.ini")
	fmt.Printf("Using port:%s  with baud rate:%d\n", commport, baudrate)

	// Get any saved data
	savedData, err := getSavedData("saved.txt")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Run command
	if strings.Index(args[1], "CAT") == 0 {
		doCAT(args[1], commport, baudrate)
	}

	err = setSavedData("saved.txt", savedData)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
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

func getSavedData(filename string) (map[int]string, error) {

	var storedCmds map[int]string

	// See if file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return storedCmds, nil
	}

	// Open file
	var filePtr *os.File
	filePtr, err = os.Open(filename)
	if err != nil {
		return storedCmds, err
	}
	defer filePtr.Close()

	// Read in data
	// Format:  xxx:<data>

	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		ln := scanner.Text()

		ids := ln[:3]
		id, err := strconv.Atoi(ids)
		if err != nil {
			continue
		}

		trimedData := strings.TrimSpace(ln[4:])
		storedCmds[id] = trimedData
	}

	return storedCmds, nil
}

func setSavedData(filename string, data map[int]string) error {

	// Open file for writing
	var filePtr *os.File
	var err error
	filePtr, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	defer filePtr.Close()

	for key, val := range data {
		ids := fmt.Sprintf("%03d", key)
		filePtr.WriteString(ids + ":" + val + "\n")
	}

	return nil
}

func split(r rune) bool {
	return r == ';'
}

func doCAT(cmd string, commport string, baudrate int) {
	flds := strings.FieldsFunc(cmd, split)
	isFirst := true
	for _, v := range flds {
		if isFirst {
			isFirst = false
			continue
		} else {
			v = v + ";"
			fmt.Printf("%s\n", v)
			agg_serial_com.SerialComm_Write(commport, baudrate, []byte(v))
		}
	}
}
