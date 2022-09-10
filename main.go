package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	padding  = "="
)

func main() {
	fileName := "test.txt"
	Encode(fileName)
}

// 指定されたファイルの中身をエンコードする。
func Encode(fileName string) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to open '%s': No such file or directory", fileName)
		os.Exit(1)
	}
	var builder strings.Builder

	max := len(buf) / 3
	if len(buf)%3 != 0 {
		max += 1
	}

	cnt := 0
	// 最終行以外をエンコード。
	for i := 0; i < max-1; i++ {
		b := buf[3*i : 3*i+3]
		// fmt.Println(b)

		if e := encode3bytes(b, 3, &builder); e != nil {
			fmt.Println("unexpected error occured while encoding")
			os.Exit(1)
		}
		cnt += 1
	}
	// 最終行をエンコード。
	switch len(buf) % 3 {
	case 0:
		b := buf[3*(max-1) : 3*max]
		if e := encode3bytes(b, 3, &builder); e != nil {
			fmt.Println("unexpected error occured while encoding")
			os.Exit(1)
		}
	case 1:
		b := append(buf[3*(max-1):3*max-2], make([]byte, 2)...)
		if e := encode3bytes(b, 1, &builder); e != nil {
			fmt.Println("unexpected error occured while encoding")
			os.Exit(1)
		}
	case 2:
		b := append(buf[3*(max-1):3*max-1], make([]byte, 1)...)
		if e := encode3bytes(b, 2, &builder); e != nil {
			fmt.Println("unexpected error occured while encoding")
			os.Exit(1)
		}
	}

	fmt.Println(builder.String())
}

// encode from 3 bytes binary data to 4 strings
func encode3bytes(bytes []byte, v int, builder *strings.Builder) error {

	data := int(bytes[2]) + int(bytes[1])*256 + int(bytes[0])*256*256

	tmp := [4]string{}
	for i := 0; i < 4; i++ {
		// 下から順に、8 bit ずつエンコードする
		tmp[i] = string(alphabet[data&63])
		data >>= 6
	}
	for i := 0; i < 4; i++ {
		// 逆順に詰めていく
		if i > v {
			fmt.Fprint(builder, "=")
		} else {
			fmt.Fprint(builder, tmp[3-i])
		}
	}
	return nil
}
