package file_explorer

import (
	"fmt"
	"github.com/pigfall/tzzGoUtil/process"
	//"path/filepath"
	"strings"
)

// path format : C:/Users
func OpenExplorer(path string) error {
	fmt.Println(strings.ReplaceAll(path, "/", "\\"))
	//_, errOut, err := process.ExeOutput("powershell.exe", "-c", fmt.Sprintf("(cd %s) -or (explorer.exe .)", path))
	_, errOut, err := process.ExeOutput("explorer.exe", strings.ReplaceAll(path, "/", "\\"))
	if err != nil {
		return fmt.Errorf("%w,%s", err, errOut)
	}
	return nil
}

func SupportOpenExplorer() bool {
	return true
}
