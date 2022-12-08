from collections import namedtuple
from queue import PriorityQueue
import math

def read(name):
    with open(name, 'r') as file:
        return [list(map(int, l.strip())) for l in file.readlines()]

data = read('input.txt')

QEntry = namedtuple("QEntry", ["cost", "node"])

class Graph:
    def __init__(self, nodes):
        self._nodes = nodes
        self._rows = len(nodes)
        self._cols = len(nodes[0])
        self.start = (0,0)
        self.end = (self._rows-1,self._cols-1)

    def level(self, n):
        return self._nodes[n[0]][n[1]]

    def neighbors(self, n):
        if n[0] > 0:
            yield(n[0]-1, n[1])
        if n[0] < self._rows-1:
            yield(n[0]+1, n[1])
        if n[1] > 0:
            yield(n[0], n[1]-1)
        if n[1] < self._cols-1:
            yield(n[0], n[1]+1)

    def shortest_path_cost(self):
        # implements Dijkstra's Algorithm
        q = PriorityQueue()
        # we are only interested in the cost, not the actual path, so we don't
        # store where we visited a node *from*. Furthermore, the edge-weight
        # only depends on the target node, so we know that the first time we
        # see a node, we did so following the shortest path. So we can mark it
        # as visited as soon as we see it the first time, making sure we only
        # queue every node at most once.
        visited = set()

        q.put(QEntry(cost=0, node=self.start))
        while not q.empty():
            e = q.get()
            if e.node == self.end:
                return e.cost
            for n in self.neighbors(e.node):
                if n in visited:
                    continue
                visited.add(n)
                q.put(QEntry(cost=e.cost+self.level(n), node=n))
        raise ValueError("no path found")

def expand(data):
    out = []
    for i in range(5):
        for row in data:
            outrow = []
            for j in range(5):
                for col in row:
                    v = col+i+j
                    while v > 9:
                        v -= 9
                    outrow.append(v)
            out.append(outrow)
    return out

g = Graph(data)
print(f"Shortest path has total risk {g.shortest_path_cost()}")
data = expand(data)
g = Graph(data)
print(f"Shortest path in full cave has total risk {g.shortest_path_cost()}")
