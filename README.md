# 🌀 Recursive Maze — Code Olympics 2026

**By Priyanshi Varshney**
**Language:** Go
**Domain:** Game
**Constraint:** Recursive Algorithms Only
**Code Size:** 396 Lines

> *"One maze. Three recursive AI solvers. One player trying to beat them."*

---

## Overview

Recursive Maze is a fully playable terminal-based maze game built for **Code Olympics 2026** under a unique combination of constraints:

* Go as the assigned language
* Maximum 400 lines of code
* Game domain
* Recursive algorithms only

Each run generates a brand-new maze using recursive backtracking. Three pathfinding algorithms—**BFS, DFS, and A***—solve the maze simultaneously while visualizing their exploration strategies. After watching the algorithms race, the player can enter the maze and attempt to reach the exit before time runs out.

The project combines algorithm visualization, game design, recursion, and concurrency into a single interactive experience.

---

## Features

### 🗺️ Random Maze Generation

* Every maze is generated using recursive backtracking.
* Each run creates a completely different layout.
* No prebuilt maps or stored levels.

### 🤖 Parallel Algorithm Race

Three recursive pathfinding algorithms solve the same maze simultaneously:

* Breadth First Search (BFS)
* Depth First Search (DFS)
* A* Search

All three execute concurrently and can be replayed visually.

### 🎮 Playable Maze Game

* WASD controls
* Arrow key support
* 60-second timer
* Live movement tracking
* Score calculation
* Best score tracking

### 🌫️ Fog of War Mode

A second gameplay mode where only nearby cells are visible, forcing exploration and memory-based navigation.

### 🔥 Exploration Heatmap

Visualizes how heavily different regions of the maze were explored by the algorithms.

### 📊 Animated Replay

Watch BFS, DFS, and A* race through the maze step-by-step and compare their behavior visually.

---

## Constraint Compliance

| Constraint                | Requirement | Implementation                              |
| ------------------------- | ----------- | ------------------------------------------- |
| Language                  | Go          | Implemented entirely in Go                  |
| Domain                    | Game        | Fully playable terminal game                |
| Line Limit                | ≤ 400 Lines | Final implementation: 396 lines             |
| Recursive Algorithms Only | Required    | Recursive maze generation, BFS, DFS, and A* |

---

## Innovation Highlights

Most pathfinding projects only display the final solution.

Recursive Maze exposes the entire decision-making process.

Instead of hiding computation, the project turns algorithms into gameplay:

* Three AI solvers race in real time
* Exploration patterns generate a heatmap
* Players directly compare themselves against optimal solutions
* Fog-of-war transforms pathfinding into exploration

The result is both a game and an educational visualization tool.

---

## Constraint-Driven Engineering Decisions

The final architecture emerged directly from the imposed constraints.

### Recursive BFS

BFS is normally implemented using loops and an explicit queue.

To satisfy the recursive-only requirement, the queue itself became recursive state:

```go
bfsR(queue, previous, destination)
```

Each recursive call processes the next frontier.

This produced an architecture that is dramatically different from traditional BFS implementations.

### Working Within 400 Lines

The line budget forced:

* Shared rendering functions
* Reusable helpers
* Minimal duplication
* Compact state representation

Every function had to justify its existence.

### Concurrency as a Core Feature

Rather than simulating parallel execution, the project embraces Go's strengths.

BFS, DFS, and A* run simultaneously using goroutines and synchronize through channels, creating a real algorithm race.

---

## Cross-Constraint Architecture

One of the most interesting outcomes came from the interaction between two constraints:

1. Recursive algorithms only
2. Maximum 400 lines

Normally BFS is iterative.

Converting BFS into a recursive algorithm required the queue to become part of the recursive state itself.

Because DFS and A* were also recursive, all three solvers naturally evolved into a similar structural pattern.

This produced an unexpected benefit:

* BFS, DFS, and A* became visually comparable
* Their differences became easier to understand
* The implementations became educational as well as functional

The architecture exists specifically because these constraints collided.

Without those limitations, this design would never have emerged.

---

## Rosetta Stone: Same Problem, Different Language

The project was implemented in two different languages:

### Go Version

* `main.go`

### Python Version

* `maze_python.py`

Both versions solve the same problem while reflecting different language philosophies.

| Aspect               | Go                    | Python                |
| -------------------- | --------------------- | --------------------- |
| Concurrency          | Goroutines + Channels | Threads               |
| Typing               | Static                | Dynamic               |
| Execution            | Compiled              | Interpreted           |
| State Representation | Structs               | Tuples & Dictionaries |
| Design Style         | Explicit              | Flexible              |

The purpose of the Python implementation is not syntax translation.

It demonstrates how the same algorithmic problem naturally produces different architectures when expressed through different programming languages.

---

## What Go Taught Me

I started this project knowing Go syntax.

I finished it understanding Go's philosophy.

### Goroutines Feel Effortless

Running BFS, DFS, and A* simultaneously required only a few lines of code.

Go makes concurrency feel like a natural language feature rather than an advanced topic.

### The Compiler Helps More Than It Restricts

Unused variables, incorrect types, and missing returns are caught immediately.

What initially felt strict eventually became reassuring.

### Explicit Code Ages Better

Go encourages clarity over cleverness.

Reading the project later felt surprisingly easy because the language naturally pushes developers toward maintainable code.

### Biggest Lesson

Go prioritizes simplicity, readability, and reliability.

For a project centered around algorithms and visualization, those tradeoffs proved extremely valuable.

---

## Controls

| Key   | Action     |
| ----- | ---------- |
| W / ↑ | Move Up    |
| S / ↓ | Move Down  |
| A / ← | Move Left  |
| D / → | Move Right |
| Q     | Quit       |

---

## Scoring System

```text
Score = 1000 - (PlayerSteps - OptimalSteps) × 5 + TimeRemaining
```

Where:

* Optimal path = BFS path length
* Higher score = better performance
* Scores above 800 earn an Elite rating

---

## Algorithm Comparison

| Algorithm | Path Quality   | Strategy                           |
| --------- | -------------- | ---------------------------------- |
| BFS       | Optimal        | Level-by-level exploration         |
| DFS       | Usually longer | Deep exploration with backtracking |
| A*        | Optimal        | Heuristic-guided search            |

---

## Code Quality & Verification

The project was formatted and validated using Go's standard tooling.

```bash
go fmt
go vet
```

The codebase maintains a clean structure and follows idiomatic Go practices.

---

## Project Structure

```text
Recursive-Maze/
├── main.go
├── maze_python.py
├── go.mod
├── README.md
```

---

## Demo

The demo showcases:

1. Recursive maze generation
2. BFS vs DFS vs A* race
3. Heatmap visualization
4. Standard gameplay mode
5. Fog-of-war mode
6. Final score calculation

**Demo Video:**
Add YouTube link here

---

## Final Reflection

This project started as a maze game.

It became an exploration of how constraints influence software design.

The most valuable lesson was that limitations do not simply restrict solutions—they create new ones. Recursive Maze exists because of its constraints, not despite them.

---

**Code Olympics 2026**
**Priyanshi Varshney**
