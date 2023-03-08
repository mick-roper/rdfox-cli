package parse

import (
	"bufio"
	"io"
	"log"
	"strings"
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
		t := scanner.Text()
		parts := strings.Split(t, "\t")

		log.Print(parts)
	}

	return result
}
