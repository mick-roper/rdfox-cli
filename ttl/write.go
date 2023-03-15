package ttl

import (
	"bytes"
	"io"
)

type triples map[string]map[string][]string

func Write(data triples, dst io.Writer) error {
	var buffer bytes.Buffer
	for s, duples := range data {
		buffer.WriteString(s)

		var i int
		var dCount int

		for _, objects := range duples {
			dCount += len(objects)
		}

		for p, objects := range duples {
			for _, o := range objects {
				i++
				buffer.WriteString("\n\t")
				buffer.WriteString(p)
				buffer.WriteRune('\t')
				buffer.WriteString(o)
				buffer.WriteRune('\t')

				if i < dCount {
					buffer.WriteRune(';')
				}
			}

		}

		buffer.WriteString(".\n")

		if _, err := dst.Write(buffer.Bytes()); err != nil {
			return err
		}

		buffer.Reset()
	}

	return nil
}
