@echo off
:: Check if NBODY_THREAD is defined in environment
VERIFY OTHER 2>nul
SETLOCAL ENABLEEXTENSIONS
IF ERRORLEVEL 1 ECHO Unable to enable extensions
IF DEFINED NBODY_THREAD (SET THREAD_COUNT="%NBODY_THREAD%") ELSE (Set THREAD_COUNT="4")

:: Print number of threads
ECHO Running in %THREAD_COUNT% threads, you can set up environment varible 'NBODY_THREAD' to change this

:: Run program
ECHO Generating result from test1.json file
go run proj3/main %THREAD_COUNT% < test1.json > test1.result.json 

ECHO Generating result from test2.json file
go run proj3/main %THREAD_COUNT% < test2.json > test2.result.json 

ECHO Generating result from test3.json file
go run proj3/main %THREAD_COUNT% < test3.json > test3.result.json

:: Print usage
ECHO Results generated, please using 'result_visulize.htm' in this folder to visulize result