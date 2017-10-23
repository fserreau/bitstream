package main

import (
	"errors"
	"fmt"
)

type Bitstream struct {
	data   []byte
	offset uint8
}

// InitStream : Initialisation of a bit stream with a byte array
func InitStream(byteStream []byte) *Bitstream {
	stream := new(Bitstream)
	stream.data = byteStream
	stream.offset = 0
	return stream
}

// ReadBit : Read a single bit from bitstream
func (stream *Bitstream) ReadBit() (uint8, error) {
	if len(stream.data) == 0 {
		return 2, errors.New("Stream is empty, can not read bit")
	}
	bit := stream.data[stream.offset/8] >> uint(7-(stream.offset%8)) & 0x01
	stream.offset++
	if stream.offset == 8 {
		stream.data = stream.data[1:]
		stream.offset = 0
	}
	return bit, nil
}

// ReadBits : Read a defined number of bits and return an array
func (stream *Bitstream) ReadBits(numBits int) ([]uint8, error) {
	bits := make([]uint8, 0)
	cpt := 0
	for cpt < numBits {
		bit, err := stream.ReadBit()
		if err != nil {
			return bits, err
		}
		bits = append(bits, bit)
		cpt++
	}
	return bits, nil
}

// ReadByte : get the next 8bits as a byte
func (stream *Bitstream) ReadByte() (byte, error) {
	bits, err := stream.ReadBits(8)
	if err != nil {
		return 0x00, err
	}
	return bits[7] + (2 * bits[6]) + (4 * bits[5]) + (8*bits[4] + (16 * bits[3]) + (32 * bits[2]) + (64 * bits[1]) + (128 * bits[0])), nil
}

func (stream *Bitstream) ReadGolomb(signed bool) (int, error) {
	data := make([]uint8, 0)
	for {
		bit, err := stream.ReadBit()
		if err != nil {
			return 0, err
		}
		data = append(data, bit)
		if bit == 1 {
			break
		}
	}
	fmt.Println(data)
	tlen := len(data)

	// Concatenation of the end of bit array
	i := 1
	for i < tlen {
		bit, err := stream.ReadBit()
		if err != nil {
			return 0, err
		}
		data = append(data, bit)
		i++
	}
	fmt.Println(data)

	// Casting bits into an integer
	num := 0
	for bitIdx := len(data); bitIdx > 0; bitIdx-- {
		bit := data[bitIdx-1]
		num += int(bit << uint(len(data)-(bitIdx)))
	}
	num--

	// Getting expected value
	if signed {
		if num%2 == 0 {
			num = -1 * (num / 2)
		} else {
			num = (num + 1) / 2
		}
	}
	return num, nil
}

func main() {
	bitStream := InitStream([]byte{0x3F, 0x9D, 0x10})
	for {
		num, err := bitStream.ReadGolomb(false)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(num)
	}
}
