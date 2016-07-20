package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"os"
	"time"
)

var accountingFile string
var reportingFile string

func init() {
	sgeroot := os.Getenv("SGE_ROOT")
	sgecell := os.Getenv("SGE_CELL")
	if sgeroot == "" || sgecell == "" {
		fmt.Fprintf(os.Stderr, "$SGE_ROOT or $SGE_CELL not set")
		os.Exit(1)
	}

	accountingFile = fmt.Sprintf("%s/%s/common/accounting", sgeroot, sgecell)
	reportingFile = fmt.Sprintf("%s/%s/common/reporting", sgeroot, sgecell)

}

func main() {
	printToFiles := false
	if len(os.Args) > 1 {
		printToFiles = true
	}

	t, err := tail.TailFile(reportingFile, tail.Config{Follow: true})
	if err != nil {
		os.Exit(1)
	}

	for line := range t.Lines {
		if printToFiles {
			for i := 1; i < len(os.Args); i++ {
				if file, errOpen := os.OpenFile(os.Args[i], os.O_APPEND|os.O_WRONLY, 0600); errOpen != nil {
					if f, errCreate := os.Create(os.Args[i]); errCreate != nil {
						fmt.Printf("%s Error opening file %s: %s\n", time.Now(), os.Args[i], errOpen)
					} else {
						f.WriteString(line.Text + "\n")
						f.Close() // by purpose
					}
				} else {
					file.WriteString(line.Text + "\n")
					file.Close() // by purpose
				}
			}
		} else {
			fmt.Println(line.Text)
		}
	}
}
