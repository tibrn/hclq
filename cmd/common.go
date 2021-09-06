package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tibrn/hclq/config"
)

func getOutput(obj interface{}, raw bool) (string, error) {
	if raw {
		result := ""
		arr, isArray := obj.([]interface{})
		if isArray {
			for _, item := range arr {
				// Simple output, uses built-in %v, most useful for simple types.
				result += strings.Trim(fmt.Sprintf("%v", item), `"`) + " "
			}
			result = strings.TrimRight(result, " ")
			return result, nil
		}
		return strings.Trim(fmt.Sprintf("%v", obj), `"`), nil
	}
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBody), nil
}

// getInputReader provides an io.Reader for reading from either a file
// or stdin, depending on whether or not an input file was specified.
func getInputReader() (io.Reader, error) {
	if val := config.InputFile; val != "" {
		reader, err := os.Open(val)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
	return os.Stdin, nil
}
