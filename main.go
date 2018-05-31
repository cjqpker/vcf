package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type person struct {
	name  string
	phone []string
}

func main() {
	if len(os.Args) < 2 {
		panic("usage: application path/to/file ")
	}
	fi, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	persons := make([]person, 0, 10240)
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		line := string(a)

		switch {
		case strings.HasPrefix(line, "BEGIN:"):
			persons = append(persons, person{})
		case strings.HasPrefix(line, "FN:"):
			persons[len(persons)-1].name = line[3:]
		case strings.Contains(line, "VOICE") || strings.Contains(line, "TEL;"):
			if strings.Contains(line, "TEL;") {
				newP := line[strings.LastIndex(line, ":")+1:]
				persons[len(persons)-1].phone = append(persons[len(persons)-1].phone, trimPhone(newP))
			}
		case strings.HasPrefix(line, "END:"):
			fmt.Println("processing:", persons[len(persons)-1].name, strings.Join(persons[len(persons)-1].phone, "@"))
			break
		}
	}

	fmt.Println("saving...")
	w := ""
	for _, v := range persons {
		w += v.name
		w += ","
		w += strings.Join(v.phone, ",")
		w += "\r\n"
	}
	if err := ioutil.WriteFile("./a.csv", []byte(w), 777); err != nil {
		panic(err)
	}
	fmt.Println("saved")
}

func trimPhone(p string) (s string) {
	for _, v := range p {
		if v == ' ' || v == '-' {

		} else {
			s += string(v)
		}
	}
	if strings.HasPrefix(s, "+86") {
		s = s[3:]
	}
	return
}
