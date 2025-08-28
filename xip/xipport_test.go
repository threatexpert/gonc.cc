package xip

import (
	"fmt"
	"testing"

	"github.com/spaolacci/murmur3"
)

func TestMurmur3Compatibility(t *testing.T) {
	inputs := [][]byte{
		[]byte("hello"),
		[]byte("127.0.0.1-3005"),
		[]byte("127.0.0.1-6799"),
		[]byte("192.168.0.1-65535"),
		[]byte("test123"),
		[]byte(""),
	}

	for _, in := range inputs {
		want := murmur3.Sum32(in)
		got := murmur3_x86_32(in, 0)
		if want != got {
			t.Errorf("Mismatch for input %q: want=%d got=%d", in, want, got)
		}
	}
}

func TestLow24Collision(t *testing.T) {
	seen := make(map[uint32]string)
	collisions := 0

	for port := 1; port <= 65535; port++ {
		s := fmt.Sprintf("127.0.0.1-%d", port)
		h := murmur3.Sum32([]byte(s))
		low24 := h & 0xFFFFFF // 取低 24 位

		if prev, exists := seen[low24]; exists {
			collisions++
			if collisions <= 10 { // 只打印前 10 个碰撞案例
				t.Logf("Collision #%d: %q and %q have same low24=%06x", collisions, prev, s, low24)
			}
		} else {
			seen[low24] = s
		}
	}

	t.Logf("Total collisions: %d", collisions)
	if collisions > 0 {
		t.Errorf("Found %d collisions in low24 space", collisions)
	}
}

func TestLow24Collision19216880(t *testing.T) {
	collisions := make(map[uint32]string) // low24 -> first string
	count := 0

	for i := 0; i <= 0xFFFF; i++ { // 192.168.0.0 ~ 192.168.255.255
		ip := fmt.Sprintf("192.168.%d.%d-80", byte(i>>8), byte(i&0xFF))
		hash := murmur3.Sum32([]byte(ip))
		low24 := hash & 0xFFFFFF

		if prev, exists := collisions[low24]; exists {
			count++
			fmt.Printf("Collision #%d: %q and %q have same low24=%06x\n", count, prev, ip, low24)
		} else {
			collisions[low24] = ip
		}
	}

	fmt.Printf("Total collisions in 192.168.0.0/16 with port 80: %d\n", count)
}
