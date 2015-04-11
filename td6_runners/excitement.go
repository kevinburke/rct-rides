package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/kevinburke/rct-rides/rle"
)

func main() {
	encodedBits, err := ioutil.ReadFile("rides/mischief.td6")
	fmt.Println(hex.EncodeToString(encodedBits[len(encodedBits)-4 : len(encodedBits)]))
	bitsWithoutChecksum := encodedBits[:len(encodedBits)-4]

	if err != nil {
		panic(err)
	}
	z := rle.NewReader(bytes.NewReader(bitsWithoutChecksum))
	if err != nil {
		panic(err)
	}

	var bitbuffer bytes.Buffer
	bitbuffer.ReadFrom(z)
	decrypted := bitbuffer.Bytes()

	// r is a pointer
	r := new(Ride)
	Unmarshal(decrypted, r)

	bits, err := Marshal(r)
	if err != nil {
		panic(err)
	}

	for i := range bits {
		if bits[i] != decrypted[i] {
			fmt.Printf("%d: ", i)
			fmt.Printf("byte %x differs, in my mischief it is %d, in orig it is %d\n", i, bits[i], decrypted[i])
		}
	}

	if DEBUG {
		begin := 0xa2 + 2*len(r.TrackData.Elements) - 3
		for i := begin; i < begin+DEBUG_LENGTH; i++ {
			// encode the value of i as hex
			ds := hex.EncodeToString([]byte{byte(i)})
			bitValueInHex := hex.EncodeToString([]byte{bits[i]})
			fmt.Printf("%s: %s\n", ds, bitValueInHex)
		}
	}

	paddedBits := pad(bits)

	fmt.Println(paddedBits)
	fmt.Println(decrypted)

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Write(paddedBits)
	ioutil.WriteFile("/Users/kevin/Applications/Wineskin/rct2.app/Contents/Resources/drive_c/GOG Games/RollerCoaster Tycoon 2 Triple Thrill Pack/Tracks/mymischief.TD6", buf.Bytes(), 0644)
	fmt.Println("Wrote rides/mischief.td6.out.")
}