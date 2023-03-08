package parse

import (
	"bufio"
	"io"
)

type statistics map[string]map[string]interface{}

func Stats(r io.Reader) statistics {
	result := statistics{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	const (
		stateIndex int = iota
		stateProperty
		stateValue
	)

	for scanner.Scan() {
		state := stateIndex
		delimitedString := false
		// var prevByte byte
		// var propName string
		// var value interface{}

		for _, b := range scanner.Bytes() {
			// prevByte = b

			switch b {
			case ' ':
				continue
			case '"':
				delimitedString = !delimitedString
			default:
				if state == stateValue {
					state = stateIndex
				} else {
					state++
				}
			}
		}
	}

	return result
}
