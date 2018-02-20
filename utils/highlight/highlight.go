package highlight

import (
	"strings"

	"github.com/mgutz/ansi"
)

func Wrap(output string, keywords []string) string {
	for _, keyword := range keywords {
		output = strings.Replace(output, keyword, ansi.Color(keyword, "red"), -1)
	}

	return output
}