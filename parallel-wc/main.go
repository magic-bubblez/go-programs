package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Count struct {
	Lines int
	Words int
	Bytes int
}

type Result struct {
	Filename string
	Count    Count
	Err      error
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
	if name != "" {
		fmt.Fprintf(w, "Name: %s", name)
	}
	fmt.Fprintln(w)

	if showBytes {
		fmt.Fprintf(w, "\tBytes: %8d ", c.Bytes)
	}
	if showWords {
		fmt.Fprintf(w, "\tWords: %8d ", c.Words)
	}
	if showLines {
		fmt.Fprintf(w, "\tLines: %8d ", c.Lines)
	}
	fmt.Fprintln(w)
}

func main() {
	showLines := flag.Bool("l", false, "count lines")
	showWords := flag.Bool("w", false, "count words")
	showBytes := flag.Bool("c", false, "count bytes")
	flag.Parse()

	files := flag.Args()
	ch := make(chan Result)
	var wg sync.WaitGroup

	for _, filename := range files {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			c, err := countFile(name)
			ch <- Result{Filename: name, Count: c, Err: err}
		}(filename)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	total := Count{}
	for r := range ch {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "mywc: %v\n", r.Err)
			continue
		}
		printCount(os.Stdout, r.Count, r.Filename, *showLines, *showWords, *showBytes)
		total.Add(r.Count)
	}

	if len(files) > 1 {
		printCount(os.Stdout, total, "total", *showLines, *showWords, *showBytes)
	}
}
