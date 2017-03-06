package differ

import(
	"io/ioutil"
	"github.com/pmezard/go-difflib/difflib"
	"os"
	"github.com/garciademarina/checker/pkg/diffwrite"
)
func DiffFiles(a, b string) string {
	bf1, _ := ioutil.ReadFile(a)
	bf2, _ := ioutil.ReadFile(b)
	diff := difflib.ContextDiff{
		//diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(bf1)),
		B:        difflib.SplitLines(string(bf2)),
		FromFile: a,
		ToFile:   b,
		Context:  5,
		Eol:      "\n",
	}
	output, _ := difflib.GetContextDiffString(diff)
	diffwrite.WriteContextDiffHtml(os.Stdout,diff)
	return output
}
