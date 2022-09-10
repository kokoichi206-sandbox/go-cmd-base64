package main_test

import (
	main "base64"
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

		// Act
		result := main.Encode(fileName)

		// Assert
		require.Equal(t, expected, result)
	}
}
