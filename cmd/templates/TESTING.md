# Testing

This doc describes required steps and tools to successfully implement testing strategy for all Privateer components.
Assumptions:
- VS Code is used as the local IDE
- Go built-in testing package is used

## Table of Contents

- [TLDR](#tldr)
- [Pre-requisites](#pre-requisites)
- [Environment Setup](#environment-setup)
- [How to write tests for packages](#how-to-write-tests-for-packages)
  * [Test files](#test-files)
  * [Test functions](#test-functions)
  * [Proposed test format](#proposed-test-format)
- [How to run tests locally](#how-to-run-tests-locally)
  * [Run all tests](#run-all-tests)
  * [Run all tests for a specific package](#run-all-tests-for-a-specific-package)
  * [Run a specific test for a specific function](#run-a-specific-test-for-a-specific-function)
  * [Generate test coverage report locally](#generate-test-coverage-report-locally)
- [How to debug tests](#how-to-debug-tests)
- [How to run tests in CICD pipeline](#how-to-run-tests-in-cicd-pipeline)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>


## TLDR

- Read go testing docs. [See here](#pre-requisites)
- Install *gotest* and *dlv* tools. [See here](#environment-setup)
- To write tests:
  - Open go file in VS Code, right-click on desired function and select *Go: Generate Unit Tests for Function*
    - This will automatically generate test file and test function boilerplate
  - Add code to test function
- To run tests
  - All tests within package:
  ```
  cd ./packagename/
  go test -v
  ```
  - Specific test function
  ```
  cd ./packagename/
  go test -v -run TestFunctionName
  ```
  - All tests within project
  ```
  cd ./projectroot/
  go test ./... -v
  ```
- To debug tests:
  - Open test function in VS Code and place breakpoint
  - Click *Debug test* on top of test function

## Pre-requisites

As a pre-requisite, please review the following materials to understand Go's testing appraoch and tools:

- https://golang.org/pkg/testing/ > Original reference to Go's built-in testing package
- https://gobyexample.com/testing > Clear examples about Go testing
- https://code.visualstudio.com/docs/languages/go > Details about Go extension for VS Code

## Environment Setup

A simple way to install useful tools in VS Code:

- To check which tools are installed:
  - Press *Ctrl+Shift+P* to open Command Palette
  - Type *Go: Locate Configured Go Tools*
  - See output report
- To install / update tools:
  - Press *Ctrl+Shift+P* to open Command Palette
  - Type *GO: Install/Update Tools*
  - Several checkboxes should appear; select desired  tools and click *Ok*
- Restart VS Code

Recommended tools from above list:

- *gotest* > Generate unit tests
- *dlv* > Enhanced Go debugging

## How to write tests for packages

### Test files

- For each package, a new file shall be created following this naming convention: *[packagename]_test.go*
- Example: 
  - Package: ```utils.go```
  - Test file: ```utils_test.go```

### Test functions

- For each function within a package, there shall be a corresponding test function following this naming convention: *func Test[FunctionName](t* **testing.T)*
- Example:
  - Package: ```utils.go```
  - Function: ```ReadStaticFile```
  - Test file: ```utils_test.go```
  - Test function: ```TestReadStaticFile(t *testing.T)```

### Proposed test format

Developers are free to write any code within the test function, however following a consistent structure can provide good reability, minimize communication issues among team members and ensure good coverage for all test cases.

The following structure is recommended:

- Use a *table-driven* approach: test inputs and expected outputs are listed in a table, and a single loop walks over them and performs the test logic. [See this reference.](https://gobyexample.com/testing)

- Use *t.Error* to report test failures and continue execution

- Use *t.Fatal* to stop test immediately, such as the case when test resources are not available or fail, which would make tests inconsistent

- Use *t.Skip(message)* to skip tests. Always log reason for skipping as the message.
  
  While is not realistic (nor desired) to achieve 100% test coverage, is is generally a best practice to avoid skipping tests. However, in some circumstances this is needed due to defered implementation or complex integration cases. Make sure the team in onboard with the decision of skipping a test.

- Use *subtests* when looping thru test cases. This will provide meaningful output when running tests in *verbose* mode
  - Example:
    - Instead of this
    ```
    func TestSum(t *testing.T) {
      tests := []struct {
        x int
        y int
        expected int
      }{
        {1, 1, 2},
        {1, 2, 3},
        {2, 2, 4},
        {5, 2, 7},
      }

      for _, test := range tests {
        total := Sum(test.x, test.y)
        if total != test.expected {
          t.Errorf("Sum of (%d+%d) was incorrect, got: %d, want: %d.", test.x, test.y, total, test.expected)
        }
      }
    }
    ```
    - Do this (notice the use of *t.Run*)
    ```
    func TestSum(t *testing.T) {
      tests := []struct {
        x        int
        y        int
        expected int
        testName string
      }{
        {1, 1, 2, "TestCase#1"},
        {1, 2, 3, "TestCase#2"},
        {2, 2, 4, "TestCase#3"},
        {5, 2, 7, "TestCase#4"},
      }

      for _, test := range tests {
        t.Run(test.testName, func(t *testing.T) {
          total := Sum(test.x, test.y)
          if total != test.expected {
            t.Errorf("Sum of (%d+%d) was incorrect, got: %d, want: %d.", test.x, test.y, total, test.expected)
          }
        })
      }
    }
    ```
  - Sample output:
    ```
    ❯ go test -v -run TestSum
    === RUN   TestSum
    === RUN   TestSum/TestCase#1
    === RUN   TestSum/TestCase#2
    === RUN   TestSum/TestCase#3
    === RUN   TestSum/TestCase#4
    --- PASS: TestSum (0.02s)
        --- PASS: TestSum/TestCase#1 (0.00s)
        --- PASS: TestSum/TestCase#2 (0.00s)
        --- PASS: TestSum/TestCase#3 (0.00s)
        --- PASS: TestSum/TestCase#4 (0.00s)
    PASS
    ok      github.com/privateerproj/privateer-sdk/utils 0.227s
    ```

- Autogenerate tests with *gotest tool*

  To be able to leverage the automated code generation, make sure you have *gotest tool* installed and integrated into VS Code. [See above details for installing tools](#environment-setup).

  This approach can not only save significant time, but it also assists with generating required code for complex structures.

  Steps to generate test code:
  - Open go file in VS Code, right-click on desired function and select *Go: Generate Unit Tests for Function*
  - This will automatically generate test file (if it doesn't exist) and test function boilerplate

  Understanding test boilerplate (see inline comments):
  ```
  func TestSum(t *testing.T) {
    type args struct { // This struct represents the arguments to be passed to the function under test
      x int
      y int
    }
    tests := []struct {
      name string // Name for test case
      args args   // Arguments for test case
      want int    // Expected value
    }{
      // TODO: Add test cases.
      {"TestCase1_AddingTwoPositiveNumbers_ShouldReturnSum", args{1, 1}, 2},
		  {"TestCase2_AddingZeroToNumber_ShouldReturnNumber", args{0, 3}, 3},
    }
    for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
        if got := Sum(tt.args.x, tt.args.y); got != tt.want {
          t.Errorf("Sum() = %v, want %v", got, tt.want)
        }
      })
    }
  }
  ```
  Sample output:
  ```
  ❯ go test -v
  === RUN   TestSum
  === RUN   TestSum/TestCase1_AddingTwoPositiveNumbers_ShouldReturnSum
  === RUN   TestSum/TestCase2_AddingZeroToNumber_ShouldReturnNumber
  --- PASS: TestSum (0.02s)
      --- PASS: TestSum/TestCase1_AddingTwoPositiveNumbers_ShouldReturnSum (0.00s)
      --- PASS: TestSum/TestCase2_AddingZeroToNumber_ShouldReturnNumber (0.00s)
  PASS
  ok      github.com/privateerproj/privateer-sdk/utils 0.223s
  ```


## How to run tests locally

### Run all tests

- Navigate to project's root folder and execute ```go tests ./...```
```
cd ./projectroot/
go test -v ./...
```

### Run all tests for a specific package

- Navigate to package folder and execute ```go test```

Sample
```
cd ./internal/utils
go test -v
```

### Run a specific test for a specific function

- Navigate to package folder and execute ```go test run <FunctionName>```

Sample:
```
cd ./internal/utils
go test -v -run TestReadStaticFile
```

### Generate test coverage report locally

- Generate coverage profile ```go test ./... -coverprofile coverage.out```
- Generate HTML report ```go tool cover -html coverage.out```

*Note: Please notice the above commands will generate a local file *coverage.out* with the test coverage details. Make sure you exclude this file from your next commit as it is not needed.*

## How to debug tests

During development it is very convenient to debug code execution, step into functions and check values of local variables and stack trace. This is possible in VS Code thru the use of *dlv* go package. [See above details for installing tools](#environment-setup).

To debug tests:
  - Open test function in VS Code and place breakpoint
  - Click *Debug test* on top of test function

See [this reference](https://code.visualstudio.com/docs/languages/go#_debugging) for more details on VS Code debugging.

## How to run tests in CICD pipeline

Since we are using Github Actions, the following commands shall be added to *.github/workflows/ci.yml*
```
sudo go test ./... -coverprofile coverage.out -covermode count
sudo go tool cover -func coverage.out
```
In addition to executing all tests, we are displaying the test coverage report for every function as well as total coverage.