"""
Recursive Maze — Python (Rosetta Stone)
Code Olympics 2026 | Priyanshi Varshney
Same algorithm, different language architecture.
"""
import sys, os, time, random, threading
sys.setrecursionlimit(10000)

ROWS, COLS, FOG = 21, 41, 4

maze  = [['#'] * COLS for _ in range(ROWS)]
heat  = [[0]   * COLS for _ in range(ROWS)]
seen  = [[False] * COLS for _ in range(ROWS)]

bp, dp, ap = [], [], []
bt, dt, at = 0, 0, 0
steps = score = best = 0
won = tout = False

D4 = [(-1,0),(1,0),(0,-1),(0,1)]

def ok(r, c): return 0 <= r < ROWS and 0 <= c < COLS
def heur(a, b): return abs(a[0]-b[0]) + abs(a[1]-b[1])
def rep(s, n): return s * n

# ── Maze generation ──────────────────────────────────────────────────────────

def carve(r, c):
    dirs = [(0,2),(0,-2),(2,0),(-2,0)]
    random.shuffle(dirs)
    for dr, dc in dirs:
        nr, nc = r+dr, c+dc
        if ok(nr, nc) and maze[nr][nc] == '#':
            maze[r+dr//2][c+dc//2] = ' '
            maze[nr][nc] = ' '
            carve(nr, nc)          # recurse

def gen():
    for r in range(ROWS):
        for c in range(COLS):
            maze[r][c] = '#'
    maze[1][1] = ' '
    carve(1, 1)
    maze[1][1] = 'S'
    maze[ROWS-2][COLS-2] = 'E'

# ── BFS (recursive) ──────────────────────────────────────────────────────────

bv = [[False]*COLS for _ in range(ROWS)]

def bfs_r(q, prev, end):
    if not q:
        return None
    cur, rest = q[0], q[1:]
    if cur == end:
        path, c = [], end
        while c != (1,1):
            path.insert(0, c)
            c = prev[c]
        return path
    for dr, dc in D4:
        nb = (cur[0]+dr, cur[1]+dc)
        r, c = nb
        if ok(r,c) and not bv[r][c] and maze[r][c] != '#':
            bv[r][c] = True
            heat[r][c] += 1
            prev[nb] = cur
            rest.append(nb)
    return bfs_r(rest, prev, end)

# ── DFS (recursive) ──────────────────────────────────────────────────────────

dv = [[False]*COLS for _ in range(ROWS)]

def dfs_r(r, c, end, path):
    if not ok(r,c) or dv[r][c] or maze[r][c] == '#':
        return None
    dv[r][c] = True
    heat[r][c] += 2
    path = path + [(r,c)]
    if (r,c) == end:
        return path
    for dr, dc in [(0,1),(1,0),(0,-1),(-1,0)]:
        res = dfs_r(r+dr, c+dc, end, path)
        if res:
            return res
    return None

# ── A* (recursive) ────────────────────────────────────────────────────────────

av = [[False]*COLS for _ in range(ROWS)]

def astar_r(open_list, closed, end):
    if not open_list:
        return None
    # find best node
    bi = min(range(len(open_list)), key=lambda i: open_list[i][0]+open_list[i][1])
    g, h, pos, parent = open_list[bi]
    rest = open_list[:bi] + open_list[bi+1:]
    if pos == end:
        path, node = [], (g, h, pos, parent)
        while node:
            path.insert(0, node[2])
            node = node[3]
        return path[1:]
    closed.add(pos)
    av[pos[0]][pos[1]] = True
    heat[pos[0]][pos[1]] += 1
    for dr, dc in D4:
        nb = (pos[0]+dr, pos[1]+dc)
        r, c = nb
        if ok(r,c) and nb not in closed and maze[r][c] != '#':
            rest.append((g+1, heur(nb,end), nb, (g,h,pos,parent)))
    return astar_r(rest, closed, end)

# ── Solve all (parallel threads) ─────────────────────────────────────────────

def solve_all():
    global bp, dp, ap, bt, dt, at, bv, dv, av
    bv = [[False]*COLS for _ in range(ROWS)]
    dv = [[False]*COLS for _ in range(ROWS)]
    av = [[False]*COLS for _ in range(ROWS)]
    end = (ROWS-2, COLS-2)

    def run_bfs():
        global bp, bt
        t = time.time()
        bv[1][1] = True
        bp = bfs_r([(1,1)], {}, end) or []
        bt = time.time() - t

    def run_dfs():
        global dp, dt
        t = time.time()
        r = dfs_r(1, 1, end, [])
        dp = r[1:] if r else []
        dt = time.time() - t

    def run_astar():
        global ap, at
        t = time.time()
        ap = astar_r([(0, heur((1,1),end), (1,1), None)], set(), end) or []
        at = time.time() - t

    threads = [threading.Thread(target=f) for f in (run_bfs, run_dfs, run_astar)]
    for t in threads: t.start()
    for t in threads: t.join()

# ── Fog ───────────────────────────────────────────────────────────────────────

def fog_reveal(pos):
    pr, pc = pos
    for r in range(pr-FOG, pr+FOG+1):
        for c in range(pc-FOG, pc+FOG+1):
            if ok(r,c) and abs(r-pr)+abs(c-pc) <= FOG:
                seen[r][c] = True

# ── Rendering ────────────────────────────────────────────────────────────────

def frame_algo(b, d, a):
    sb = set(bp[:b]); sd = set(dp[:d]); sa = set(ap[:a])
    buf = "\033[H\033[1;37m╔" + rep("═",COLS) + "╗\033[0m\n"
    for r in range(ROWS):
        buf += "\033[1;37m║\033[0m"
        for c in range(COLS):
            p = (r,c)
            ch = maze[r][c]
            if ch == '#':   buf += "\033[90m█\033[0m"
            elif ch == 'S': buf += "\033[1;32mS\033[0m"
            elif ch == 'E': buf += "\033[1;31mE\033[0m"
            else:
                n = (p in sb) + (p in sd) + (p in sa)
                if n == 3:       buf += "\033[1;35m◈\033[0m"
                elif p in sb:    buf += "\033[34m·\033[0m"
                elif p in sd:    buf += "\033[33m·\033[0m"
                elif p in sa:    buf += "\033[36m·\033[0m"
                else:            buf += " "
        buf += "\033[1;37m║\033[0m\n"
    return buf + "\033[1;37m╚" + rep("═",COLS) + "╝\033[0m\n"

def frame_play(pos, fog_on):
    buf = "\033[H\033[1;37m╔" + rep("═",COLS) + "╗\033[0m\n"
    for r in range(ROWS):
        buf += "\033[1;37m║\033[0m"
        for c in range(COLS):
            if fog_on and not seen[r][c]: buf += "\033[30m░\033[0m"
            elif (r,c) == pos:            buf += "\033[1;92m@\033[0m"
            elif maze[r][c] == '#':       buf += "\033[90m█\033[0m"
            elif maze[r][c] == 'S':       buf += "\033[1;32mS\033[0m"
            elif maze[r][c] == 'E':       buf += "\033[1;31mE\033[0m"
            else:                         buf += " "
        buf += "\033[1;37m║\033[0m\n"
    return buf + "\033[1;37m╚" + rep("═",COLS) + "╝\033[0m\n"

# ── Animation + heatmap ──────────────────────────────────────────────────────

def animate():
    mx = max(len(bp), len(dp), len(ap))
    print("\033[2J", end="")
    for s in range(mx+1):
        b = min(s, len(bp)); d = min(s, len(dp)); a = min(s, len(ap))
        print(frame_algo(b,d,a), end="")
        print(f"\033[34mBFS\033[0m:{b} \033[33mDFS\033[0m:{d} \033[36mA*\033[0m:{a}  {s}/{mx}")
        time.sleep(0.025)

def heatmap():
    print("\033[2J\033[H", end="")
    mx = max(max(row) for row in heat) or 1
    print("\033[1;37m╔" + rep("═",COLS) + "╗\033[0m")
    for r in range(ROWS):
        print("\033[1;37m║\033[0m", end="")
        for c in range(COLS):
            if maze[r][c] == '#': print("\033[90m█\033[0m", end=""); continue
            lvl = heat[r][c]*5//mx
            syms = ["\033[96m░\033[0m","\033[34m▒\033[0m","\033[32m▓\033[0m","\033[33m▓\033[0m","\033[31m█\033[0m"]
            print(syms[min(lvl,4)], end="")
        print("\033[1;37m║\033[0m")
    print("\033[1;37m╚" + rep("═",COLS) + "╝\033[0m")

# ── Raw keyboard input ───────────────────────────────────────────────────────

def read_key():
    import tty, termios
    fd = sys.stdin.fileno()
    old = termios.tcgetattr(fd)
    try:
        tty.setraw(fd)
        ch = sys.stdin.read(1)
        if ch == '\x1b':
            ch2 = sys.stdin.read(2)
            return ch + ch2
        return ch
    finally:
        termios.tcsetattr(fd, termios.TCSADRAIN, old)

# ── Play ─────────────────────────────────────────────────────────────────────

def play(fog_on):
    global steps, won, tout, score, best
    print("\033[2J\033[?25l", end="")
    pos = (1,1); steps = 0; won = tout = False
    end = (ROWS-2, COLS-2); t0 = time.time()
    if fog_on: fog_reveal(pos)
    while True:
        tl = 60 - int(time.time()-t0)
        if tl <= 0: tout = True; break
        col = 32 if tl > 10 else 31
        ftag = " [FOG]" if fog_on else ""
        print(frame_play(pos, fog_on), end="")
        print(f"  Steps:{steps:<4}  Time:\033[{col}m{tl}s\033[0m{ftag}  WASD=move  Q=quit")
        try:
            k = read_key()
        except Exception:
            k = input()
        np = pos
        if k in ('w','W','\x1b[A'): np = (pos[0]-1, pos[1])
        if k in ('s','S','\x1b[B'): np = (pos[0]+1, pos[1])
        if k in ('d','D','\x1b[C'): np = (pos[0], pos[1]+1)
        if k in ('a','A','\x1b[D'): np = (pos[0], pos[1]-1)
        if k in ('q','Q','\x03'):   print("\033[?25h", end=""); return
        if ok(*np) and maze[np[0]][np[1]] != '#':
            steps += 1; pos = np
            if fog_on: fog_reveal(pos)
        if pos == end:
            won = True
            v = max(50, 1000-(steps-len(bp))*5+tl)
            score = v
            if score > best: best = score
            print("\033[2J\033[H")
            print(f"\n  \033[1;92m★ YOU WIN!\033[0m  Steps:{steps:<4}  Score:\033[1;93m{score}\033[0m/1000  Optimal:{len(bp)}")
            if score > 800: print("  \033[1;92m⚡ ELITE!\033[0m")
            input("\n  [Enter]...")
            print("\033[?25h", end=""); return
    print("\033[2J\033[H\n  \033[1;31m⏱ TIME UP!\033[0m")
    input("\n  [Enter]...")
    print("\033[?25h", end="")

# ── Results ───────────────────────────────────────────────────────────────────

def results():
    print("\033[2J\033[H", end="")
    print("\033[1;37m╔────── ALGORITHM RESULTS ──────╗\033[0m")
    print(f"\033[34m  BFS\033[0m  path={len(bp):<4}  {bt*1000:.2f}ms")
    print(f"\033[33m  DFS\033[0m  path={len(dp):<4}  {dt*1000:.2f}ms")
    print(f"\033[36m  A* \033[0m  path={len(ap):<4}  {at*1000:.2f}ms")
    if won: print(f"\033[1;92m  YOU\033[0m  steps={steps:<4}  score=\033[1;93m{score}\033[0m/1000")
    if best: print(f"  \033[1;93m★ Best: {best}\033[0m")
    print("\033[1;37m╚───────────────────────────────╝\033[0m")
    s = input("\n  [R]eplay  [Enter] back: ")
    if s.lower() == 'r':
        animate()
        input("\n  [Enter]...")

# ── Reset ─────────────────────────────────────────────────────────────────────

def reset():
    global bp, dp, ap, won, tout, steps, score
    bp = dp = ap = []
    won = tout = False
    steps = score = 0
    for r in range(ROWS):
        for c in range(COLS):
            heat[r][c] = 0
            seen[r][c] = False

# ── Menu ──────────────────────────────────────────────────────────────────────

def menu():
    print("\033[2J\033[H", end="")
    print("\033[1;37m╔══════════════════════════════════╗")
    print("║   RECURSIVE MAZE — 4D 2026       ║")
    print(f"║  \033[1;93m★ Best: {best:<6}\033[1;37m                ║" if best else "║                                  ║")
    print("╠══════════════════════════════════╣")
    print("║  1.  Watch BFS / DFS / A*        ║")
    print("║  2.  Play — no fog               ║")
    print("║  3.  Play — fog of war           ║")
    print("║  4.  Quit                        ║")
    print("╚══════════════════════════════════╝\033[0m")
    try:
        return int(input("\n  Choice: "))
    except ValueError:
        return 0

# ── Main ──────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("\033[?25l", end="")
    try:
        while True:
            reset(); gen(); solve_all()
            choice = menu()
            if choice == 1:
                input("  [Enter] race...")
                animate()
                input("\n  [Enter] heatmap...")
                heatmap()
                input("\n  [Enter] results...")
                results()
            elif choice == 2:
                input("  [Enter] to play...")
                play(False); results()
            elif choice == 3:
                input("  [Enter] FOG mode...")
                play(True); results()
            elif choice == 4:
                print("\033[2J\033[H"); break
    finally:
        print("\033[?25h", end="")
