package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"proj3/data"
	"strconv"
	"testing"
	"time"
)

///////
// Auxiliary functions needed for the tests.
//////

// Genreate random float64 in range [min, max)
func randFloat64(min, max float64, rand *rand.Rand) float64 {
	return min + (max-min)*rand.Float64()
}

// Generate test data, output to *json.Encoder
func testDataGenerator(locationMin float64, locationMax float64,
	vMin float64, vMax float64, planetCount int) *data.Task {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	task := data.Task{
		Step:       0.25,
		G:          1.0,
		Iterations: 1000,
		Softning:   1e-8,
		Planets:    make([]data.Planet, 0),
	}
	for i := 0; i < planetCount; i++ {
		task.Planets = append(task.Planets, data.Planet{
			Mass: rand.Float64() + 0.5,
			Location: data.Vector3{
				X: randFloat64(locationMin, locationMax, r1),
				Y: randFloat64(locationMin, locationMax, r1),
				Z: randFloat64(locationMin, locationMax, r1),
			},
			Speed: data.Vector3{
				X: randFloat64(vMin, vMax, r1),
				Y: randFloat64(vMin, vMax, r1),
				Z: randFloat64(vMin, vMax, r1),
			},
		})
	}
	return &task
}

// Compare two results
func compareResults(a *data.Result, b *data.Result) bool {
	if a.TimeStamp != b.TimeStamp {
		return false
	}
	if len(a.Planets) != len(b.Planets) {
		return false
	}
	for i := range a.Planets {
		result := a.Planets[i].Mass == b.Planets[i].Mass
		result = result && (a.Planets[i].Location.X == b.Planets[i].Location.X)
		result = result && (a.Planets[i].Location.Y == b.Planets[i].Location.Y)
		result = result && (a.Planets[i].Location.Z == b.Planets[i].Location.Z)
		result = result && (a.Planets[i].Speed.X == b.Planets[i].Speed.X)
		result = result && (a.Planets[i].Speed.Y == b.Planets[i].Speed.Y)
		result = result && (a.Planets[i].Speed.Z == b.Planets[i].Speed.Z)
		if !result {
			return false
		}
	}
	return true
}

//////
// Beginning of actual nbody tests
//////

// Validate result by comparing result in seqential version and parallel version
func TestValidation(t *testing.T) {
	// Check if numOfThreads already set in env varible
	numOfThreads, present := os.LookupEnv("TEST_NBODY_THREAD_COUNT")
	if !present {
		numOfThreads = "2"
	} else if numOfThreads == "1" || numOfThreads == "0" {
		numOfThreads = "2"
	}
	// Check if planetCount already set in env varible
	planetCountStr, present := os.LookupEnv("TEST_NBODY_PLANET_COUNT")
	if !present {
		planetCountStr = "500"
	}
	// Check if planetCount is an valid number
	planetCount, errAtoi := strconv.Atoi(planetCountStr)
	if errAtoi != nil {
		t.Fatal("<TestValidation>: atoi error in getting planetCount, please check env varible TEST_NBODY_PLANET_COUNT")
	}

	// Generate Test task
	task := testDataGenerator(-3, 3, -1.5, 1.5, planetCount)

	// Reserve space for seq and par results
	resultSeq := make([]data.Result, 0)
	resultPar := make([]data.Result, 0)

	// =================== get result from seq version ===================
	{
		ctx, cancel1 := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel1()
		cmd := exec.CommandContext(ctx, "go", "run", "proj3/main", "0")

		stdin, errIn := cmd.StdinPipe()
		if errIn != nil {
			t.Fatal("<TestValidation>: stdin error in getting stdin pipe")
		}

		stdout, errOut := cmd.StdoutPipe()
		if errOut != nil {
			t.Fatal("<TestValidation>: stdout error in getting stdout pipe")
		}

		// Start Commandline
		if err := cmd.Start(); err != nil {
			t.Fatal("<TestValidation> cmd.Start error in executing test")
		}

		done := make(chan bool)

		var errData error

		go func() {
			defer func() { done <- true }()
			// Write data to stdin
			encoder := json.NewEncoder(stdin)
			errData = encoder.Encode(task)
			stdin.Close()
			if errData != nil {
				return
			}
			// Read data from stdout
			decoder := json.NewDecoder(stdout)
			for {
				var result data.Result
				errData := decoder.Decode(&result)
				if errData != nil {
					return
				}
				resultSeq = append(resultSeq, result)
			}
		}()

		<-done
		if errData != nil && errData != io.EOF {
			t.Fatal("<TestValidation> testDataGenerator error in stdpipe")
		}

		if err := cmd.Wait(); err != nil {
			t.Errorf("The automated test timed out. You may have a deadlock, starvation issue and/or you did not implement" +
				"the necessary code for passing this test.")
		}
	}
	//  =================== get result from parallel version ===================
	{
		ctx, cancel2 := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel2()
		cmd := exec.CommandContext(ctx, "go", "run", "proj3/main", numOfThreads)

		stdin, errIn := cmd.StdinPipe()
		if errIn != nil {
			t.Fatal("<TestValidation>: stdin error in getting stdin pipe")
		}

		stdout, errOut := cmd.StdoutPipe()
		if errOut != nil {
			t.Fatal("<TestValidation>: stdout error in getting stdout pipe")
		}

		// Start Commandline
		if err := cmd.Start(); err != nil {
			t.Fatal("<TestValidation> cmd.Start error in executing test")
		}

		done := make(chan bool)

		var errData error

		go func() {
			defer func() { done <- true }()
			// Write data to stdin
			encoder := json.NewEncoder(stdin)
			errData := encoder.Encode(task)
			stdin.Close()
			if errData != nil {
				return
			}
			// Read data from stdout
			decoder := json.NewDecoder(stdout)
			for {
				var result data.Result
				errData := decoder.Decode(&result)
				if errData != nil {
					return
				}
				resultPar = append(resultPar, result)
			}
		}()

		<-done
		if errData != nil && errData != io.EOF {
			t.Fatal("<TestValidation> testDataGenerator error in stdpipe")
		}
		if err := cmd.Wait(); err != nil {
			t.Errorf("The automated test timed out. You may have a deadlock, starvation issue and/or you did not implement" +
				"the necessary code for passing this test.")
		}
	}

	// Compare results
	if len(resultPar) != len(resultSeq) {
		t.Fatal("<TestValidation> two version have different length of output")
	}
	for i := range resultPar {
		if !compareResults(&resultPar[i], &resultSeq[i]) {
			t.Fatal("<TestValidation> two version have different results in output")
		}
	}
}

// Benchmark result, show time after finish genrate the task
// It will log the running time to test
// And it will save result to file if specified in ${TEST_NBODY_BENCHMARK_RESULT_PATH}
func TestBenchmark(t *testing.T) {
	// Check if numOfThreads already set in env varible
	numOfThreads, present := os.LookupEnv("TEST_NBODY_THREAD_COUNT")
	if !present {
		numOfThreads = "0"
	}
	// Check if planetCount already set in env varible
	planetCountStr, present := os.LookupEnv("TEST_NBODY_PLANET_COUNT")
	if !present {
		planetCountStr = "1000"
	}
	planetCount, errAtoi := strconv.Atoi(planetCountStr)
	if errAtoi != nil {
		t.Fatal("<TestBenchmark>: atoi error in getting planetCount, please check env varible TEST_NBODY_PLANET_COUNT")
	}

	// Genreate the task
	task := testDataGenerator(-3, 3, -1.5, 1.5, planetCount)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "proj3/main", numOfThreads)

	stdin, errIn := cmd.StdinPipe()
	if errIn != nil {
		t.Fatal("<TestBenchmark>: stdin error in getting stdin pipe")
	}

	if err := cmd.Start(); err != nil {
		t.Fatal("<TestBenchmark> cmd.Start error in executing test")
	}

	// Chanel for waiting io operation
	done := make(chan bool)

	var errData error
	var start time.Time

	go func() {
		encoder := json.NewEncoder(stdin)
		errData = encoder.Encode(task)
		start = time.Now()
		stdin.Close()
		done <- true
	}()

	<-done
	// Checking error
	if errData != nil {
		t.Fatal("<TestBenchmark> testDataGenerator error in stdpipe")
	}

	if err := cmd.Wait(); err != nil {
		t.Errorf("The automated test timed out. You may have a deadlock, starvation issue and/or you did not implement" +
			"the necessary code for passing this test.")
	} else {
		duration := time.Since(start).Seconds()
		t.Log("Result:", duration, "s")
		// Check if resultPath defined in env varible
		resultPath, present := os.LookupEnv("TEST_NBODY_BENCHMARK_RESULT_PATH")
		if present {
			file, err := os.OpenFile(resultPath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				t.Fatal("<TestBenchmark> Failed to save file")
			} else {
				fmt.Fprintln(file, duration)
				file.Close()
			}
		}
	}
}
