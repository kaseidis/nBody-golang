# MPCS 52060 Project 3

## How to run


## Describe

- Describe in detailed of your system and the problem it is trying to solve.

- A description of how you implemented your parallel solution.

- Describe the challenges you faced while implementing the system. What aspects of the system might make it difficult to parallelize? In other words, what to you hope to learn by doing this assignment?


# How to test

- Specifications of the testing machine you ran your experiments on (i.e. Core Architecture (Intel/AMD), Number of cores, operating system, memory amount, etc.)

- What are the hotspots and bottlenecks in your sequential program? Were you able to parallelize the hotspots and/or remove the bottlenecks in the parallel version?

- What limited your speedup? Is it a lack of parallelism? (dependencies) Communication or synchronization overhead? As you try and answer these questions, we strongly prefer that you provide data and measurements to support your conclusions.

## Write up
- Run experiments with data you generate for both the sequential and parallel versions. As with the data provided by prior assignments, the data should vary the granularity of your parallel system. For the parallel version, make sure you are running your experiments with at least producing work for N threads, where N = {2,4,6,8,12}. You can go lower/larger than those numbers based on the machine you are running your system on. You are not required to run project 3 on the Peanut cluster. You can run it on your local machine and base your N threads on the number of logical cores you have on your local machine. If you choose to run your system on your local machine then please state that in your report and the your machine specifications as well.

- Produce speedup graph(s) for those data sets.