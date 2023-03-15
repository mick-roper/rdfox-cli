package ttl

import (
	"bytes"
	"io"
)

type triples map[string]map[string]string

func Write(data triples, dst io.Writer) error {
	var buffer bytes.Buffer
	for s, duples := range data {
		buffer.WriteString(s)

		var x int

		for p, o := range duples {
			x++

			buffer.WriteString("\n\t")
			buffer.WriteString(p)
			buffer.WriteRune('\t')
			buffer.WriteString(o)
			buffer.WriteRune(' ')

			if x < len(duples) {
				buffer.WriteRune(';')
			} else {
				buffer.WriteRune('.')
			}
		}

		buffer.WriteRune('\n')

		if _, err := dst.Write(buffer.Bytes()); err != nil {
			return err
		}

		buffer.Reset()
	}

	return nil
}
