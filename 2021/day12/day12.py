from collections import defaultdict
from collections import Counter

def read(name):
    _map = defaultdict(lambda: set())
    with open(name, 'r') as file:
        for l in file.readlines():
            start, end = l.strip().split("-")
            _map[start].add(end)
            _map[end].add(start)
    return _map

def find_paths(_map, abort):
    paths = []
    def dfs(node='start', path=['start']):
        if node == 'end':
            paths.append(list(path)+[node])
            return
        if abort(path):
            return

        for neighbor in _map[node]:
            path.append(neighbor)
            dfs(neighbor, path)
            path.pop()

    dfs()
    return paths

def part1(path):
    visited = Counter(p for p in path if p.islower())
    return any(v > 1 for v in visited.values())

def part2(path):
    visited = Counter(p for p in path if p.islower())
    if visited['start'] > 1:
        return True
    if any(v > 2 for v in visited.values()):
        return True
    return sum(1 for v in visited.values() if v > 1) > 1

_map = read('input.txt')
print(f"Paths for part 1: {sum(1 for _ in find_paths(_map, part1))}")
print(f"Paths for part 2: {sum(1 for _ in find_paths(_map, part2))}")
