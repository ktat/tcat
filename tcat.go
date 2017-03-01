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

type tcat struct {
	formatStr string
	number    *int64
	opt       map[string]bool
}

func main() {
	number := int64(1)

	tcat := tcat{
		number: &number,
		opt:    make(map[string]bool),
	}

	files := tcat.parseOptions()

	if len(files) == 0 {
		in := make(chan string)
		go readStdin(in)
		for {
			l, ok := <-in
			if ok == false {
				break
			} else {
				fmt.Print(tcat.timedText(l))
			}
		}
	} else {
		for _, file := range files {
			fp, ioerr := os.Open(file)
			if ioerr != nil {
				panic("cannot read file")
			}
			tcat.readFileAndPrint(fp)
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

func (tcat *tcat) readFileAndPrint(fp *os.File) {
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
		fmt.Print(tcat.timedText(string(line) + "\n"))
	}
	ioerr := fp.Close()
	if ioerr != nil {
		panic(ioerr)
	}
}

func (tcat *tcat) timedText(line string) string {
	format := tcat.formatStr
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
	if tcat.opt["np"] {
		line = tcat.escapeString(line)
	}
	if tcat.opt["ends"] {
		line = regexp.MustCompile("([\r\n])$").ReplaceAllString(line, "$$$1")
	}
	if tcat.opt["num"] {
		line = fmt.Sprintf("%6d\t", *tcat.number) + line
		*tcat.number++
	}
	for k, v := range time2ref {
		format = v.ReplaceAllString(format, k)
	}
	t := time.Now().Format(format)

	return t + line
}

func (tcat *tcat) escapeString(str string) string {
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
		} else if ch == '\t' && tcat.opt["tabs"] {
			out += "^I"
		} else if ch == '\n' {
			out += "\n"
		} else {
			out += "^" + string(rune(ch+64))
		}
	}
	return out
}

func (tcat *tcat) parseOptions() []string {
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

	if *vET {
		tcat.opt["np"] = true
		tcat.opt["ends"] = true
		tcat.opt["tabs"] = true
	} else {
		if *vE {
			tcat.opt["np"] = true
			tcat.opt["ends"] = true
		}
		if *vT {
			tcat.opt["np"] = true
			tcat.opt["tabs"] = true
		}
	}
	if *showTabs {
		tcat.opt["tabs"] = true
	}
	if *showEnds {
		tcat.opt["ends"] = true
	}
	if *showNp {
		tcat.opt["np"] = true
	}
	if *num {
		tcat.opt["num"] = true
	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	tcat.formatStr = *format + *delimiter
	return flag.Args()
}
