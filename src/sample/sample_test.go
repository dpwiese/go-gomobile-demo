package sample

import (
	"encoding/json"
	"reflect"
	"testing"
)

type resultType struct {
	Result resultBodyType `json:"result"`
}

type resultBodyType struct {
	Iterations  int     `json:"iterations"`
	Residual    float64 `json:"residual"`
	ElapsedTime float64 `json:"time"`
}

func TestConcatenateStrings(t *testing.T) {
	type testCases struct {
		string1  string
		string2  string
		expected string
	}

	var tests = []testCases{
		{string1: "hello", string2: "world", expected: "helloworld"},
		{string1: "foo", string2: "bar", expected: "foobar"},
	}

	for ic, tc := range tests {
		var result = ConcatenateStrings(tc.string1, tc.string2)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("[Test #%d] Expected is: %s, got: %s", ic+1, tc.expected, result)
		}
	}
}

func TestAddTwoNumbers(t *testing.T) {
	type testCases struct {
		num1     float64
		num2     float64
		expected float64
	}

	var tests = []testCases{
		{num1: 1, num2: 2, expected: 3},
		{num1: 1, num2: 0, expected: 1},
	}

	for ic, tc := range tests {
		var result = AddTwoNumbers(tc.num1, tc.num2)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("[Test #%d] Expected is: %s, got: %s", ic+1, tc.expected, result)
		}
	}
}

func TestIncrementSliceElements(t *testing.T) {
	type testCases struct {
		str      string
		expected string
	}

	var tests = []testCases{
		{str: "mpmZmZmZ8T+amZmZmZkBQA==", expected: "zczMzMzMAECamZmZmZkJQA=="},
	}

	for ic, tc := range tests {
		var result = IncrementSliceElements(tc.str)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("[Test #%d] Expected is: %s, got: %s", ic+1, tc.expected, result)
		}
	}
}

func TestMakeRange(t *testing.T) {
	type rangeArray struct {
		min float64
		max float64
		n   int
	}

	type testCases struct {
		test     rangeArray
		expected []float64
	}

	var tests = []testCases{
		{test: rangeArray{min: 1, max: 3, n: 3}, expected: []float64{1, 2, 3}},
		{test: rangeArray{min: 1, max: 3, n: 2}, expected: []float64{1, 3}},
		{test: rangeArray{min: 1, max: 3, n: 5}, expected: []float64{1, 1.5, 2, 2.5, 3}},
		{test: rangeArray{min: 1, max: 3, n: 1}, expected: []float64{}},
		{test: rangeArray{min: 1, max: 3, n: 0}, expected: []float64{}},
	}

	for ic, tc := range tests {
		var expected = makeRange(tc.test.min, tc.test.max, tc.test.n)
		if !reflect.DeepEqual(expected, tc.expected) {
			t.Errorf("[Test #%d] Expected is: %f, got: %f", ic+1, tc.expected, expected)
		}
	}
}

func TestFindMax(t *testing.T) {
	type testCases struct {
		test     []float64
		expected float64
	}

	var tests = []testCases{
		{test: nil, expected: 0},
		{test: []float64{-1, 2, 3, 4}, expected: 4},
		{test: []float64{5, 1, 2, -3}, expected: 5},
	}

	for ic, tc := range tests {
		var expected = findMax(tc.test)
		if expected != tc.expected {
			t.Errorf("[Test #%d] Expected is: %f, got: %f", ic+1, tc.expected, expected)
		}
	}
}

func TestSolveBVPWithoutInputs(t *testing.T) {
	type testCases struct {
		expected int
	}

	var tests = []testCases{
		{expected: 35680},
	}

	for ic, tc := range tests {
		var expected = SolveBVPWithoutInputs()

		var resultOutput resultType
		json.Unmarshal([]byte(expected), &resultOutput)

		iters := resultOutput.Result.Iterations

		if iters != tc.expected {
			t.Errorf("[Test #%d] Expected is: %d, got: %d", ic+1, tc.expected, iters)
		}
	}
}

func TestSolveBVPWithInputs(t *testing.T) {
	type testCases struct {
		configString string
		expected     int
	}

	var tests = []testCases{
		{configString: "{\"solverConfig\":{\"epsilon\": 0.1, \"maxIterations\": 250000, \"maxResidual\": 1.0e-11, \"domain\": {\"min\": -1.0, \"max\": 1.0 }}}", expected: 35680},
	}

	for ic, tc := range tests {
		var expected = SolveBVPWithInputs(tc.configString)

		var resultOutput resultType
		json.Unmarshal([]byte(expected), &resultOutput)

		iters := resultOutput.Result.Iterations
		if iters != tc.expected {
			t.Errorf("[Test #%d] Expected is: %f, got: %f", ic+1, tc.expected, iters)
		}
	}
}
