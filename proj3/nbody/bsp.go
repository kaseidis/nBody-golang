package nbody

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"proj3/data"
	"runtime"
	"sync"
	"sync/atomic"
)

// =================== Data structures ===================
type bspContext struct {
	input          *json.Decoder
	output         *json.Encoder
	currentTask    *data.Task
	numOfThreads   int
	lockedThreadP1 int
	lockedThreadP2 int
	currentIter    int
	exitCounter    *int32
	done           bool
	cond           *sync.Cond
}

// ================ BSP General Function ================

// Run BSP parallel version of program
func RunBsp(input *json.Decoder, output *json.Encoder, numOfThreads int) {
	ctx := initBSPContext(input, output, numOfThreads) // Initialize BSP context
	for idx := 0; idx < numOfThreads-1; idx++ {
		go executeBSP(idx, ctx)
	}
	executeBSP(numOfThreads-1, ctx)
}

// Block and aggreate bsp threads
func bspBlocking(tid int, numOfThreads int,
	cond *sync.Cond, lockCount *int, aggregate func()) {
	// Blocking, wait all planets has location updated
	cond.L.Lock()
	if *lockCount < numOfThreads-1 {
		// If current thread isn't last thread to be locked
		*lockCount++
		cond.Wait()
		*lockCount--
	} else {
		// If current thread is last thread finished job,
		// Aggeraget result and wake up other threads
		aggregate()
		cond.Broadcast()
	}
	cond.L.Unlock()
}

// ==================== Util Function ====================

// Generate range of job that need to be processed
func generateRange(tid int, totalThreads int, totalJob int) (int, int) {
	slice := float64(totalJob) / float64(totalThreads)
	lower := slice * float64(tid)
	upper := slice * float64(tid+1)
	return int(math.Round(lower)), int(math.Round(upper))
}

// Init shared context
func initBSPContext(input *json.Decoder, output *json.Encoder, numOfThreads int) *bspContext {
	var exitCounter int32
	var task data.Task
	// Read first task
	if err := input.Decode(&task); err != nil {
		// Unexpected error, print error
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}
		// There is no task, end program directly
		os.Exit(0)
	}
	return &bspContext{
		input:          input,
		output:         output,
		currentTask:    &task,
		exitCounter:    &exitCounter,
		numOfThreads:   numOfThreads,
		lockedThreadP1: 0,
		lockedThreadP2: 0,
		currentIter:    1,
		done:           false,
		cond:           sync.NewCond(&sync.Mutex{}),
	}
}

// ============== nbody BSP Implementation ==============

// Output results and update shared status
func updateContext(ctx *bspContext) {
	// Construct result
	var result data.Result
	result.Planets = ctx.currentTask.Planets
	result.TimeStamp = float64(ctx.currentIter) * ctx.currentTask.Step
	// Encode json
	ctx.output.Encode(result)
	// Increase Iterater
	ctx.currentIter++
	// Check should we process next task
	if ctx.currentIter > ctx.currentTask.Iterations {
		ctx.currentIter = 1
		// Decode json
		if err := ctx.input.Decode(ctx.currentTask); err != nil {
			// Unexpected error, print error
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err)
			}
			ctx.done = true
		}
	}
}

// excuteBSP for BSP model, that calculate nbody problem
func executeBSP(tid int, ctx *bspContext) {
	totalThreads := ctx.numOfThreads
	for {
		// Caculate range this thread should process
		totalJob := len(ctx.currentTask.Planets)
		minRange, maxRange := generateRange(tid, totalThreads, totalJob)
		// Check if we need terminate
		if ctx.done {
			atomic.AddInt32(ctx.exitCounter, 1)
			// If it is main thread, wait until all thread finished
			if tid == ctx.numOfThreads-1 {
				for atomic.LoadInt32(ctx.exitCounter) != int32(ctx.numOfThreads) {
					runtime.Gosched()
				}
			}
			return
		}

		// Super Step 1 (Phase 1)
		{
			// Caclulate force for assigned planets
			for i := minRange; i < maxRange; i++ {
				for j, planetB := range ctx.currentTask.Planets {
					if i != j {
						ctx.currentTask.Planets[i].UpdateForce(&planetB,
							ctx.currentTask.G, ctx.currentTask.Softning)
					}
				}
			}
			// Blocking
			bspBlocking(tid, ctx.numOfThreads, ctx.cond,
				&ctx.lockedThreadP1, func() {})
		}

		// Super Step 2 (Phase 2)
		{
			// Update Planets location
			for i := minRange; i < maxRange; i++ {
				ctx.currentTask.Planets[i].UpdateSpeed(ctx.currentTask.Step)
				ctx.currentTask.Planets[i].UpdateLocation(ctx.currentTask.Step)
			}
			// Blocking
			bspBlocking(tid, ctx.numOfThreads, ctx.cond,
				&ctx.lockedThreadP2, func() {
					updateContext(ctx)
				})
		}
	}
}
