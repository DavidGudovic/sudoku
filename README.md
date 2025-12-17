# Sudoku

This repository contains a Sudoku library implemented in GO. It provides various functionalities to work with Sudoku
puzzles, including brute force or logical solving, generating, grading puzzles, detailed readable solving steps with affected cells, removed candidates, placed values and reason cells, reasoning.

The project is my first attempt at Go, other than a REST Api in Golang course, and I created it to deepen my understanding
of the language.

## Features

### Board representation
-- Efficient representation of Sudoku boards with bitmasks for candidates.
-- Support for standard 9x9 Sudoku puzzles.
-- Automated constraint propagation.
-- Automated detection of invalid/solved boards.
-- Serialization and deserialization of boards to/from strings.

### Bruto force solver (Optimized Backtracking, Knuth's Algorithm X with Dancing Links)
-- Fast and efficient brute-force solver using optimized backtracking techniques.
--- Minimal remaining values (MRV) heuristic for variable selection.
--- Constraint propagation to reduce search space.
--- Forward checking to prevent invalid assignments.
-- Knuth's Algorithm X with Dancing Links implementation for exact cover problems.
--- Work in progress

### Logical Step Solver (Human like techniques producing a step-by-step solution)
-- Implements various human-like solving techniques, including:
--- Last Digit
--- Naked Pairs
--- Naked Triples ...
--- Hidden Singles
--- Hidden Pairs
--- Hidden Triples ...
--- Pointing Pairs
--- X-Wings
--- Skyscrapers
--- Two-String Kites
--- More features to be added soon.

### PeerSet and PeerQuery API
-- Efficient querying of peers (rows, columns, boxes) for a given cell.
-- PeerSet structure to manage and manipulate sets of using a [9]uint16 bitmask representation of the board.
-- Immutable, composable, type safe Query API for expressing complex techniques in an almost natural language.
-- # Type the code as you think about the problem #

### Generator - Work in progress
### Grader - Work in progress

## Outside Core Functionality but planned

### Web Interface
-- Play Sudoku online
-- Generate puzzles
-- Solve puzzles
-- Step through a solution with detailed explanations and highlighting
-- Grade puzzles

### Terminal app server over SSH using the Charm libraries
-- Play Sudoku in the terminal over an SSH connection

## Performance
-- The solver is optimized for performance and can solve even the hardest puzzles in <1 ms on modern hardware.
-- Uses bitmasks to represent candidates for constraint propagation and general board state.
-- Uses bitmasks to represent the entire coordinate system of the board and query it into PeerSet's for O(1) access to peers.
-- Very low memory overhead due to efficient data structures.
-- Benchmark tests included to measure and improve performance.