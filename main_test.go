package main_test

import (
	main "base64"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	fileNames := []string{
		"test_data/test0.txt",
		"test_data/test1.txt",
		"test_data/cat.jpg",
	}

	for _, fileName := range fileNames {
		// Arrange
		data, err := exec.Command("base64", fileName).Output()
		require.NoError(t, err)
		// Without newline
		expected := string(data[:len(data)-1])
		// Encode の引数に渡すために、ファイルを読み込みバイト列にする。
		buf, err := ioutil.ReadFile(fileName)
		require.NoError(t, err)

		// Act
		result := main.Encode(buf)

		// Assert
		require.Equal(t, expected, result)
	}
}
