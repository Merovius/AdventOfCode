from re import fullmatch
from collections import defaultdict

def parse(name):
    lines = open(name).readlines()
    groups = [fullmatch('(\d+),(\d+) -> (\d+),(\d+)', l.strip()).groups() for l in lines]
    return [((int(g[0]), int(g[1])),(int(g[2]),int(g[3]))) for g in groups]

def print_hits(hits):
    max_x, max_y = 0, 0
    for p in hits:
        max_x = max(max_x, p[0])
        max_y = max(max_y, p[1])
    for y in range(0, max_y+1):
        for x in range(0, max_x+1):
            if hits[(x,y)] == 0:
                print('.', end='')
            else:
                print(str(hits[(x,y)]), end='')
        print('')

def sign(n):
    if n<0:
        return -1
    return int(n>0)

def calculate_hits(lines, skip_diagonals=True):
    def sign(n):
        if n < 0:
            return -1
        return int(n>0)
    hits = defaultdict(lambda: 0)
    for l in lines:
        (xa, ya) = l[0]
        (xb, yb) = l[1]
        if skip_diagonals and xa != xb and ya != yb:
            continue
        xs = sign(xb-xa)
        ys = sign(yb-ya)
        x, y = xa, ya
        while True:
            hits[(x,y)] += 1
            if x == xb and y == yb:
                break
            x += xs
            y += ys
    return hits

def count_multihits(hits):
    return sum(1 for n in hits.values() if n > 1)

lines = parse('input.txt')
hits = calculate_hits(lines, True)
print(f"Number of overlap-points: {count_multihits(hits)}")
hits = calculate_hits(lines, False)
print(f"Number of overlap-points with diagonals: {count_multihits(hits)}")
