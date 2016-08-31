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

type DataMap map[string]float64

func New(buffer []byte, inputs int) ([]*TestData, float64, DataMap, DataMap) {
	////////////////////////
	// TODO :: Fix this shit
	////////////////////////
	lines := bytes.Split(buffer, []byte("\n"))
	////////////////////////
	testData := []*TestData{}
	inputMap := make(DataMap)
	assertMap := make(DataMap)
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
						}
					} else {
						num = NumberFromString(string(val), j)
						if j < inputs {
							data.Input = append(data.Input, num)
							inputMap[string(val)] = num
						} else {
							data.Assert = num
							data.AssertStr = string(val)
							assertMap[data.AssertStr] = data.Assert
						}
					}
				}
				testData = append(testData, data)
			}
		}
	}

	threshold := math.MaxFloat64
	for i, iv := range assertMap {
		for j, jv := range assertMap {
			if i != j {
				diff := math.Abs(iv - jv)
				if diff < threshold {
					threshold = diff
				}
			}
		}
	}
	return testData, threshold / 2.0, inputMap, assertMap
}

var base = 8.0
var memo = make(map[int]DataMap)

func NumberFromString(str string, j int) float64 {
	if memo[j] == nil {
		memo[j] = make(DataMap)
	}
	if memo[j][str] == 0 {
		memo[j][str] = float64(len(memo[j])+1) * base
	}
	return memo[j][str]
}
