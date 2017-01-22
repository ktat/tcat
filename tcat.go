package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

func main() {
	format, opt, file := parseOptions()
	in := make(chan string)
	number := 1

	if file == "" {
		go readStdin(in)
	} else {
		var fp *os.File
		var ioerr error
		fp, ioerr = os.Open(file)
		if ioerr != nil {
			panic("cannot read file")
		}
		go readFile(in, fp)
	}

	for {
		l, ok := <-in
		if ok == false {
			break
		} else {
			fmt.Print(timedText(format, l, &number, opt))
		}
	}

}

func readStdin(in chan string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var s = scanner.Text()
		s += "\n"
		in <- s
	}
	close(in)
}

func readFile(in chan string, fp *os.File) {
	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, ioerr := reader.ReadLine()
		if ioerr != nil {
			if ioerr != io.EOF {
				panic(ioerr)
			} else if ioerr == io.EOF {
				break
			}
		}
		l := string(line) + "\n"
		in <- l
	}
	close(in)
	ioerr := fp.Close()
	if ioerr != nil {
		panic(ioerr)
	}
}

func timedText(format string, line string, n *int, opt map[string]bool) string {
	time2ref := map[string]*regexp.Regexp{
		"01":       regexp.MustCompile("%m"),
		"02":       regexp.MustCompile("%d"),
		"15":       regexp.MustCompile("%H"),
		"04":       regexp.MustCompile("%M"),
		"05":       regexp.MustCompile("%S"),
		"2006":     regexp.MustCompile("%Y"),
		"15:04:05": regexp.MustCompile("%T"),
		"Mon":      regexp.MustCompile("%W"),
		"Jan":      regexp.MustCompile("%[hb]"),
		"January":  regexp.MustCompile("%B"),
	}
	if opt["np"] {
		line = escapeString(line, opt)
	}
	if opt["ends"] {
		line = regexp.MustCompile("([\r\n])$").ReplaceAllString(line, "$$$1")
	}
	if opt["num"] {
		line = fmt.Sprintf("%6d\t", *n) + line
		*n++
	}
	for k, v := range time2ref {
		format = v.ReplaceAllString(format, k)
	}
	t := time.Now().Format(format)

	return t + line
}

func escapeString(str string, opt map[string]bool) string {
	// from source code of cat
	out := ""
	for i := 0; i < len(str); i++ {
		ch := str[i]
		if ch >= 32 {
			if ch < 127 {
				out += string(ch)
			} else if ch == 127 {
				out += "^?"
			} else {
				out += "M-"
				if ch >= 128+32 {
					if ch < 128+127 {
						out += string(rune(ch - 128))
					} else {
						out += "^?"
					}
				} else {
					out += "^" + string(rune(ch-128+64))
				}
			}
		} else if ch == '\t' && opt["tabs"] {
			out += "^I"
		} else if ch == '\n' {
			out += "\n"
		} else {
			out += "^" + string(rune(ch+64))
		}
	}
	return out
}

func parseOptions() (string, map[string]bool, string) {
	format := flag.String("f", "%Y-%m-%d %T", "time fomat.")
	delimiter := flag.String("d", ": ", "delimiter.")
	showTabs := flag.Bool("T", false, "display TAB characters as ^I")
	showNp := flag.Bool("v", false, "use ^ and M- notation, except for LFD and TAB")
	vT := flag.Bool("t", false, "equivalent to -vT")
	showEnds := flag.Bool("E", false, "display $ at end of each line")
	vE := flag.Bool("e", false, "equivalent to -vE")
	vET := flag.Bool("A", false, "equivalent to -vET")
	num := flag.Bool("n", false, "number all output lines")
	help := flag.Bool("h", false, "show this help")

	flag.Parse()

	opt := make(map[string]bool)

	if *vET {
		opt["np"] = true
		opt["ends"] = true
		opt["tabs"] = true
	} else {
		if *vE {
			opt["np"] = true
			opt["ends"] = true
		}
		if *vT {
			opt["np"] = true
			opt["tabs"] = true
		}
	}
	if *showTabs {
		opt["tabs"] = true
	}
	if *showEnds {
		opt["ends"] = true
	}
	if *showNp {
		opt["np"] = true
	}
	if *num {
		opt["num"] = true

	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	var file string
	args := flag.Args()
	if len(args) > 0 {
		file = args[0]
	}

	return *format + *delimiter, opt, file
}
