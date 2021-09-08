package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func main() {
	r := csv.NewReader(os.Stdin)
	r.Comma = ';'
	all, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"Date", "Payee", "Memo", "Outflow", "Inflow"})
	for _, line := range all[1:] {
		outflow := line[7]
		inflow := line[6]
		date, err := time.Parse("02.01.06", line[3])
		if err != nil {
			log.Fatal(err)
		}
		ss := strings.Fields(line[5])
		payee := "etc"
		for _, s := range ss {
			if len(s) == 7 && strings.HasPrefix(s, "MCC") {
				payee = s
				break
			}
		}
		memo := ""
		if len(ss) > 2 {
			mm := strings.FieldsFunc(ss[1], func(c rune) bool {
				return !unicode.IsLetter(c) && !unicode.IsNumber(c)
			})
			if len(mm) > 3 {
				memo = mm[3]
			}
		}
		if memo == "" {
			rurI := findRUR(ss)
			if rurI > 0 {
				memo = strings.Join(ss[2:rurI-3], "")
			}
		}
		if memo == "" {
			memo = line[5]
		}

		w.Write([]string{date.Format("01/02/2006"), payee, memo, outflow, inflow})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func findRUR(s []string) int {
	for i, v := range s {
		if v == "RUR" {
			return i
		}
	}
	return 0
}
