# Codeforces CLI
This is a command-line tool for testing and submitting solutions to programming problems on Codeforces.

![demo](demo.svg)

## Features:
* Supports time and memory limits
* Supports submitting solutions
* Tests cases are run in parallel, taking advantage of multiple threads
## Installation
1. Install [Go](https://golang.org/dl/) and [GCC](https://gcc.gnu.org/)
1. Clone the repository: `git clone https://github.com/markojukic/codeforces_cli`
1. CD into the repository: `cd codeforces_cli`
1. To build, run: `make`
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
## Usage
To see information about flags, run:
```
./codeforces_cli --help
```
