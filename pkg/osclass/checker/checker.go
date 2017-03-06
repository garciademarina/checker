package checker

import (
	"bytes"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"strings"
	"github.com/garciademarina/checker/pkg/osclass/out"
)

func ShowWarnings(s, filename string) int {

	bf1, _ := ioutil.ReadFile(s)
	lines := difflib.SplitLines(string(bf1))

	lines_map := make([]string, len(lines), len(lines))
	linesWithError := make([]int, 0, 0)

	for num, line := range lines {
		// input warning, a and other elements warning, img warning
		switch {
		case strings.Contains(line, "<?php") && (strings.Contains(line, "value=") ||
			strings.Contains(line, "placeholder=") ||
			strings.Contains(line, "title=") ||
			strings.Contains(line, "alt=") ||
			strings.Contains(line, "content=")):

			linesWithError = append(linesWithError, num)
			lines_map = out.PrintLineWarning(num, line, lines, lines_map)
		default:
		}
	}
	// make sure lines have red color
	var buffer bytes.Buffer
	for _, v := range linesWithError {
		line := lines[v]
		buffer.WriteString(out.Cbg.Sprintf("%04d", v))
		buffer.WriteString(out.R.Sprintf("%v", line))
		if strings.Contains(line, "osc_esc_html") {
			buffer.Reset()
			buffer.WriteString(out.Cbg.Sprintf("%04d", v))
			buffer.WriteString(out.Y.Sprintf("%v", line))
		}
		lines_map[v] = buffer.String()
		buffer.Reset()

	}

	if len(lines_map) == 0 {
		ok := emoji.Sprint(":sunny:")
		out.G.Printf("No known error found %v\n", ok)
	} else {
		out.PrintLineFill(out.C, "-", filename)
		for _, v := range lines_map {
			if v != "" {
				fmt.Print(v)
			}
		}
		out.PrintLineFill(out.C, "-", "end file "+filename)
	}
	return 1
}
