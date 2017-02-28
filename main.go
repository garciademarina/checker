package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/pmezard/go-difflib/difflib"
)

const (
	lineWidth = 100
)

var (
	dir1 = flag.String("dir1", "1", "Directory 1")
	dir2 = flag.String("dir2", "2", "Directory 2")
	c    = color.New(color.FgCyan)
	r    = color.New(color.FgRed)
	g    = color.New(color.FgGreen)
	y    = color.New(color.FgYellow)
	bold = color.New(color.Bold)
	cbg  = color.New(color.FgBlack, color.BgWhite)
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: foobar \n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = func() { usage() }
	flag.Parse()

	fmt.Println(`
 _____ _           _                    ___   ___
|     | |_ ___ ___| |_ ___ ___    _ _  |   | |_  |
|   --|   | -_|  _| '_| -_|  _|  | | |_| | |_ _| |_
|_____|_|_|___|___|_,_|___|_|     \_/|_|___|_|_____|
                                                    `)

	fmt.Printf("Diff directories: %v %v\n", *dir1, *dir2)

	listfiles1, listfiles2, _ := createMap(*dir1, *dir2)

	parseDirectory(listfiles1, listfiles2)

}

func createMap(dir1, dir2 string) (map[string]string, map[string]string, error) {
	files1 := make(map[string]string)
	files2 := make(map[string]string)

	filepath.Walk(dir1, func(path string, f os.FileInfo, err error) error {
		key_path := path[len(dir1):len(path)]
		if key_path == "" {
			return nil
		}
		files1[key_path] = path
		return nil
	})
	filepath.Walk(dir2, func(path string, f os.FileInfo, err error) error {
		key_path := path[len(dir2):len(path)]
		if key_path == "" {
			return nil
		}
		files2[key_path] = path
		return nil
	})
	return files1, files2, nil
}

func parseDirectory(m1, m2 map[string]string) {
	m2_deleted := make([]string, 0, 50)
	m2_added := make([]string, 0, 50)

	for k := range m1 {
		fmt.Printf("File %v\n", k)
		if m2[k] == "" {
			m2_deleted = append(m2_deleted, strings.Trim(k, " "))
		} else {
			cont1, _ := ioutil.ReadFile(m1[k])
			cont2, _ := ioutil.ReadFile(m2[k])

			diff := difflib.ContextDiff{
				//diff := difflib.UnifiedDiff{
				A:        difflib.SplitLines(string(cont1)),
				B:        difflib.SplitLines(string(cont2)),
				FromFile: "Original",
				ToFile:   "Current",
				Context:  5,
				Eol:      "\n",
			}
			output, _ := difflib.GetContextDiffString(diff)
			//output, _ := difflib.GetUnifiedDiffString(diff)
			diffOption := false
			if output != "" {
				diffOption = true
			}

			action := scanAction(diffOption)
			switch action {
			case "c":
				printCodeWarnings(k, string(cont1))
			case "d":
				if diffOption {
					fmt.Printf("\n")
					_printLine(cbg, k)
					printColorDiff(output)
					_printLine(cbg, "end file "+k)
					fmt.Printf("\n")
				}
			case "":
				c.Printf("Skyped\n")
			}
		}
	}

	for k := range m2 {
		if m1[k] == "" {
			m2_added = append(m2_added, strings.Trim(k, " "))
		}
	}

	//for _, v := range m2_deleted {
	//	r.Printf("- %v\n", v)
	//}
	//for _, v := range m2_added {
	//	g.Printf("- %v\n", v)
	//}

	fmt.Printf("\nTotal files directory 1: %v\n", len(m1))
	fmt.Printf("Total files directory 2: %v\n", len(m1))
	r.Printf("\nMissed files: %v\n", len(m2_deleted))
	g.Printf("Added files: %v\n", len(m2_added))

}

func scanAction(diffOption bool) string {
	var option string
	// show help section
	fmt.Println("  c, check known errors")
	if diffOption {
		fmt.Println("  d, diff files")
	}
	fmt.Println("  ⏎ , do nothing (enter key)")

	fmt.Print("➜ ")
	fmt.Scanln(&option)
	return fmt.Sprintf(option)
}

func printColorDiff(s string) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			r.Println(line)
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			g.Println(line)
		}
		if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			bold.Println(line)
		}
		if strings.HasPrefix(line, "!") {
			y.Println(line)
		}
	}
}

func printCodeWarnings(filename, s string) {
	lines := difflib.SplitLines(string(s))
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
			lines_map = _printLineWarning(num, line, lines, lines_map)
		default:
		}
	}
	// make sure lines have red color
	var buffer bytes.Buffer
	for _, v := range linesWithError {
		line := lines[v]
		buffer.WriteString(cbg.Sprintf("%04d", v))
		buffer.WriteString(r.Sprintf("%v", line))
		if strings.Contains(line, "osc_esc_html") {
			buffer.Reset()
			buffer.WriteString(cbg.Sprintf("%04d", v))
			buffer.WriteString(y.Sprintf("%v", line))
		}
		lines_map[v] = buffer.String()
		buffer.Reset()

	}

	if len(lines_map) == 0 {
		ok := emoji.Sprint(":sunny:")
		g.Printf("No known error found %v\n", ok)
	} else {
		_printLineFill(c, "-", filename)
		for _, v := range lines_map {
			if v != "" {
				fmt.Print(v)
			}
		}
		_printLineFill(c, "-", "end file "+filename)
	}

}

func _printLine(cc *color.Color, filename string) {

	llen := lineWidth
	llen = llen - len(filename)

	for i := 0; i < llen/2; i++ {
		cc.Print(" ")
	}
	cc.Print(filename)
	for i := 0; i < llen/2; i++ {
		cc.Print(" ")
	}
	fmt.Printf("\n")
}

func _printLineWarning(n int, line string, lines []string, lines_map []string) []string {

	var buffer bytes.Buffer
	for i := 3; i >= 1; i-- {
		pre_n := n - i
		if pre_n >= 0 || pre_n < len(lines) {
			pre_line := lines[pre_n]
			buffer.WriteString(cbg.Sprintf("%04d", pre_n))
			buffer.WriteString(fmt.Sprintf("%v", pre_line))
			lines_map[pre_n] = buffer.String()
			buffer.Reset()
		}
	}
	buffer.WriteString(cbg.Sprintf("%04d", n))
	buffer.WriteString(r.Sprintf("%v", line))
	lines_map[n] = buffer.String()
	buffer.Reset()
	for i := 1; i <= 3; i++ {
		pre_n := n + i
		if pre_n >= 0 || pre_n < len(lines) {
			pre_line := lines[pre_n]
			buffer.WriteString(cbg.Sprintf("%04d", pre_n))
			buffer.WriteString(fmt.Sprintf("%v", pre_line))
			lines_map[pre_n] = buffer.String()
			buffer.Reset()
		}
	}
	return lines_map
}

func _printLineFill(cc *color.Color, char, filename string) {
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
