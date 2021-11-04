package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	hexBinData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Println("open hex file err:", err)
		os.Exit(1)
	}

	encoding, encodingName, certain := charset.DetermineEncoding(hexBinData, "")

	var hexStr string

	if !certain {
		fmt.Println("can not detect hex file encoding, treat it as utf-8 file")
		hexStr = string(hexBinData)
	} else {
		decoder := encoding.NewDecoder()
		utf8bytes, err := decoder.Bytes(hexBinData)
		if err != nil {
			fmt.Println("process hex file err, read as ", encodingName, err)
			os.Exit(2)
		}
		hexStr = string(utf8bytes)
	}

	replacer := strings.NewReplacer(
		" ", "",
		"\n", "",
		"\r", "",
		"0x", "",
		"0X", "",
		"h", "",
		"H", "",
		":", "",
		",", "",
		"/", "",
		"[", "",
		"]", "",
		"'", "",
		"\"", "",
		".", "")

	hexData := replacer.Replace(hexStr)

	binData, err := hex.DecodeString(hexData)
	if err != nil {
		fmt.Println("hex decode err:", err, hexData)
		os.Exit(3)
	}

	binFilename := os.Args[1] + ".bin"

	err = ioutil.WriteFile(binFilename, binData, 0666)
	if err != nil {
		fmt.Println("save bin file err:", err)
		os.Exit(4)
	}

}

func usage() {
	fmt.Println(
		`yhex2bin hex_filename [bin_filename]`)
}
