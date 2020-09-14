# Codeforces CLI
This is a command-line tool for testing and submitting solutions to programming problems on Codeforces.
## Features:
* Supports time and memory limits
* Supports submitting solutions
* Tests cases are run in parallel, taking advantage of multiple threads
## Installation
1. Download and install [Go](https://golang.org/dl/) and [GCC](https://gcc.gnu.org/)
1. Clone the repository: `git clone https://github.com/markojukic/codeforces_cli`
1. CD into the repository: `cd codeforces_cli`
1. Build: `make`
1. Add `codeforces_cli` to `$PATH` variable
## Configuration
### Example `config.json`
```json
{
    "codeforces": {
        "username": "Your codeforces username",
        "password": "Your codeforces password"
    }
}
```
## Example usage
```
$ ./codeforces_cli -site codeforces -contest 1 -problem A -file examples/1A.cpp
Problem 1 A
Time limit: 	1s
Memory limit: 	256MB
   Test case|   Verdict|      Time|
           1|        OK|   2.506ms|
```