// ************************************************************
//
// sample.go
// Daniel Wiese
// 28 January, 2019
//
// ************************************************************

package sample

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"time"
)

type configDatumType struct {
	SolverConfig solverConfigType `json:"solverConfig"`
}

type solverConfigType struct {
	Epsilon       float64    `json:"epsilon"`
	MaxIterations int        `json:"maxIterations"`
	MaxResidual   float64    `json:"maxResidual"`
	Domain        domainType `json:"domain"`
}

type domainType struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// Set to enable print for debugging (swap `_ `and `out` to enable/disable)
var _ io.Writer = ioutil.Discard
var out = os.Stdout

func makeRange(min float64, max float64, n int) []float64 {
	if n <= 1 {
		return []float64{}
	}
	delta := (max - min) / float64(n-1)
	a := make([]float64, n)
	for i := range a {
		a[i] = min + float64(i)*delta
	}
	return a
}

func findMax(x []float64) float64 {
	var n float64
	for _, v := range x {
		if v > n {
			n = v
		}
	}
	return n
}

func solveBVP(xMin float64, xMax float64, eEpsilon float64, nCount int, eResidMax float64) string {
	const nInt int = 101

	deltaX := ((xMax - xMin) / float64(nInt-1))

	eResid := 1.0
	iCount := 0
	wOmega := 1.99

	var uOld [nInt]float64
	var uNew [nInt]float64
	var err [nInt - 1]float64
	var iResid []float64
	var a [nInt]float64
	var b [nInt]float64
	var c [nInt - 1]float64
	var d [nInt]float64

	start := time.Now()

	// Initial guess for u vector
	for i := 0; i < nInt; i++ {
		uOld[i] = 1.0
	}

	for eResid > eResidMax && iCount < nCount {
		for i := 1; i < nInt-1; i++ {
			a[i] = (-uOld[i] / (2 * deltaX)) - (eEpsilon / (math.Pow(deltaX, 2)))
			b[i] = ((2 * eEpsilon) / (math.Pow(deltaX, 2)))
			c[i] = (uOld[i] / (2 * deltaX)) - (eEpsilon / (math.Pow(deltaX, 2)))
			d[i] = 0
		}

		a[0] = 0
		b[0] = 1
		c[0] = 0
		d[0] = 1

		a[nInt-1] = 0
		b[nInt-1] = 1
		d[nInt-1] = -1

		// Start tridiagonal solver
		for i := 1; i < nInt; i++ {
			b[i] = b[i] - c[i-1]*a[i]/b[i-1]
			d[i] = d[i] - d[i-1]*a[i]/b[i-1]
		}

		uNew[nInt-1] = d[nInt-1] / b[nInt-1]

		for i := nInt - 2; i > -1; i-- {
			uNew[i] = (d[i] - c[i]*uNew[i+1]) / b[i]
		}
		// End tridiagonal solver

		iCount = iCount + 1

		err[0] = 0

		for i := 1; i < nInt-1; i++ {
			err[i] = uNew[i]*((uNew[i+1]-uNew[i-1])/(2*deltaX)) - eEpsilon*((uNew[i-1]-(2*uNew[i])+uNew[i+1])/(math.Pow(deltaX, 2)))
		}

		// uOld = uOld + wOmega * (uNew - uOld);
		var s [nInt]float64
		for i := range uOld {
			s[i] = uOld[i] + (uNew[i]-uOld[i])*wOmega
		}
		uOld = s

		var absErr [nInt - 1]float64

		for i := range absErr {
			absErr[i] = math.Abs(err[i])
		}

		eResid = findMax(absErr[:])
		iResid = append(iResid, eResid)
	}

	elapsed := time.Since(start)
	ms := float64(elapsed / time.Millisecond)

	return fmt.Sprintf("{\"result\":{\"iterations\":%d, \"residual\":%e, \"time\":%d}}", iCount, eResid, int(ms))

}

// ConcatenateStrings does the obvious
func ConcatenateStrings(string1 string, string2 string) string {
	fmt.Fprintf(out, "[ GO ] ConcatenateStrings called\n")

	return string1 + string2
}

// AddTwoNumbers does the obvious
func AddTwoNumbers(num1 float64, num2 float64) float64 {
	fmt.Fprintf(out, "[ GO ] AddTwoNumbers called\n")

	return num1 + num2
}

func float64FromByteSlice(byteSlice []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(byteSlice))
}

func float64SliceFromBase64String(str string) []float64 {
	// Decode base64 string to byte slice
	byteSlice, _ := b64.StdEncoding.DecodeString(str)

	var floatSlice []float64

	for i := 0; i < len(byteSlice)/8; i++ {
		floatSlice = append(floatSlice, float64FromByteSlice(byteSlice[8*i:8*(i+1)]))
	}

	return floatSlice
}

func base64StringFromFloat64Slice(arr []float64) string {
	// Convert to byte slice
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, arr)
	data := buf.Bytes()

	// Encode byte slice as base64 string
	return b64.StdEncoding.EncodeToString([]byte(data))
}

// IncrementSliceElements takes a base64 string representing an array of floats decodes it and increments each element of the array by 1
func IncrementSliceElements(str string) string {
	fmt.Fprintf(out, "[ GO ] IncrementSliceElements called\n")

	byteSlice, _ := b64.StdEncoding.DecodeString(str)

	var floatSlice []float64

	for i := 0; i < len(byteSlice)/8; i++ {
		floatSlice = append(floatSlice, float64FromByteSlice(byteSlice[8*i:8*(i+1)])+1)
	}

	return base64StringFromFloat64Slice(floatSlice)
}

// SolveBVPWithoutInputs solves BVP with fixed parameters
func SolveBVPWithoutInputs() string {
	const xMin float64 = -1.0
	const xMax float64 = 1.0

	const eEpsilon float64 = 0.1
	const nCount int = 250000
	const eResidMax float64 = 1.0e-11

	return solveBVP(xMin, xMax, eEpsilon, nCount, eResidMax)
}

// SolveBVPWithInputs requires user specified parameters to solve BVP
func SolveBVPWithInputs(configString string) string {
	// Find config settings
	var configDatum configDatumType
	json.Unmarshal([]byte(configString), &configDatum)

	fmt.Fprintf(out, "[ GO ] configDatum: %s\n", configDatum)

	// Get the stuff from config
	xMin := configDatum.SolverConfig.Domain.Min
	xMax := configDatum.SolverConfig.Domain.Max

	eEpsilon := configDatum.SolverConfig.Epsilon
	nCount := configDatum.SolverConfig.MaxIterations
	eResidMax := configDatum.SolverConfig.MaxResidual

	return solveBVP(xMin, xMax, eEpsilon, nCount, eResidMax)
}
