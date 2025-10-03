package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

var dbPath string

func init() {
	flag.StringVar(&dbPath, "db", "tokens.db", "path to BoltDB file")
}

func tokenToIP(tok uint32) string {
	b1 := byte((tok >> 16) & 0xFF)
	b2 := byte((tok >> 8) & 0xFF)
	b3 := byte(tok & 0xFF)
	return fmt.Sprintf("127.%d.%d.%d", b1, b2, b3)
}

func main() {
	flag.Parse()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("database file not found: %s", dbPath)
	}

	fmt.Println("opening DB...")

	db, err := bolt.Open(dbPath, 0600, &bbolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	fmt.Println("listing DB...")

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tokens"))
		if b == nil {
			fmt.Println("no tokens bucket found")
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			// k = token (big endian uint32)
			if len(k) != 4 {
				return nil
			}
			tok := binary.BigEndian.Uint32(k)
			ip := tokenToIP(tok)

			// v = serialized entry: host|port|expiresAtUnix
			// 这里简单假设你存储时是 "host:port|unixTime"
			data := string(v)

			var hostPort string
			var expireUnix int64
			n, _ := fmt.Sscanf(data, "%s|%d", &hostPort, &expireUnix)
			if n < 2 {
				fmt.Printf("token=%d ip=%s raw=%s\n", tok, ip, data)
				return nil
			}

			exp := time.Unix(expireUnix, 0)
			ttlRemain := time.Until(exp).Round(time.Second)

			fmt.Printf("token=%d ip=%s hostPort=%s expireAt=%s (remain=%s)\n",
				tok, ip, hostPort, exp.Format(time.RFC3339), ttlRemain)
			return nil
		})
	})
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
}
