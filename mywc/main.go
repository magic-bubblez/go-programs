package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Count struct {
	Lines int
	Words int
	Bytes int
}

func (c *Count) Add(other Count) {
	c.Lines += other.Lines
	c.Words += other.Words
	c.Bytes += other.Bytes
}

func count(r io.Reader) (Count, error) {
	c := Count{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c.Lines++
		c.Words += len(strings.Fields(line))
		c.Bytes += len(line) + 1
	}
	if err := scanner.Err(); err != nil {
		return Count{}, fmt.Errorf("reading input: %w", err)
	}

	return c, nil
}

func countFile(filename string) (Count, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Count{}, fmt.Errorf("countFile %s: %w", filename, err)
	}
	defer f.Close()
	return count(f)
}

func printCount(w io.Writer, c Count, name string, showLines, showWords, showBytes bool) {
	if !showLines && !showWords && !showBytes {
		showLines, showWords, showBytes = true, true, true
	}
	if showLines {
		fmt.Fprintf(w, "%8d", c.Lines)
	}
	if showWords {
		fmt.Fprintf(w, "%8d", c.Words)
	}
	if showBytes {
		fmt.Fprintf(w, "%8d", c.Bytes)
	}
	if name != "" {
		fmt.Fprintf(w, " %s", name)
	}
	fmt.Fprintln(w)
}

func main() {
	showLines := flag.Bool("l", false, "count lines")
	showWords := flag.Bool("w", false, "count words")
	showBytes := flag.Bool("c", false, "count bytes")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		c, err := count(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mywc: %v\n", err)
			os.Exit(1)
		}
		printCount(os.Stdout, c, "", *showLines, *showWords, *showBytes)
		return
	}
	total := Count{}

	for _, filename := range files {
		c, err := countFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mywc: %v\n", err)
			continue
		}
		printCount(os.Stdout, c, filename, *showLines, *showWords, *showBytes)
		total.Add(c)
	}
}
