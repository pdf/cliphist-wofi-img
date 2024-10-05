package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const tmpDir = `/tmp/cliphist-wofi-img`

var imgRegexp = regexp.MustCompile(`^\[\[ binary.*\b(?P<ext>jpg|jpeg|png|bmp)\b`)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, `Usage:`, os.Args[0], `<cliphistEntry>`)
		os.Exit(1)
	}

	input := os.Args[1]
	split := strings.SplitN(input, "\t", 2)

	if len(split) != 2 {
		fmt.Fprintf(os.Stderr, "Incorrect number of fields (wanted 2): got %d\n", len(split))
		fmt.Println(input)
		os.Exit(1)
	}

	matches := imgRegexp.FindStringSubmatch(split[1])
	if len(matches) < 2 {
		fmt.Fprintln(os.Stderr, "Image not found, exiting")
		fmt.Println(input)
		os.Exit(0)
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Fprintln(os.Stderr, `Failed removing tmp dir:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		fmt.Fprintln(os.Stderr, `Failed creating tmp dir:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	cmd := exec.Command(`cliphist`, `decode`)

	suffix := matches[1]
	imgPath := filepath.Join(tmpDir, split[0]+`.`+suffix)
	f, err := os.Create(imgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, `Failed creating output image:`, err)
		fmt.Println(input)
		os.Exit(1)
	}
	defer f.Close()

	in, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, `Failed connecting stdin:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	go func() {
		defer in.Close()
		if _, err := io.WriteString(in, input); err != nil {
			fmt.Fprintln(os.Stderr, `Failed sending input:`, err)
			fmt.Println(input)
			os.Exit(1)
		}
	}()

	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, `Failed executing cliphist:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	if _, err := f.Write(out); err != nil {
		fmt.Fprintln(os.Stderr, `Failed writing image:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	fmt.Println(`img:` + imgPath + `:text:` + input)
}
