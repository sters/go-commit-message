package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/koki-develop/go-fzf"
)

type commitType struct {
	key         string
	description string

	printable string
}

var commitTypes = []commitType{
	{key: "add", description: "Add something"},
	{key: "change", description: "Change something"},
	{key: "deprecate", description: "Mark as deprecate something"},
	{key: "remove", description: "Remove something"},
	{key: "hotfix", description: "Make a hotfix"},
	{key: "fix", description: "Just bug fix"},
	{key: "security", description: "Fix security related issue"},
	{key: "refactor", description: "Refactoring code"},
	{key: "chore", description: "Something not related the feature (e.g. CI/CD, tools, libraries)"},
	{key: "doc", description: "Documentation only changes"},
	{key: "test", description: "Adding or Fix testcase"},
	{key: "style", description: "Just update coding style that should be not changed a feature"},
	{key: "db", description: "DB only changes (e.g. migrations, SQL)"},
}

func main() {
	f, err := fzf.New()
	if err != nil {
		log.Fatal(err)
	}

	setPrintableCommitTypes()

	idxs, err := f.Find(
		commitTypes,
		func(i int) string {
			return commitTypes[i].printable
		},
	)
	if err != nil {
		if errors.Is(err, fzf.ErrAbort) {
			os.Exit(1)
		}
		log.Fatal(err)
	}
	if len(idxs) != 1 {
		fmt.Println("Please choose single item.")
		os.Exit(1)
	}

	choosen := commitTypes[idxs[0]]
	fmt.Printf("Type: %s\n", choosen.key)

	s := bufio.NewScanner(os.Stdin)
	fmt.Printf("Scope: ")
	s.Scan()
	scope := s.Text()

	fmt.Printf("Description: ")
	s.Scan()
	description := s.Text()

	fmt.Printf("\n%s(%s): %s\n", choosen.key, scope, description)
}

func setPrintableCommitTypes() {
	buf := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	for _, t := range commitTypes {
		fmt.Fprintf(w, "%s\t-\t%s\n", t.key, t.description)
	}
	w.Flush()

	for i, p := range strings.Split(buf.String(), "\n") {
		if i == len(commitTypes) {
			break
		}
		commitTypes[i].printable = p
	}
}
