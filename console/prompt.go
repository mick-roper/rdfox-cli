package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StringPrompt(label string) string {
	var s string
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(label)
		s, _ = reader.ReadString('\n')
		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s)
}

func BoolPrompt(label string) bool {
	const suffix = " (type 'yes' to continue)"
	if !strings.HasSuffix(label, suffix) {
		label = fmt.Sprint(label, suffix)
	}

	s := strings.ToLower(StringPrompt(label))

	return s == "yes"
}
