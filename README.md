# go-cmd-base64

## Usage

``` sh
$ go run main.go -h
Usage:  base64 [-h] <file_name>
  -h, --help   display this message

$ make build

# Encode
$ cat test.txt
ab% 
$ ./base64 test.txt
YWI=

# Decode
$ ./base64 test.txt > encoded_test.txt
$ ./base64 -d encoded_test.txt
ab%
```

## LICENSE

under [MIT License](./LICENSE).
