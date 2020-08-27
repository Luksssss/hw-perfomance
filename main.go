package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/valyala/fastjson"
	insaneJSON "github.com/vitkovskii/insane-json"
)

func main() {
	// f, err := ioutil.ReadFile("logs")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// joinNumbers(f)
	// fmt.Println("OK")
}

func joinNumbers(logs []byte, ch chan string) {

	resultList := make([]string, 0, 64)

	lastI := 0

	for {
		endLine := bytes.IndexByte(logs[lastI:], '\n')

		if endLine < 0 {
			break
		}

		endLine += lastI

		msg := getMsgFastJSON(logs[lastI:endLine])

		lastI = endLine + 1
		// msg := getMsgInsaneJSON(logs[lastI:i])
		if msg == "" {
			continue
		}
		joinMessage(msg, &resultList)

	}

	ch <- strings.Join(resultList, ", ")
}

func joinMessage(m string, resultList *[]string) {
	lastI := -1
	for i, c := range m {
		if isNumber(c) {
			continue
		}
		if i-lastI > 1 {
			*resultList = append(*resultList, m[lastI+1:i])
		}
		lastI = i
	}
}

func isNumber(c rune) bool {

	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func divide(logs []byte, parts int) [][]byte {
	pos := 0
	step := len(logs) / parts
	result := make([][]byte, 0)

	for pos+step < len(logs) {
		start := pos
		finish := pos + step
		for logs[finish] != '\n' {
			finish++
		}
		finish++
		result = append(result, logs[start:finish])
		pos = finish
	}

	if pos != len(logs) {
		result = append(result, logs[pos:])
	}

	return result
}

func getMsgFastJSON(line []byte) string {
	return fastjson.GetString(line, "message")
}

func getMsgInsaneJSON(line []byte) string {
	root, err := insaneJSON.DecodeBytes(line)

	if err != nil {
		log.Println(line, err)
		return ""
	}
	msg := root.Dig("message").AsString()
	defer insaneJSON.Release(root)

	if msg == "" {
		return ""
	}

	return msg
}
