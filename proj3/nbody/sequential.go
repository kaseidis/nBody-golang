package nbody

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"proj3/data"
)

// Run seqential version of program
func RunSeqential(input *json.Decoder, output *json.Encoder) {
	var task data.Task
	for {
		// Decode json
		if err := input.Decode(&task); err != nil {
			// Unexpected error, print error
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
		// Run task
		runTaskSeqential(task, output)
	}
}

// Run single task
func runTaskSeqential(task data.Task, output *json.Encoder) {
	var result data.Result
	for iter := 1; iter <= task.Iterations; iter++ {
		// Calculates force
		for i := range task.Planets {
			for j, planetB := range task.Planets {
				if i != j {
					task.Planets[i].UpdateForce(&planetB, task.G, task.Softning)
				}
			}
		}
		// Update Planets location
		for i := range task.Planets {
			task.Planets[i].UpdateSpeed(task.Step)
			task.Planets[i].UpdateLocation(task.Step)
		}
		// Output Result for current iteration
		result.Planets = task.Planets
		result.TimeStamp = float64(iter) * task.Step
		output.Encode(result)
	}
}
