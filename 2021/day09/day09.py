from functools import reduce
from operator import mul

def read(name):
    with open(name, 'r') as file:
        return [[int(c) for c in l.strip()] for l in file.readlines()]

data = read('input.txt')

class Point(tuple):
    def __new__(self, r, c):
        return tuple.__new__(Point, (r, c))

    @property
    def val(self):
        return data[self[0]][self[1]]

    @property
    def neighbors(self):
        r, c = self
        if r > 0:
            yield Point(r-1, c)
        if r < len(data)-1:
            yield Point(r+1, c)
        if c > 0:
            yield Point(r, c-1)
        if c < len(data[r])-1:
            yield Point(r, c+1)

def points():
    for r in range(len(data)):
        for c in range(len(data[r])):
            yield Point(r, c)

def risk_level(data):
    return sum(1+p.val for p in points() if all(p.val < q.val for q in p.neighbors))

def basin_sizes(data):
    labels = {}
    sizes = []
    def dfs(p, label=None):
        if p in labels or p.val >= 9:
            return
        if label is None:
            label = len(sizes)
            sizes.append(0)
        labels[p] = label
        sizes[label] += 1
        for q in p.neighbors:
            dfs(q, label)

    for p in points():
        dfs(p)

    return reduce(mul, sorted(sizes)[-3:])

print(f"Risk level: {risk_level(data)}")
print(f"Product of largest basin sizes: {basin_sizes(data)}")

