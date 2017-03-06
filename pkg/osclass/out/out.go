package out

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	lineWidth = 100
)

var (
	C    = color.New(color.FgCyan)
	R    = color.New(color.FgRed)
	G   = color.New(color.FgGreen)
	Y    = color.New(color.FgYellow)
	Bold = color.New(color.Bold)
	Cbg  = color.New(color.FgBlack, color.BgWhite)
)


func PrintLine(cc *color.Color, filename string) {

	llen := lineWidth
	llen = llen - len(filename)

	cc.Print(filename)
	for i := len(filename); i < llen+len(filename); i++ {
		cc.Print(" ")
	}
	fmt.Printf("\n")
}

func PrintLineFill(cc *color.Color, char, filename string) {
	llen := lineWidth
	llen = llen - len(filename)

	for i := 0; i < llen/2; i++ {
		cc.Print(char)
	}
	cc.Print(filename)
	for i := 0; i < llen/2; i++ {
		cc.Print(char)
	}
	fmt.Printf("\n")
}

func PrintLineWarning(n int, line string, lines []string, lines_map []string) []string {

	var buffer bytes.Buffer
	for i := 3; i >= 1; i-- {
		pre_n := n - i
		if pre_n >= 0 || pre_n < len(lines) {
			pre_line := lines[pre_n]
			buffer.WriteString(Cbg.Sprintf("%04d", pre_n))
			buffer.WriteString(fmt.Sprintf("%v", pre_line))
			lines_map[pre_n] = buffer.String()
			buffer.Reset()
		}
	}
	buffer.WriteString(Cbg.Sprintf("%04d", n))
	buffer.WriteString(R.Sprintf("%v", line))
	lines_map[n] = buffer.String()
	buffer.Reset()
	for i := 1; i <= 3; i++ {
		pre_n := n + i
		if pre_n >= 0 || pre_n < len(lines) {
			pre_line := lines[pre_n]
			buffer.WriteString(Cbg.Sprintf("%04d", pre_n))
			buffer.WriteString(fmt.Sprintf("%v", pre_line))
			lines_map[pre_n] = buffer.String()
			buffer.Reset()
		}
	}
	return lines_map
}

func PrintStats(data map[string][]string) {
	// show digest
	aux := [][]string{
		[]string{"ADDED", strconv.Itoa(len(data["added"]))},
		[]string{"REMOVED", strconv.Itoa(len(data["removed"]))},
		[]string{"MODIFIED", strconv.Itoa(len(data["diff"]))},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"TYPE", "HITS"})

	for _, v := range aux {
		table.Append(v)
	}

	Bold.Printf("\nOverview\n")
	table.Render() // Send output
	Bold.Printf("\nFiles\n")
	if len(data["added"])>0 {
		for _, file := range data["added"] {
			G.Printf(" + %v\n", file)
		}
	}
	if len(data["removed"])>0 {
		for _, file := range data["removed"] {
			R.Printf(" - %v\n", file)
		}
	}
	if len(data["diff"])>0 {
		for _, file := range data["diff"] {
			Y.Printf(" · %v\n", file)
		}
	}


}