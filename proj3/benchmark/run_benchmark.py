#!/bin/python3
#
# SBATCH --mail-user=shchen16@cs.uchicago.edu
# SBATCH --mail-type=ALL
# SBATCH --job-name=proj3_benchmark
# SBATCH --output=./slurm/out/%j.%N.stdout
# SBATCH --error=./slurm/out/%j.%N.stderr
# SBATCH --chdir=/home/shchen16/class_work/MPCS-52060_Parallel/nBody-golang/proj3/benchmark
# SBATCH --partition=debug
# SBATCH --nodes=1
# SBATCH --ntasks=1
# SBATCH --exclusive
# SBATCH --cpus-per-task=16
# SBATCH --mem-per-cpu=300

from posix import environ
import matplotlib
import matplotlib.pyplot as plt
import subprocess
import sys
import os


# Load module for test
subprocess.call("module load golang/1.16.2", shell=True)

# Global set up for test script
matplotlib.use('Agg')
save_path = "./slurm/out/"
seq_result = {}
threadCountEnv = "TEST_NBODY_THREAD_COUNT"
planetCountEnv = "TEST_NBODY_PLANET_COUNT"
filePathEnv = "TEST_NBODY_BENCHMARK_RESULT_PATH"
filePath = "./test_result.txt"
relativePath = "../main/"
planetCount_list = [500, 1000, 2000, 4000]
threadCount_list = [2, 4, 6, 8, 12]

# Run Parallel Test
def run_test(planetCount, threadCount):
    total_time = 0
    # Set argument for benchmark
    os.environ[threadCountEnv] = str(threadCount)
    os.environ[planetCountEnv] = str(planetCount)
    # Run five time and record total time
    for _ in range(5):
        subprocess.call(
            ['go', 'test', relativePath, "-run", "TestBenchmark"], 
            stdout=subprocess.PIPE
        )
        # Read result from benchmark
        f=open(relativePath+filePath, 'r')
        total_time += float(f.readline())
        f.close()

    # Clear argument for benchmark
    del os.environ[threadCountEnv]
    del os.environ[planetCountEnv]

    # Return results
    return total_time / 5

# Create env varible for filePath
os.environ[filePathEnv] = filePath


# Run seq tests
print("[Main] Run Seq Tests")
sys.stdout.flush()
for planetCount in planetCount_list:
    seq_result[planetCount] = run_test(planetCount,0)
    print("[Seq]", planetCount, "planets took", seq_result[planetCount], 
          "seconds to finish")
    sys.stdout.flush()

# Run parallel tests and create graph
print("[Main] Run Parallel Tests")
sys.stdout.flush()

# Setup environment
result = {}
# Create graph
plt.figure()
plt.xlabel('Number of threads')
plt.ylabel('Speed up')
# Check each block_size
for planetCount in planetCount_list:
    # Init varible
    result[planetCount] = {}
    # Run test for different threads count
    for threads in threadCount_list:
        result[planetCount][threads] = run_test(planetCount, threads)
        result[planetCount][threads] = seq_result[planetCount] / result[planetCount][threads]
        # Print Status
        print("[Parallel] ", planetCount, "thread_count=" +
              threads, "speedup", result[planetCount][threads])
        sys.stdout.flush()
    plt.plot(*zip(*sorted(result[planetCount].items())), label="Planet Count="+str(planetCount))
# Save graph
print("[Parallel]", planetCount, "finished, Creating graph")
sys.stdout.flush()
plt.legend()
plt.title("N-Body Speed Up Graph")
plt.savefig(save_path+"speedup.png")
# Print Status
print("[Parallel] Graph saved", "speedup.png")
sys.stdout.flush()

# Delete env varible for filePath
del os.environ[filePathEnv]

