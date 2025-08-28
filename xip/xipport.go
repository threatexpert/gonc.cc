package xip

import "encoding/binary"

// MurmurHash3 x86_32 实现（seed 可设为 0）
func murmur3_x86_32(data []byte, seed uint32) uint32 {
	var (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
	)

	length := len(data)
	h1 := seed
	var k1 uint32
	nblocks := length / 4

	// body
	for i := 0; i < nblocks; i++ {
		// little endian read
		k1 = binary.LittleEndian.Uint32(data[i*4 : i*4+4])

		k1 *= c1
		k1 = (k1 << 15) | (k1 >> (32 - 15))
		k1 *= c2

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> (32 - 13))
		h1 = h1*5 + 0xe6546b64
	}

	// tail
	tail := data[nblocks*4:]
	var k uint32
	switch len(tail) {
	case 3:
		k ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k ^= uint32(tail[0])
		k *= c1
		k = (k << 15) | (k >> (32 - 15))
		k *= c2
		h1 ^= k
	}

	// finalization
	h1 ^= uint32(length)
	// fmix32
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}
