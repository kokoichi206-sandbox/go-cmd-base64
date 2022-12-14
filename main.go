package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// Print usage
var Usage = func() {
	fmt.Println("Usage:  base64 [-hd] <file_name>")
	flag.PrintDefaults()
}

// Options
type Params struct {
	IsHelp   bool
	IsDecode bool
	Args     []string
}

func init() {
	flag.BoolVarP(&params.IsHelp, "help", "h", false, "display this message")
	flag.BoolVarP(&params.IsDecode, "decode", "d", false, "decodes input")

	flag.Parse()

	params.Args = flag.Args()
}

var params Params

func main() {
	if params.IsHelp {
		Usage()
		os.Exit(0)
	}

	var buf []byte
	if len(params.Args) != 0 {
		// 引数にファイル名が渡された場合。
		buf = readFile(params.Args[0])
	} else {
		// なければ、標準入力を待ち受ける。
		buf = readStandardInput()
	}

	if len(buf) == 0 {
		fmt.Println("no input was passed")
		Usage()
		os.Exit(1)
	}

	if params.IsDecode {
		fmt.Print(Decode(buf))
	} else {
		fmt.Println(Encode(buf))
	}
}

func readFile(fileName string) []byte {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to open '%s': No such file or directory", fileName)
		os.Exit(1)
	}
	return buf
}

func readStandardInput() []byte {
	scanner := bufio.NewScanner(os.Stdin)
	// 1行分スキャン
	scanner.Scan()
	return scanner.Bytes()
}

// 渡されたバイト列をエンコードする。
func Encode(buf []byte) string {

	var builder strings.Builder

	max := len(buf) / 3
	if len(buf)%3 != 0 {
		max += 1
	}

	// 最終行以外をエンコード。
	for i := 0; i < max-1; i++ {
		b := buf[3*i : 3*i+3]

		if e := encode3bytes(b, 3, &builder); e != nil {
			fmt.Println("unexpected error occured while encoding")
			os.Exit(1)
		}
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

	return builder.String()
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

// 渡されたバイト列ををデコードする。
func Decode(buf []byte) string {

	sb := string(buf)
	max := len(sb) / 4

	var res bytes.Buffer
	// 最終行以外をエンコード。
	for i := 0; i < max-1; i++ {
		b := sb[4*i : 4*i+4]

		data := alphabetMap[string(b[3])] + alphabetMap[string(b[2])]*64 + alphabetMap[string(b[1])]*64*64 + alphabetMap[string(b[0])]*64*64*64

		b1 := data / (256 * 256)
		b2 := (data % (256 * 256)) / 256
		b3 := data % 256
		res.Write([]byte{byte(b1), byte(b2), byte(b3)})
	}

	b := sb[4*(max-1) : 4*max]
	data := alphabetMap[string(b[3])] + alphabetMap[string(b[2])]*64 + alphabetMap[string(b[1])]*64*64 + alphabetMap[string(b[0])]*64*64*64
	b1 := data / (256 * 256)
	b2 := (data % (256 * 256)) / 256
	b3 := data % 256
	if b2 == 0 {
		res.Write([]byte{byte(b1)})
	} else if b3 == 0 {
		res.Write([]byte{byte(b1), byte(b2)})
	} else {
		res.Write([]byte{byte(b1), byte(b2), byte(b3)})
	}

	return res.String()
}
