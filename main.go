package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

const ROWS, COLS, FOG = 21, 41, 4

var (
	maze               [ROWS][COLS]byte
	heat               [ROWS][COLS]int
	seen               [ROWS][COLS]bool
	rng                = rand.New(rand.NewSource(time.Now().UnixNano()))
	bp, dp, ap         []pt
	bv, dv, av         [ROWS][COLS]bool
	bt, dt, at         time.Duration
	steps, score, best int
	won, tout          bool
)

type pt struct{ r, c int }
type nd struct {
	p    pt
	g, h int
	up   *nd
}

var d4 = []pt{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ok(r, c int) bool { return r >= 0 && r < ROWS && c >= 0 && c < COLS }
func heur(a, b pt) int { return abs(a.r-b.r) + abs(a.c-b.c) }
func rep(s string, n int) string {
	o := ""
	for i := 0; i < n; i++ {
		o += s
	}
	return o
}

func carve(r, c int) {
	for _, i := range rng.Perm(4) {
		d := []pt{{0, 2}, {0, -2}, {2, 0}, {-2, 0}}[i]
		nr, nc := r+d.r, c+d.c
		if ok(nr, nc) && maze[nr][nc] == '#' {
			maze[r+d.r/2][c+d.c/2] = ' '
			maze[nr][nc] = ' '
			carve(nr, nc)
		}
	}
}

func gen() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			maze[r][c] = '#'
		}
	}
	maze[1][1] = ' '
	carve(1, 1)
	maze[1][1] = 'S'
	maze[ROWS-2][COLS-2] = 'E'
}

func bfsR(q []pt, prev map[pt]pt, end pt) []pt {
	if len(q) == 0 {
		return nil
	}
	cur, rest := q[0], q[1:]
	if cur == end {
		var path []pt
		for c := end; c != (pt{1, 1}); c = prev[c] {
			path = append([]pt{c}, path...)
		}
		return path
	}
	for _, d := range d4 {
		nb := pt{cur.r + d.r, cur.c + d.c}
		if ok(nb.r, nb.c) && !bv[nb.r][nb.c] && maze[nb.r][nb.c] != '#' {
			bv[nb.r][nb.c] = true
			heat[nb.r][nb.c]++
			prev[nb] = cur
			rest = append(rest, nb)
		}
	}
	return bfsR(rest, prev, end)
}

func dfsR(r, c int, end pt, path []pt) []pt {
	if !ok(r, c) || dv[r][c] || maze[r][c] == '#' {
		return nil
	}
	dv[r][c] = true
	heat[r][c] += 2
	path = append(path, pt{r, c})
	if (pt{r, c}) == end {
		return path
	}
	for _, d := range []pt{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		if res := dfsR(r+d.r, c+d.c, end, path); res != nil {
			return res
		}
	}
	return nil
}

func astarR(open []*nd, closed map[pt]bool, end pt) []pt {
	if len(open) == 0 {
		return nil
	}
	bi := 0
	for i, n := range open {
		if n.g+n.h < open[bi].g+open[bi].h {
			bi = i
		}
	}
	cur := open[bi]
	open = append(open[:bi], open[bi+1:]...)
	if cur.p == end {
		var path []pt
		for n := cur; n != nil; n = n.up {
			path = append([]pt{n.p}, path...)
		}
		return path[1:]
	}
	closed[cur.p] = true
	av[cur.p.r][cur.p.c] = true
	heat[cur.p.r][cur.p.c]++
	for _, d := range d4 {
		nb := pt{cur.p.r + d.r, cur.p.c + d.c}
		if ok(nb.r, nb.c) && !closed[nb] && maze[nb.r][nb.c] != '#' {
			open = append(open, &nd{p: nb, g: cur.g + 1, h: heur(nb, end), up: cur})
		}
	}
	return astarR(open, closed, end)
}

func solveAll() {
	done := make(chan int, 3)
	end := pt{ROWS - 2, COLS - 2}
	go func() { t := time.Now(); bv[1][1] = true; bp = bfsR([]pt{{1, 1}}, map[pt]pt{}, end); bt = time.Since(t); done <- 1 }()
	go func() { t := time.Now(); if r := dfsR(1, 1, end, nil); r != nil { dp = r[1:] }; dt = time.Since(t); done <- 1 }()
	go func() { t := time.Now(); ap = astarR([]*nd{{p: pt{1, 1}, h: heur(pt{1, 1}, end)}}, map[pt]bool{}, end); at = time.Since(t); done <- 1 }()
	<-done; <-done; <-done
}

func fogReveal(pos pt) {
	for r := pos.r - FOG; r <= pos.r+FOG; r++ {
		for c := pos.c - FOG; c <= pos.c+FOG; c++ {
			if ok(r, c) && abs(r-pos.r)+abs(c-pos.c) <= FOG {
				seen[r][c] = true
			}
		}
	}
}

func frameAlgo(b, d, a int) string {
	sets := [3]map[pt]bool{{}, {}, {}}
	for i := 0; i < b && i < len(bp); i++ { sets[0][bp[i]] = true }
	for i := 0; i < d && i < len(dp); i++ { sets[1][dp[i]] = true }
	for i := 0; i < a && i < len(ap); i++ { sets[2][ap[i]] = true }
	buf := "\033[H\033[1;37m╔" + rep("═", COLS) + "╗\033[0m\n"
	for r := 0; r < ROWS; r++ {
		buf += "\033[1;37m║\033[0m"
		for c := 0; c < COLS; c++ {
			p := pt{r, c}
			switch maze[r][c] {
			case '#':
				buf += "\033[90m█\033[0m"
			case 'S':
				buf += "\033[1;32mS\033[0m"
			case 'E':
				buf += "\033[1;31mE\033[0m"
			default:
				n := 0
				for _, s := range sets {
					if s[p] { n++ }
				}
				switch {
				case n == 3:     buf += "\033[1;35m◈\033[0m"
				case sets[0][p]: buf += "\033[34m·\033[0m"
				case sets[1][p]: buf += "\033[33m·\033[0m"
				case sets[2][p]: buf += "\033[36m·\033[0m"
				default:         buf += " "
				}
			}
		}
		buf += "\033[1;37m║\033[0m\n"
	}
	return buf + "\033[1;37m╚" + rep("═", COLS) + "╝\033[0m\n"
}

func framePlay(pos pt, fogOn bool) string {
	buf := "\033[H\033[1;37m╔" + rep("═", COLS) + "╗\033[0m\n"
	for r := 0; r < ROWS; r++ {
		buf += "\033[1;37m║\033[0m"
		for c := 0; c < COLS; c++ {
			switch {
			case fogOn && !seen[r][c]: buf += "\033[30m░\033[0m"
			case pt{r, c} == pos:     buf += "\033[1;92m@\033[0m"
			case maze[r][c] == '#':   buf += "\033[90m█\033[0m"
			case maze[r][c] == 'S':   buf += "\033[1;32mS\033[0m"
			case maze[r][c] == 'E':   buf += "\033[1;31mE\033[0m"
			default:                  buf += " "
			}
		}
		buf += "\033[1;37m║\033[0m\n"
	}
	return buf + "\033[1;37m╚" + rep("═", COLS) + "╝\033[0m\n"
}

func animate() {
	mx := len(bp)
	for _, p := range [][]pt{dp, ap} {
		if len(p) > mx { mx = len(p) }
	}
	fmt.Print("\033[2J")
	for s := 0; s <= mx; s++ {
		b, d, a := s, s, s
		if b > len(bp) { b = len(bp) }
		if d > len(dp) { d = len(dp) }
		if a > len(ap) { a = len(ap) }
		fmt.Print(frameAlgo(b, d, a))
		fmt.Printf("\033[34mBFS\033[0m:%d \033[33mDFS\033[0m:%d \033[36mA*\033[0m:%d  %d/%d\n", b, d, a, s, mx)
		time.Sleep(25 * time.Millisecond)
	}
}

func heatmap() {
	fmt.Print("\033[2J\033[H")
	mx := 1
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			if heat[r][c] > mx { mx = heat[r][c] }
		}
	}
	fmt.Println("\033[1;37m╔" + rep("═", COLS) + "╗\033[0m")
	for r := 0; r < ROWS; r++ {
		fmt.Print("\033[1;37m║\033[0m")
		for c := 0; c < COLS; c++ {
			if maze[r][c] == '#' { fmt.Print("\033[90m█\033[0m"); continue }
			switch heat[r][c] * 5 / mx {
			case 0:  fmt.Print("\033[96m░\033[0m")
			case 1:  fmt.Print("\033[34m▒\033[0m")
			case 2:  fmt.Print("\033[32m▓\033[0m")
			case 3:  fmt.Print("\033[33m▓\033[0m")
			default: fmt.Print("\033[31m█\033[0m")
			}
		}
		fmt.Println("\033[1;37m║\033[0m")
	}
	fmt.Println("\033[1;37m╚" + rep("═", COLS) + "╝\033[0m")
	fmt.Println("\033[96m░\033[0m=cold \033[34m▒\033[0m=light \033[32m▓\033[0m=mid \033[33m▓\033[0m=heavy \033[31m█\033[0m=hot")
}

func key() byte {
	fd := int(os.Stdin.Fd())
	old, _ := term.MakeRaw(fd)
	defer term.Restore(fd, old)
	b := make([]byte, 3)
	n, _ := os.Stdin.Read(b)
	if n == 3 && b[0] == 27 && b[1] == 91 { return b[2] }
	if n == 1 { return b[0] }
	return 0
}

func play(fogOn bool) {
	fmt.Print("\033[2J\033[?25l")
	defer fmt.Print("\033[?25h")
	pos := pt{1, 1}
	steps, won, tout = 0, false, false
	end := pt{ROWS - 2, COLS - 2}
	t0 := time.Now()
	if fogOn { fogReveal(pos) }
	for {
		tl := 60 - int(time.Since(t0).Seconds())
		if tl <= 0 { tout = true; break }
		col := 32
		if tl <= 10 { col = 31 }
		ftag := ""
		if fogOn { ftag = " [FOG]" }
		fmt.Print(framePlay(pos, fogOn))
		fmt.Printf("  Steps:%-4d  Time:\033[%dm%ds\033[0m%s  WASD=move  Q=quit\n", steps, col, tl, ftag)
		np := pos
		k := key()
		if k == 'w' || k == 'W' || k == 65 { np = pt{pos.r - 1, pos.c} }
		if k == 's' || k == 'S' || k == 66 { np = pt{pos.r + 1, pos.c} }
		if k == 'd' || k == 'D' || k == 67 { np = pt{pos.r, pos.c + 1} }
		if k == 'a' || k == 'A' || k == 68 { np = pt{pos.r, pos.c - 1} }
		if k == 'q' || k == 'Q' || k == 3  { return }
		if ok(np.r, np.c) && maze[np.r][np.c] != '#' { steps++; pos = np; if fogOn { fogReveal(pos) } }
		if pos == end {
			won = true
			v := 1000 - (steps-len(bp))*5 + tl
			if v < 50 { v = 50 }
			score = v
			if score > best { best = score }
			fmt.Print("\033[2J\033[H")
			fmt.Printf("\n  \033[1;92m★ YOU WIN!\033[0m  Steps:%-4d  Score:\033[1;93m%d\033[0m/1000  Optimal:%d\n", steps, score, len(bp))
			if score > 800 { fmt.Println("  \033[1;92m⚡ ELITE!\033[0m") }
			fmt.Print("\n  [Enter]...")
			fmt.Scanln()
			return
		}
	}
	fmt.Print("\033[2J\033[H\n  \033[1;31m⏱ TIME UP!\033[0m\n\n  [Enter]...")
	fmt.Scanln()
}

func results() {
	fmt.Print("\033[2J\033[H")
	fmt.Println("\033[1;37m╔────── ALGORITHM RESULTS ──────╗\033[0m")
	fmt.Printf("\033[34m  BFS\033[0m  path=%-4d  %v\n", len(bp), bt)
	fmt.Printf("\033[33m  DFS\033[0m  path=%-4d  %v\n", len(dp), dt)
	fmt.Printf("\033[36m  A* \033[0m  path=%-4d  %v\n", len(ap), at)
	if won { fmt.Printf("\033[1;92m  YOU\033[0m  steps=%-4d  score=\033[1;93m%d\033[0m/1000\n", steps, score) }
	if best > 0 { fmt.Printf("  \033[1;93m★ Best: %d\033[0m\n", best) }
	fmt.Println("\033[1;37m╚───────────────────────────────╝\033[0m")
	fmt.Print("\n  [R]eplay  [Enter] back: ")
	var s string
	fmt.Scanln(&s)
	if s == "r" || s == "R" { animate(); fmt.Print("\n  [Enter]..."); fmt.Scanln() }
}

func reset() {
	bp, dp, ap = nil, nil, nil
	won, tout = false, false
	steps, score = 0, 0
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			heat[r][c] = 0; bv[r][c] = false; dv[r][c] = false; av[r][c] = false; seen[r][c] = false
		}
	}
}

func menu() int {
	fmt.Print("\033[2J\033[H")
	fmt.Println("\033[1;37m╔══════════════════════════════════╗")
	fmt.Println("║   RECURSIVE MAZE — 4D 2026       ║")
	if best > 0 {
		fmt.Printf("║  \033[1;93m★ Best: %-6d\033[1;37m                ║\n", best)
	} else {
		fmt.Println("║                                  ║")
	}
	fmt.Println("╠══════════════════════════════════╣")
	fmt.Println("║  1.  Watch BFS / DFS / A*        ║")
	fmt.Println("║  2.  Play — no fog               ║")
	fmt.Println("║  3.  Play — fog of war           ║")
	fmt.Println("║  4.  Quit                        ║")
	fmt.Println("╚══════════════════════════════════╝\033[0m")
	fmt.Print("\n  Choice: ")
	var c int
	fmt.Scan(&c)
	fmt.Scanln()
	return c
}

func main() {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	for {
		reset()
		switch menu() {
		case 1:
			gen(); solveAll()
			fmt.Print("  [Enter] race..."); fmt.Scanln()
			animate()
			fmt.Print("\n  [Enter] heatmap..."); fmt.Scanln()
			heatmap()
			fmt.Print("\n  [Enter] results..."); fmt.Scanln()
			results()
		case 2:
			gen(); solveAll()
			fmt.Print("  [Enter] to play..."); fmt.Scanln()
			play(false); results()
		case 3:
			gen(); solveAll()
			fmt.Print("  [Enter] FOG mode..."); fmt.Scanln()
			play(true); results()
		case 4:
			fmt.Print("\033[2J\033[H"); return
		}
	}
}