# Wordle Solver

An entropy-based [wordle](https://www.powerlanguage.co.uk/wordle) solver.

## Install

```bash
# Clone the repository
git clone https://github.com/Jack97/wordle-solver.git /your/install/path

# Navigate to the repository
cd /your/install/path

# Build the binary
go build -o ./output/wordle-solver
```

## Usage

Run the binary to start an interactive solver session. For example:

```
$ ./output/wordle-solver

Guess #1: soare
Enter feedback (Grey = 0, Yellow = 1, Green = 2):

00000

Guess #2: clint
Enter feedback (Grey = 0, Yellow = 1, Green = 2):

01102

Guess #3: limit
Enter feedback (Grey = 0, Yellow = 1, Green = 2):

22002

Guess #4: light
Enter feedback (Grey = 0, Yellow = 1, Green = 2):

22222

Completed the wordle in 4/6 guesses
```

## Todo

* Optimise
* Benchmark against all possible answers
