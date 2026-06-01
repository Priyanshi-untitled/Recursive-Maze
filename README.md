# 🌀 Recursive Maze — Code Olympics 2026

**By Priyanshi Varshney** | Language: Go | Line Limit: 400 | Domain: Game

---

## What is this?

A fully playable terminal maze game written in **Go**, under 400 lines, built for Code Olympics 2026.

The maze is randomly generated every round. Three classic pathfinding algorithms — **BFS, DFS, and A\*** — solve it simultaneously using recursion. Then *you* try to beat them.

---

## Constraints Met

| Constraint | Value |
|---|---|
| Language | Go |
| Line Limit | ≤ 400 lines |
| Domain | Game |
| Restriction | Recursive algorithms only |

---

## Features

- 🗺️ **Random maze generation** — different every run using recursive backtracking
- 🤖 **3 AI solvers run in parallel** — BFS, DFS, A* (all recursive, concurrent goroutines)
- 🎮 **Playable** — move with WASD or arrow keys, 60-second timer
- 🌫️ **Fog of War mode** — only see nearby cells as you explore
- 🔥 **Heatmap view** — see which cells the algorithms visited most
- 📊 **Animated replay** — watch all three algorithms race step by step
- 🏆 **Scoring** — based on steps taken vs optimal path + time remaining

---

## How to Run

**Requirements:** Go 1.21+

```bash
git clone https://github.com/YOUR_USERNAME/recursive-maze
cd recursive-maze
go mod tidy
go run main.go
```

> ⚠️ Run in a real terminal (CMD on Windows, Terminal on Mac/Linux). VS Code's integrated terminal may not support raw keyboard input.

---

## Controls

| Key | Action |
|---|---|
| W / ↑ | Move Up |
| S / ↓ | Move Down |
| A / ← | Move Left |
| D / → | Move Right |
| Q | Quit |

---

## Menu Options

1. **Watch BFS / DFS / A\*** — animated race between all three algorithms
2. **Play (no fog)** — full maze visible, race to the exit
3. **Play (fog of war)** — limited visibility, harder challenge
4. **Quit**

---

## Scoring

`Score = 1000 - (your_steps - optimal_steps) × 5 + time_remaining`

- Optimal path = BFS solution length
- Score above 800 = ⚡ ELITE rating

---

## Algorithm Comparison

| Algorithm | Path | Speed | Strategy |
|---|---|---|---|
| BFS | Shortest | Fast | Level-by-level |
| DFS | Longer | Fast | Go deep first |
| A* | Shortest | Fastest | Heuristic guided |

---

## Project Structure

```
recursive-maze/
├── main.go       # entire game — 396 lines
├── go.mod
└── README.md
```

---

## Demo

▶️ [YouTube Demo](https://YOUR_YOUTUBE_LINK)

---

*Code Olympics 2026 submission — Priyanshi Varshney*
