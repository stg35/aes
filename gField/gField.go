package gField

import "errors"

func MultiplyTwoBytes(a byte, b byte) byte{
	pos := findOnesPositions(b)
	sum := byte(0)
	var ax byte
	for i := range pos {
		if pos[i] == 0 {
			sum ^= a
			continue
		}
		ax = a
		for j := 0; j < pos[i]; j++ {
			ax = MultiplyX(ax)
		}
		sum ^= ax
	}
	return sum
}

func FindInverseElement(b byte) (byte, error) {
	if b == 0 {
		return 0, nil
	}

	for i := 1; i < 256; i++ {
		if MultiplyTwoBytes(byte(i), b) == 1 {
			return byte(i), nil
		}
	}
	return 0, errors.New("can't find inverse element")
}

func MultiplyX(b byte) byte{
	if highBit := b & 0x80; highBit == 0 {
		return b << 1
	} else {
		return (b << 1)^0x1b
	}
}

func findOnesPositions(b byte) []int{
	counter := 0
	res := []int{}
	for b > 0 {
		if b % 2 != 0 {
			res = append(res, counter)
		}
		b /= 2
		counter++
	}
	return res
}

