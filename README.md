# Sudoku

A high-performance Sudoku library implemented in Go. This library provides comprehensive functionality for working with Sudoku puzzles, including brute-force and logical solving, puzzle generation, difficulty grading, and detailed step-by-step solutions with explanations.

The project is my first attempt at Go beyond a REST API course, created to deepen my understanding of the language through practical implementation of complex algorithms and data structures.

## Quick Start

```go
import (
    "github.com/DavidGudovic/sudoku/internal/core/board"
    "github.com/DavidGudovic/sudoku/internal/core/solver"
)

// Create a board from a string (0 represents empty cells)
puzzle, _ := board.FromString(
    "530070000" +
    "600195000" +
    "098000060" +
    "800060003" +
    "400803001" +
    "700020006" +
    "060000280" +
    "000419005" +
    "000080079",
)

// Note: Candidates can also be serialized in the string representation using * (e.g., "130*5*6*7..." means candidates 5, 6, and 7 are possible for the 3rd cell [0, 2])

// Solve using brute force (backtracking)
bruteForceSolver := solver.NewBruteForceSolver()
steps, _ := bruteForceSolver.Solve(puzzle)

// Solve using human-like logical techniques
logicalSolver := solver.NewLogicalSolver()
steps, _ := logicalSolver.Solve(puzzle)

// Take a single step for interactive solving
step, _ := logicalSolver.TakeAStep(puzzle)
```

## Features

### Board Representation

Efficient representation of Sudoku boards with automatic constraint management:

- Bitmask-based candidate tracking (uint16 per cell)
- Support for standard 9x9 Sudoku puzzles
- Automated constraint propagation (Can be turned off)
- Automatic detection of invalid/solved board states 
- Serialization and deserialization to/from strings

```go
// Create coordinates and set values
coords, _ := board.NewCoordinates(4, 5)
puzzle.SetValueOnCoords(coords, 7)

// Check board state
state := puzzle.State() // Invalid, Unsolved, or Solved

// Examine cell candidates
cell := puzzle.CellAt(coords)
candidates := cell.Candidates()
if candidates.Contains(7) {
    // 7 is a valid candidate for this cell
}
```

### Brute Force Solver

Efficient brute-force solver using optimized backtracking:

- **Minimum Remaining Values (MRV)** heuristic for smart variable selection
- **Constraint propagation** to reduce search space
- **Forward checking** to prevent invalid assignments
- Solves even the hardest puzzles in under 1ms on modern hardware

```go
solver := solver.NewBruteForceSolver()
steps, err := solver.Solve(puzzle)
if err != nil {
    // Puzzle is unsolvable
}
```

> **Note**: Knuth's Algorithm X with Dancing Links implementation is in progress.

### Logical Step Solver

Human-like solving techniques producing step-by-step solutions:

**Implemented Techniques:**
- Last Digit (Row, Column, Box variants)
- Naked Single
- Naked Pair
- Naked Triple
- Naked Quad
- Hidden Single
- Hidden Pair
- Hidden Triple
- Hidden Quad

**Planned Techniques:**
- Locked Candidates
- X-Wing
- Y-Wing
- W-Wing
- Swordfish
- Skyscraper
- Two-String Kite
- AIC's 

Each technique returns a detailed `Step` structure containing all information needed for educational UIs:

```go
type Step struct {
    Description       string              // Human-readable explanation
    Technique         string              // Technique name
    AffectedCells     []Coordinates       // Cells where changes occurred
    RemovedCandidates CandidateSet        // Candidates eliminated
    PlacedValue       *PlacedValue        // Value placed (if any)
    ReasonCells       []Coordinates       // Cells that justify this step
}
```

This makes the library ideal for building educational Sudoku applications that explain solving strategies step-by-step.

### PeerSet and PeerQuery API

Composable, type-safe query API for expressing complex solving techniques in near-natural language:

- Efficient querying of peers (rows, columns, boxes) for any cell
- Bitmask-based set operations ([9]uint16 representation)
- Immutable operations with method chaining
- O(1) containment checks and set operations
- Compiler enforced type safety to prevent invalid queries

```go

// Find all peers of a cell across row, column, and box
peers := Peers.Of(coords).Across(
    Row,
    Column,
    Box,
)

// Use variadic syntax for all scopes
allPeers := Peers.Of(coords).Across(AllScopes...)

// Find peers in shared scopes between multiple coordinates
sharedPeers := Peers.Of(coords1, coords2, coords3, ...).AcrossSharedScopes()

// Query all cells in a scope
columnPeers := Peers.InScope(Column, 5)

// Filter peers containing specific candidates
candidates, _ := board.NewCandidateSet(1, 2, 3)
candidatePeers := peers.ContainingCandidates(puzzle, candidates)

// Set operations for combining and filtering
union := peers1.Union(peers2)               // Combine two peer sets
intersection := peers1.Intersection(peers2) // Only peers in both sets
difference := peers1.Except(peers2)         // Peers in first but not second

// Incremental set building
expanded := peers.With(newCoord)          // Add a single coordinate
reduced := peers.Without(coord)           // Remove a single coordinate

expanded := peers.Including(newCoords...)       // Add multiple coordinates
reduced := peers.Excluding(coords...)           // Remove multiple coordinates

// Complex query combining multiple operations
filtered := Peers.Of(coords).
Across(Row).
NotContainingCandidates(puzzle, candidates).
Except(otherPeers).
Including(additionalCoord)
```

The PeerSet API enables techniques to be written as you think about the problem, without wrestling with nested loops and index calculations.

### Generator

> **Status**: Work in progress

### Grader

> **Status**: Work in progress

## Custom Solvers

Create custom solvers with specific technique combinations:

```go
import "github.com/DavidGudovic/sudoku/internal/core/techniques"

customSolver := solver.NewSudokuSolver([]techniques.Technique{
    techniques.LastDigit,
    techniques.NakedSingle,
    techniques.NakedPair,
    // Add your own techniques
})

steps, _ := customSolver.Solve(puzzle)
```

## Planned Features

### Web Interface
- Play Sudoku online with interactive UI
- Generate puzzles of varying difficulty
- Step-by-step solution viewer with detailed explanations
- Visual highlighting of affected cells and candidates
- Difficulty grading for custom puzzles

### Terminal Application
- SSH-based terminal UI using Charm libraries
- Interactive puzzle solving in the terminal
- Real-time hints and technique explanations

## Performance
- **Speed**: Solves even the hardest puzzles in ~ 1ms on modern hardware
- **Efficiency**: Bitmasks for O(1) candidate and peer operations
- **Memory**: Minimal overhead due to efficient data structures
- **Benchmarks**: Comprehensive benchmark suite included

```
cpu: AMD Ryzen 7 9800X3D 8-Core Processor
```

- **Backtracking alone**: 
```
BenchBacktracking/Easy-16         	           62132    19149 ns/op	 80 B/op   1 allocs/op
BenchBacktracking/Vicious-16      	            1720   697738 ns/op	 97 B/op   3 allocs/op
BenchBacktracking/Beyond_Hell-16  	             981  1205182 ns/op	 99 B/op   3 allocs/op
BenchBacktracking/AlEscargot_(2006)-16          2917   410264 ns/op	 80 B/op   1 allocs/op
```

- **Calculating board state**: 
```
BenchBoard_State/Known_unsolved_board-16      9422727   128.0 ns/op	  0 B/op   0 allocs/op
BenchBoard_State/Known_solved_board-16        5626734   211.1 ns/op	  0 B/op   0 allocs/op
```
## License

This project is open source and available under the MIT License.