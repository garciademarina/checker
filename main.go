package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/garciademarina/checker/pkg/diffwrite"
	"github.com/garciademarina/checker/pkg/osclass/checker"
	"github.com/garciademarina/checker/pkg/osclass/differ"
	"github.com/garciademarina/checker/pkg/osclass/out"
)

const (
	lineWidth = 100
)

var (
	dir1   = flag.String("dir1", "1", "Directory 1")
	dir2   = flag.String("dir2", "2", "Directory 2")
	action = flag.String("f", "", "Show stats overview")
	c      = color.New(color.FgCyan)
	r      = color.New(color.FgRed)
	g      = color.New(color.FgGreen)
	y      = color.New(color.FgYellow)
	bold   = color.New(color.Bold)
	cbg    = color.New(color.FgBlack, color.BgWhite)
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: foobar \n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = func() { usage() }
	flag.Parse()



	listfiles1, listfiles2, _ := createMap(*dir1, *dir2)

	if *action == "stats" {
		stats := overview(listfiles1, listfiles2)
		out.PrintStats(stats)
		os.Exit(0)
	}
	if *action == "diffall" {
		diffall(listfiles1, listfiles2)
		os.Exit(0)
	}

	fmt.Println(`
 _____ _           _                    ___   ___
|     | |_ ___ ___| |_ ___ ___    _ _  |   | |_  |
|   --|   | -_|  _| '_| -_|  _|  | | |_| | |_ _| |_
|_____|_|_|___|___|_,_|___|_|     \_/|_|___|_|_____|
                                                    `)

	fmt.Printf("Diff directories: %v %v\n", *dir1, *dir2)

	var option string

	help_promt := func() string {
		fmt.Println("stats, Stats overview")
		fmt.Println("manual, Manualy check all files")
		fmt.Println("auto, Check all files automatically (diff+warnings)")
		fmt.Print("➜ ")
		fmt.Scanln(&option)
		return option
	}

	// Scan user input
	invalid := true
	for invalid {
		invalid = false
		action := help_promt()
		switch action {
		case "stats":
			stats := overview(listfiles1, listfiles2)
			out.PrintStats(stats)
		case "manual":
			parseDirectory(listfiles1, listfiles2)
		case "auto":

		default:
			out.R.Printf("No valid input, try again.\n")
			invalid = true
		}
	}
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

func overview(m1, m2 map[string]string) map[string][]string {
	ret := make(map[string][]string)

	for filename := range m1 {

		if m2[filename] == "" {
			ret["removed"] = append(ret["removed"], filename)
		} else {
			if differ.DiffFiles(m1[filename], m2[filename]) != "" {
				ret["diff"] = append(ret["diff"], m2[filename])
			}
		}
		//checker.ShowWarnings(m2[filename], filename);
	}
	for filename := range m2 {
		if m1[filename] == "" {
			ret["added"] = append(ret["added"], m2[filename])
		}
	}

	return ret
}

func parseDirectory(m1, m2 map[string]string) {
	for filename := range m1 {

		if m2[filename] != "" {
			// Diff files
			var diffOption bool
			// map[filename] = path_file_name
			output := differ.DiffFiles(m1[filename], m2[filename])
			if output != "" {
				diffOption = true
			}
			// Scan user input
			invalid := true
			for invalid {
				fmt.Printf("File %v\n", filename)
				invalid = false
				action := scanAction(diffOption)
				switch action {
				case "c":
					checker.ShowWarnings(m2[filename], filename)
				case "d":
					diffwrite.ColorDiff(output, filename)
				case "":
					c.Printf("Skyped\n")
				default:
					out.R.Printf("No valid input, try again.\n")
					invalid = true
				}
			}
		}
	}

}

func diffall(m1, m2 map[string]string) {
	header := `<html>
	<head>
	</head>
	<body>`
	footer := `</body>
	<html>`
	fmt.Print(header)
	fmt.Println(`
		<pre>
 _____ _           _                    ___   ___
|     | |_ ___ ___| |_ ___ ___    _ _  |   | |_  |
|   --|   | -_|  _| '_| -_|  _|  | | |_| | |_ _| |_
|_____|_|_|___|___|_,_|___|_|     \_/|_|___|_|_____|
</pre>`)
	for filename := range m1 {

		if m2[filename] != "" {
			// map[filename] = path_file_name
			differ.DiffFiles(m1[filename], m2[filename])
			//if output != "" {
			//	diffwrite.HtmlColorDiff(output, filename)
			//}
		}
	}
	fmt.Print(footer)
}

func scanAction(diffOption bool) string {
	var option string
	// show help section
	fmt.Println("  c, check known errors")
	if diffOption {
		fmt.Println("  d, diff files")
	} else {
		fmt.Println("  d, diff files (Without changes)")
	}
	fmt.Println("  ⏎ , do nothing (enter key)")

	fmt.Print("➜ ")
	fmt.Scanln(&option)
	return fmt.Sprintf(option)
}
