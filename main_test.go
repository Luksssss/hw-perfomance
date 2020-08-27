package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
)

func BenchmarkCount(b *testing.B) {
	file, err := os.Open("logs")
	if err != nil {
		panic(err.Error())
	}

	logs, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	parts := divide(logs, runtime.GOMAXPROCS(0))

	allStr := make([]string, 0, len(parts))

	b.ReportAllocs()
	b.SetBytes(int64(len(logs)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		ch := make(chan string, len(parts))

		for j := 0; j < len(parts); j++ {
			go func(index int) {
				joinNumbers(parts[index], ch)
			}(j)
		}

		for s := range ch {
			allStr = append(allStr, s)

			if len(allStr) == len(parts) {
				close(ch)
			}
		}

		// итоговая строка
		strings.Join(allStr, ", ")
		allStr = allStr[:0]

	}

}
