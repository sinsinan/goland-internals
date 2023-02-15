package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	opFuncMap := map[string]func(float64, float64) float64{}

	populateOpFuncMap(&opFuncMap)

	for {
		fmt.Print("-> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("exit", strings.ToLower(text)) == 0 {
			break
		}
		v1, op, v2, err := parseExp(text)
		if err != nil {
			panic(err)
		}
		fmt.Printf("v1: %f, op: %s, v2: %f, out: %f\n", v1, op, v2, opFuncMap[op](v1, v2))
	}
}

func populateOpFuncMap(m *map[string]func(float64, float64) float64) {
	filepath.Walk("./plugins", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		if !info.IsDir() && strings.HasSuffix(path, ".so") {
			p, err := plugin.Open(path)
			if err != nil {
				panic(err)
			}
			key, err := p.Lookup("Operator")
			if err != nil {
				panic(err)
			}
			value, err := p.Lookup("Operate")
			if err != nil {
				panic(err)
			}
			(*m)[key.(func() string)()] = value.(func(float64, float64) float64)
		}
		return nil
	})
}

func parseExp(s string) (float64, string, float64, error) {
	expParser := regexp.MustCompile(`^([0-9]+) *([+-/*]) *([0-9]+) *$`)
	parts := expParser.FindStringSubmatch(s)
	if len(parts) != 4 {
		return 0, "", 0, errors.New(fmt.Sprintf("Invalid part length for expression: %s", s))
	}
	v1, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, "", 0, err
	}
	v2, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return 0, "", 0, err
	}
	return v1, parts[2], v2, nil
}
