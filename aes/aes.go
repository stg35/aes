package aes

import (
	"log"
	"math/bits"
	"github.com/stg35/aes/gField"
)

var (
	Nb = 4
	Nk = 4
	Nr = 10
	Rcon = [][]byte{
		[]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36},
		[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}
)


func SBox(b byte) byte{
	m := byte(0xf8)
	res := byte(0)
	inv, err := gField.FindInverseElement(b) 
	if err != nil {
		log.Fatal("err")
	}
	for i := 0; i < 8; i++ {
		res = (res << 1) | (byte(bits.OnesCount8(m & inv))%2)
		if lowBit := m & 0x01; lowBit == 1 {
			m = (m >> 1) | 0x80
		} else {
			m = (m >> 1)
		}
	}
	return res ^ 0x63
}

func KeyExpansion(key []byte) [][]byte {
	if len(key) < 4 * Nk {
		for i := 0; i < 4 * Nk - len(key); i++ {
			key = append(key, 0x01)
		}
	}

	keySchedule := [][]byte{[]byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}}
	for i := 0; i < 4; i ++ {
		for j := 0; j < Nk; j++ {
			keySchedule[i] = append(keySchedule[i], key[i+4*j])
		}
	}

	for col := Nk; col < Nb*(Nr+1); col++ {
		if col % Nk == 0 {
			a := []byte{}
			for row := 1; row < 4; row++ {
				a = append(a, keySchedule[row][col-1])
			}
			a = append(a, keySchedule[0][col-1])

			for i := range a {
				a[i] = SBox(a[i])
			}

			for row := 0; row < 4; row++ {
				keySchedule[row] = append(keySchedule[row], keySchedule[row][col-4]^a[row]^Rcon[row][col/Nk-1])
			}
		} else {
			for row := 0; row < 4; row++ {
				keySchedule[row] = append(keySchedule[row], keySchedule[row][col-4]^keySchedule[row][col-1])
			}
		}
	}
	return keySchedule
} 

func SubBytes(state [][]byte) [][]byte {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] = SBox(state[i][j])
		}
	}

	return state
}

func ShiftRows(state [][]byte) [][]byte {
	for counter, row := range state {
		for i := 0; i < counter; i++ {
			row = shiftSlice(row)
		}
	}
	return state
}

func shiftSlice(b []byte) []byte {
	x, s := b[0], b[1:]
	s = append(s, x)
	return s
}

func MixColumns(state [][]byte) [][]byte {
	for i := 0; i < Nb; i++ {
		state[0][i] = gField.MultiplyTwoBytes(state[0][i], 0x02) ^ gField.MultiplyTwoBytes(state[1][i], 0x03) ^ state[2][i] ^ state[3][i]
		state[1][i] = state[0][i] ^ gField.MultiplyTwoBytes(state[1][i], 0x02) ^ gField.MultiplyTwoBytes(state[2][i], 0x03) ^ state[3][i]
		state[2][i] = state[0][i] ^ state[1][i] ^ gField.MultiplyTwoBytes(state[2][i], 0x02) ^ gField.MultiplyTwoBytes(state[3][i], 0x03)
		state[3][i] = gField.MultiplyTwoBytes(state[0][i], 0x03) ^ state[1][i] ^ state[2][i] ^ gField.MultiplyTwoBytes(state[3][i], 0x02)
	} 
	return state
}

func AddRoundKey(state [][]byte, keySchedule [][]byte, round int) [][]byte {
	for col := 0; col < Nk; col++ {
		state[0][col] = state[0][col]^keySchedule[0][Nb*round + col]
        state[1][col] = state[1][col]^keySchedule[1][Nb*round + col]
        state[2][col] = state[2][col]^keySchedule[2][Nb*round + col]
        state[3][col] = state[3][col]^keySchedule[3][Nb*round + col]
	}
	return state
}

func Encryption(b []byte, key []byte) []byte {
	state := [][]byte{ []byte{}, []byte{}, []byte{}, []byte{}}
	result := []byte{}
	for i := 0; i < 4; i++ {
		for j := 0; j < Nb; j++ {
			state[i] = append(state[i], b[i+4*j])
		}
	}

	keySchedule := KeyExpansion(key)
	state = AddRoundKey(state, keySchedule, 0)

	for round := 1; round < Nk; round++ {
		state = SubBytes(state)
		state = ShiftRows(state)
		state = MixColumns(state)
		state = AddRoundKey(state, keySchedule, round)
	}

	state = SubBytes(state)
	state = ShiftRows(state)
	state = AddRoundKey(state, keySchedule, Nr)

	for i := range state {
		for j := range state[i] {
			result = append(result, state[i][j])
		}
	}

	return result
}

func InvShiftRows(state [][]byte) [][]byte {
	for counter, row := range state {
		for i := 0; i < counter; i++ {
			row = InvShiftSlice(row)
		}
	}
	return state
}

func InvShiftSlice(b []byte) []byte {
	a := []byte{}
	x, s :=	b[len(b)-1], b[:len(b)-1]
	a = append(a, x)
	a = append(a, s...)
	return a
}

func InvSubBytes(state [][]byte) [][]byte {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] = InvSBox(state[i][j])
		}
	}

	return state
}

func InvSBox(b byte) byte{
	for i := byte(0); i <= byte(255); i++ {
		if SBox(i) == b {
			return byte(i)
		}
	}
	return 0
}

func InvMixColumns(state [][]byte) [][]byte {
	for i := 0; i < Nb; i++ {
		state[0][i] = gField.MultiplyTwoBytes(state[0][i], 0x0e)^gField.MultiplyTwoBytes(state[1][i], 0x0b)^gField.MultiplyTwoBytes(state[2][i], 0x0d)^gField.MultiplyTwoBytes(state[3][i], 0x09)
        state[1][i] = gField.MultiplyTwoBytes(state[0][i], 0x09)^gField.MultiplyTwoBytes(state[1][i], 0x0e)^gField.MultiplyTwoBytes(state[2][i], 0x0b)^gField.MultiplyTwoBytes(state[3][i], 0x0d)
        state[2][i] = gField.MultiplyTwoBytes(state[0][i], 0x0d)^gField.MultiplyTwoBytes(state[1][i], 0x09)^gField.MultiplyTwoBytes(state[2][i], 0x0e)^gField.MultiplyTwoBytes(state[3][i], 0x0b)
        state[3][i] = gField.MultiplyTwoBytes(state[0][i], 0x0b)^gField.MultiplyTwoBytes(state[1][i], 0x0d)^gField.MultiplyTwoBytes(state[2][i], 0x09)^gField.MultiplyTwoBytes(state[3][i], 0x0e)
	} 
	return state
}

func Decryption(b []byte, key []byte) []byte {
	state := [][]byte{ []byte{}, []byte{}, []byte{}, []byte{}}
	result := []byte{}
	for i := 0; i < 4; i++ {
		for j := 0; j < Nb; j++ {
			state[i] = append(state[i], b[i+4*j])
		}
	}

	keySchedule := KeyExpansion(key)
	state = AddRoundKey(state, keySchedule, Nr)

	round := Nr - 1
	for round >= 1 {
		state = InvSubBytes(state)
		state = InvShiftRows(state)
		state = InvMixColumns(state)
		state = AddRoundKey(state, keySchedule, round)

		round--
	}

	state = InvSubBytes(state)
	state = InvShiftRows(state)
	state = AddRoundKey(state, keySchedule, round)

	for i := range state {
		for j := range state[i] {
			result = append(result, state[i][j])
		}
	}

	return result
}