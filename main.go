package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const tmpDir = `/tmp/cliphist`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, `Usage:`, os.Args[0], `<cliphistEntry>`)
		os.Exit(1)
	}

	input := os.Args[1]
	split := strings.SplitN(input, "\t", 2)

	if len(split) != 2 {
		fmt.Println(input)
		os.Exit(1)
	}

	if strings.Index(split[1], `binary data image/`) != 0 {
		fmt.Println(input)
		os.Exit(1)
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Fprintln(os.Stderr, `Failed removing tmp dir:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	if err := os.Mkdir(tmpDir, 0700); err != nil {
		fmt.Fprintln(os.Stderr, `Failed creating tmp dir:`, err)
		fmt.Println(input)
		os.Exit(1)
	}

	cmd := exec.Command(`cliphist`, `decode`)

	suffix := split[1][strings.IndexRune(split[1], '/')+1:]
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

	//fmt.Println(split[0] + "\t" + `img:` + imgPath + `:text:` + input)
	fmt.Println(`img:` + imgPath + `:text:` + input)
}
