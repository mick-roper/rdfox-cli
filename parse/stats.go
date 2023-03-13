package parse

import (
	"bufio"
	"io"
	"strings"
)

type statistics map[string]map[string]interface{}

func Stats(r io.Reader) statistics {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	stats := statistics{}
	thisComponent := ""
	scanner.Split(bufio.ScanLines)

	i := -1

	for scanner.Scan() {
		i++
		if i == 0 {
			continue
		}

		t := scanner.Text()

		parts := strings.SplitN(t, "\t", 3)
		_, p, v := parts[0], strings.Trim(parts[1], "\""), strings.Trim(parts[2], "\"")

		if p == "Component name" {
			thisComponent = v
			stats[thisComponent] = map[string]interface{}{}
			continue
		}

		stats[thisComponent][p] = v
	}

	return stats
}
