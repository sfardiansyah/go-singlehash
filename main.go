package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"tester/bloom"
	"time"
)

const (
	w = 200000
	// w = 200
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	m := 80000
	// m := 3
	s := make([]string, m)
	n := 40000
	// n := 2
	var tmp int

	f, err := os.Open("data/uniqueelement.txt")
	check(err)
	defer f.Close()

	fout, err := os.Create("data/dat1.dat")
	check(err)
	defer fout.Close()

	fout2, err := os.Create("data/dat2.dat")
	check(err)
	defer fout2.Close()

	r := bufio.NewScanner(f)

	for i := 0; i < m && r.Scan(); i++ {
		s[i] = r.Text()
	}

	for k := 2; k < 21; k++ {
		sh := bloom.NewBF(w, k, 0)
		lh := bloom.NewBF(w, k, 1)

		c1p := 0
		c2p := 0

		for i := 0; i < n; i++ {
			sh.Insert([]byte(s[i]), uint(len(s[i])))
			lh.Insert([]byte(s[i]), uint(len(s[i])))
		}

		for i := n; i < m; i++ {
			tmp = sh.Query([]byte(s[i]), uint(len(s[i])))
			if tmp == 0 {
				c1p++
			}

			tmp = lh.Query([]byte(s[i]), uint(len(s[i])))
			if tmp == 0 {
				c2p++
			}
		}

		fout.WriteString(strconv.Itoa(k) + " " + fmt.Sprintf("%f", 1-(float64(c1p)/float64(m-n))) + " " + fmt.Sprintf("%f", 1-(float64(c2p)/float64(m-n))) + "\n")

		start := time.Now()
		for i := 0; i < m; i++ {
			for j := 0; j < 100; j++ {
				tmp = sh.Query([]byte(s[i]), uint(len(s[i])))
			}
		}
		finish := time.Now()
		dur := finish.Sub(start)
		d1 := dur.Nanoseconds

		start = time.Now()
		for i := 0; i < m; i++ {
			for j := 0; j < 100; j++ {
				tmp = lh.Query([]byte(s[i]), uint(len(s[i])))
			}
		}
		finish = time.Now()
		dur = finish.Sub(start)
		d2 := dur.Nanoseconds

		fout2.WriteString(strconv.Itoa(k) + " " + fmt.Sprintf("%f", (float64(m)/float64(d1()/1000000))) + " " + fmt.Sprintf("%f", (float64(m)/float64(d2()/1000000))) + "\n")
	}
}
