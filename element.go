package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
)

type element struct {
	t    elementType
	name string
	data interface{}
}

type elementType byte

const (
	ET_double elementType = 0x01
	ET_string elementType = 0x02
	ET_array  elementType = 0x04
	ET_int32  elementType = 0x10
)

func readElementList(docb *bufio.Reader, bytesLeft int32) []element {
	elements := []element{}
	for bytesLeft > 0 {
		el, bytesRead := readElement(docb)
		bytesLeft -= bytesRead
		elements = append(elements, el)
	}
	return elements
}

func readElement(docb *bufio.Reader) (element, int32) {
	el := element{}
	var bytesRead, temp int32

	etb, err := docb.ReadByte()
	el.t = elementType(etb)
	bytesRead++
	if err != nil {
		panic(err)
	}
	el.name, temp = readCString(docb)
	bytesRead += temp
	switch el.t {
	case ET_double:
		data := make([]byte, 8)
		_, err := docb.Read(data)
		if err != nil {
			panic(err)
		}
		el.data = math.Float64frombits(binary.LittleEndian.Uint64(data))
		bytesRead += 8
	case ET_string:
		el.data, temp = readString(docb)
		bytesRead += temp
	case ET_array:
		el.data, temp = readArray(docb)
		bytesRead += temp
	case ET_int32:
		el.data = readInt32(docb)
		bytesRead += 4
	default:
		panic(fmt.Sprintf("Element type not implemented yet: %d", el.t))
	}
	return el, bytesRead
}
