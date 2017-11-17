package bitstream

import (
	"errors"
)

// Bitstream : core object
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
func (stream *Bitstream) ReadBits(numBits uint) ([]uint8, error) {
	bits := make([]uint8, 0)
	cpt := uint(0)
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

// ReadBitsAsInt : Read a defined number of bits and return as an unsigned int
func (stream *Bitstream) ReadBitsAsInt(numBits uint) (uint, error) {
	bits, err := stream.ReadBits(numBits)
	if err != nil {
		return 0x00, err
	}

	// Casting bits to an integer
	num := uint(0)
	for bitIdx := len(bits); bitIdx > 0; bitIdx-- {
		bit := bits[bitIdx-1]
		num += uint(bit << uint(len(bits)-(bitIdx)))
	}

	// sending data back
	return num, nil
}

// ReadBytes: Read a defined number of bytes and returns the value as uint array
func (stream *Bitstream) ReadBytes(numBytes uint) ([]uint, error) {
	bytes := make([]uint, 0)
	cpt := uint(0)
	for cpt < numBytes {
		bits, err := stream.ReadBitsAsInt(8)
		if err != nil {
			return bytes, err
		}
		bytes = append(bytes, bits)
		cpt++
	}
	return bytes, nil
}

// ReadGolomb : Returns a int value (signed or not) from bitstream using Exp-Golomb code
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

// Remains function send the number of present bits in the stream
func (stream *Bitstream) Remains() int {
	return len(stream.data) - int(stream.offset)
}
