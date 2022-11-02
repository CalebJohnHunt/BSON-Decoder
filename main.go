package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
)

var ctob map[byte]byte = map[byte]byte{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'a': 0xa,
	'b': 0xb,
	'c': 0xc,
	'd': 0xd,
	'e': 0xe,
	'f': 0xf,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Args: write ... | read | test")
	}
	switch os.Args[1] {
	case "write":
		bs := os.Args[2]
		arr := []byte{}
		for i := len(bs) - 1; i >= 0; {
			b := ctob[bs[i]]
			i--
			if i >= 0 {
				b |= (ctob[bs[i]] << 4)
				i--
			}
			arr = append([]byte{b}, arr...)
		}
		os.WriteFile("document.bson", arr, fs.FileMode(0777))
	case "read":
		f, _ := os.Open("document.bson")
		defer f.Close()
		docb := bufio.NewReader(f)
		fmt.Println(readDoc(docb))
	case "test":
		f, err := os.Open(("document.bson"))
		if err != nil {
			panic(err)
		}
		docb := bufio.NewReader(f)
		fmt.Println(docb.ReadString(0x00))
	default:
		fmt.Println("No recognized argument")
	}
}

type document struct {
	size     int32
	elements []element
}

func readDoc(docb *bufio.Reader) document {
	doc := document{}

	// size
	doc.size = readInt32(docb)
	bytesLeft := doc.size - 4
	doc.elements = readElementList(docb, bytesLeft-1)
	eof, err := docb.ReadByte()
	if err != nil {
		panic(err)
	}
	if eof != 0x00 {
		panic("Did not end with eof")
	}

	return doc
}

func readInt32(rd *bufio.Reader) int32 {
	b := make([]byte, 4)
	_, err := rd.Read(b)
	if err != nil {
		panic(err)
	}
	return int32(binary.LittleEndian.Uint32(b))
}

func readArray(docb *bufio.Reader) (document, int32) {
	doc := readDoc(docb)
	return doc, doc.size
}

func readString(docb *bufio.Reader) (string, int32) {
	readInt32(docb)
	a, b := readCString(docb)
	return a, b + 4 // + 4 for the int32
}

func readCString(docb *bufio.Reader) (string, int32) {
	bs, err := docb.Peek(1)
	if err != nil {
		panic(err)
	}
	// Empty string
	if bs[0] == 0x00 {
		return "", 1
	}
	s, err := docb.ReadString(0x00)
	if err != nil {
		panic(err)
	}
	return s[:len(s)-1], int32(len(s))
}
