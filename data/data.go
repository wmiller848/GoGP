package data

import (
	"bytes"
	"math"
	"strconv"
)

type TestData struct {
	Input     []float64
	Assert    float64
	AssertStr string
}

func New(buffer []byte, inputs int) ([]*TestData, float64, map[string]float64) {
	////////////////////////
	// TODO :: Fix this shit
	////////////////////////
	lines := bytes.Split(buffer, []byte("\n"))
	////////////////////////
	testData := []*TestData{}
	asserts := make(map[string]float64)
	for i, _ := range lines {
		if len(lines[i]) > 0 {
			vals := bytes.Split(lines[i], []byte(","))
			if len(vals) >= inputs {
				data := &TestData{}
				for j, val := range vals {
					num, err := strconv.ParseFloat(string(val), 64)
					if err == nil {
						if j < inputs {
							data.Input = append(data.Input, num)
						} else {
							data.Assert = num
							data.AssertStr = string(val)
							asserts[data.AssertStr] = data.Assert
						}
					} else {
						if string(val) == "true" {
							num = 1
						} else if string(val) == "false" {
							num = -1
						} else {
							num = NumberFromString(string(val))
						}
						if j < inputs {
							data.Input = append(data.Input, num)
						} else {
							data.Assert = num
							data.AssertStr = string(val)
							asserts[data.AssertStr] = data.Assert
						}
					}
				}
				testData = append(testData, data)
			}
		}
	}
	threshold := math.MaxFloat64
	for i, iv := range asserts {
		for j, jv := range asserts {
			if i != j {
				diff := math.Abs(iv - jv)
				if diff < threshold {
					threshold = diff
				}
			}
		}
	}
	return testData, threshold / 2.0, memo
}

var base = 8.0
var memo = make(map[string]float64)

func NumberFromString(str string) float64 {
	if memo[str] == 0 {
		memo[str] = float64(len(memo)+1) * base
	}
	return memo[str]
}
