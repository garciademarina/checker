package diffwrite

import (
	"strings"
	"fmt"
	"html"
	"regexp"
	"io"
	"bufio"
	"github.com/fatih/color"
	"github.com/garciademarina/checker/pkg/osclass/out"
	. "github.com/pmezard/go-difflib/difflib"
)

var (
	c    = color.New(color.FgCyan)
	r    = color.New(color.FgRed)
	g    = color.New(color.FgGreen)
	y    = color.New(color.FgYellow)
	bold = color.New(color.Bold)
	cbg  = color.New(color.FgBlack, color.BgWhite)
)

func ColorDiff(output, k string) {
	out.PrintLine(cbg, k)
	printColorDiff(output)
	fmt.Printf("\n")
}
func HtmlColorDiff(output, k string) {
	out.PrintLine(cbg, k+"<br>")
	printHtmlColorDiff(output)
	out.PrintLine(cbg, " end file "+k+"<br>")
	fmt.Printf("<br>")
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

func printHtmlColorDiff(s string) {
	lines := strings.Split(s, "\n")
	fmt.Printf(`<div style="position: relative;
		margin-top: 16px;
		margin-bottom: 16px;
		border: 1px solid #ddd;
		border-radius: 3px;"</div>`)
	for _, line := range lines {
		re := regexp.MustCompile(`\r?\n`)
		line = re.ReplaceAllString(line, "")
		line = html.EscapeString(line)
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			r.Printf("<span style='background-color: #ffecec;'>%v</span><br>",line)
		} else {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				g.Printf("<span style='background-color: #eaffea;'>%v</span><br>", line)
			} else {
				if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
					bold.Printf("<span style='font-weight:bold;'>%v</span><br>", line)
				} else {
					if line == "!" {
						fmt.Printf("<span style=''>%v</span><br>", line);
					} else {
						if strings.HasPrefix(line, "!") {
							fmt.Printf("<span style='background-color: #ffffec;'>%v</span><br>", line)
						} else {
							fmt.Printf("<span style='background-color: #ecffec;'>%v</span><br>", line)
						}
					}
				}
			}
		}
	}
	fmt.Printf("</div>")
}


// define custom html output
// 2 columns layout
func WriteContextDiffHtml(writer io.Writer, diff ContextDiff) error {
	buf := bufio.NewWriter(writer)
	defer buf.Flush()
	var diffErr error
	wf := func(format string, args ...interface{}) {
		_, err := buf.WriteString(fmt.Sprintf(format, args...))
		if diffErr == nil && err != nil {
			diffErr = err
		}
	}
	ws := func(s string) {
		_, err := buf.WriteString(s)
		if diffErr == nil && err != nil {
			diffErr = err
		}
	}

	if len(diff.Eol) == 0 {
		diff.Eol = "\n"
	}

	prefix := map[byte]string{
		'i': "+ ",
		'd': "- ",
		'r': "! ",
		'e': "  ",
	}

	started := false
	m := NewMatcher(diff.A, diff.B)
	for _, g := range m.GetGroupedOpCodes(diff.Context) {
		if !started {
			started = true
			fromDate := ""
			if len(diff.FromDate) > 0 {
				fromDate = "\t" + diff.FromDate
			}
			toDate := ""
			if len(diff.ToDate) > 0 {
				toDate = "\t" + diff.ToDate
			}
			if diff.FromFile != "" || diff.ToFile != "" {
				wf("<div style='border: 1px solid #ddd;padding: 5px 10px;background-color: #fafbfc;border-bottom: 1px solid #e1e4e8;border-top-left-radius: 2px;border-top-right-radius: 2px;'>%s%s%s", diff.FromFile, fromDate, diff.Eol)
				wf("--- %s%s%s</div>", diff.ToFile, toDate, diff.Eol)
			}
		}

		//first, last := g[0], g[len(g)-1]
		first, _ := g[0], g[len(g)-1]
		// The tags are characters, with these meanings:
		//
		// 'r' (replace):  a[i1:i2] should be replaced by b[j1:j2]
		//
		// 'd' (delete):   a[i1:i2] should be deleted, j1==j2 in this case.
		//
		// 'i' (insert):   b[j1:j2] should be inserted at a[i1:i1], i1==i2 in this case.
		//
		// 'e' (equal):    a[i1:i2] == b[j1:j2]
		var index1 int
		//range1 := formatRangeContext(first.I1, last.I2)
		//wf("*** %s ****%s", range1, diff.Eol)
		ws("<div>")

		ws("<table style='border:1px black solid;width:50%;float:left;'>")
		ws("	<tbody>")
		for _, c := range g {
			if c.Tag == 'r' || c.Tag == 'd' {
				for index, cc := range g {
					index1++
					if cc.Tag == 'i' {
						ws("	    <tr>")
						ws("	        <td style='width: 1%;min-width: 50px;'>")
						wf("%d",first.I1+index1+1);
						ws("	        </td>")
						ws("	        <td>&nbsp;</td>")
						ws("	    </tr>")
						continue
					}
					for _, line := range diff.A[cc.I1:cc.I2] {
						line = html.EscapeString(line)
						ws("	    <tr>")
						ws("	        <td style='width: 1%;min-width: 50px;'>")
						ws(string(first.I1+index));
						wf("%d",first.I1+index1+1);
						ws("	        </td>")
						ws("	        <td>")
						if cc.Tag == 'r' {
							ws("<span style='background-color: #ffffec;'>");
						}
						if cc.Tag == 'd' {
							ws("<span style='background-color: #ffecec;'>");
						}
						if cc.Tag == 'i' {
							ws("<span style='background-color: #eaffea;'>");
						}
						ws(prefix[cc.Tag] + line )
						if cc.Tag == 'i' || cc.Tag == 'd' || cc.Tag == 'r' {
							ws("</span>")
						}

						ws("	        </td>")
						ws("	    </tr>")
					}
				}
				break
			}
		}
		ws("	</tbody>")
		ws("<table>")

		var index2 int
		//range2 := formatRangeContext(first.J1, last.J2)
		//wf("--- %s ----%s", range2, diff.Eol)
		ws("<table style='border:1px black solid;width:50%;'>")
		ws("	<tbody>")
		for _, c := range g {
			if c.Tag == 'r' || c.Tag == 'i' {
				for _, cc := range g {
					index2++
					if cc.Tag == 'd' {
						ws("	    <tr>")
						ws("	        <td style='width: 1%;min-width: 50px;'>")
						wf("%d",first.J1+index2+1);
						ws("	        </td>")
						ws("	        <td>&nbsp;</td>")
						ws("	    </tr>")
						continue
					}
					for _, line := range diff.B[cc.J1:cc.J2] {
						line = html.EscapeString(line)
						ws("	    <tr>")
						ws("	        <td style='width: 1%;min-width: 50px;'>")
						wf("%d",first.J1+index2+1);
						ws("	        </td>")
						ws("	        <td>")
						if cc.Tag == 'r' {
							ws("<span style='background-color: #ffffec;'>");
						}
						if cc.Tag == 'd' {
							ws("<span style='background-color: #ffecec;'>");
						}
						if cc.Tag == 'i' {
							ws("<span style='background-color: #eaffea;'>");
						}
						ws(prefix[cc.Tag] + line )
						if cc.Tag == 'i' || cc.Tag == 'd' || cc.Tag == 'r' {
							ws("</span>")
						}
						ws("	        </td>")
						ws("	    </tr>")
					}
				}
				break
			}
		}
		ws("	</tbody>")
		ws("</table>")
		ws("</div><div style='clear:both;'></div>")
	}
	return diffErr
}

// Convert range to the "ed" format.
func formatRangeContext(start, stop int) string {
	// Per the diff spec at http://www.unix.org/single_unix_specification/
	beginning := start + 1 // lines start numbering with one
	length := stop - start
	if length == 0 {
		beginning -= 1 // empty ranges begin at line just before the range
	}
	if length <= 1 {
		return fmt.Sprintf("%d", beginning)
	}
	return fmt.Sprintf("%d,%d", beginning, beginning+length-1)
}


