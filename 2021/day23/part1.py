from dataclasses import dataclass, field
from typing import Any
from queue import PriorityQueue

# Occupiable spaces are numbered left to right, top to bottom:
#•••••••••••••••••••
#•0 1. 2. 3. 4. 5 6•
#••••7 •8 •9 •10••••
#   •11•12•13•14•
#   •••••••••••••

# Path lengths of direct connections between spaces:
edges = {
    (0,1): 1,
    (1,0): 1, (1,2): 2, (1,7): 2,
    (2,1): 2, (2,3): 2, (2,7): 2, (2,8): 2,
    (3,2): 2, (3,4): 2, (3,8): 2, (3,9): 2,
    (4,3): 2, (4,5): 2, (4,9): 2, (4,10): 2,
    (5,4): 2, (5,10): 2, (5,6): 1,
    (6,5): 1,
    (7,1): 2, (7,2): 2, (7,11): 1,
    (8,2): 2, (8,3): 2, (8,12): 1,
    (9,3): 2, (9,4): 2, (9,13): 1,
    (10,4): 2, (10,5): 2, (10,14): 1,
    (11,7): 1,
    (12,8): 1,
    (13,9): 1,
    (14,10): 1,
}
# Movement cost per amphipod type
cost = {'A': 1, 'B': 10, 'C': 100, 'D': 1000}
# Homes of the respective amphipod type
homes = {'A': {7,11}, 'B': {8,12}, 'C': {9,13}, 'D': {10,14}}
def is_hallway(c):
    return c < 7

@dataclass(order=True)
class QEntry:
    cost: int
    state: Any=field(compare=False)

class State:
    # Create a new State, by giving a list of cells with their content
    def __init__(self, cells):
        if isinstance(cells, dict):
            l = [None for _ in range(15)]
            for i, a in cells.items():
                l[i] = a
            self._cells = l
            self._moves = []
            return
        if len(cells) != 15:
            raise ValueError
        self._cells = list(cells)
        self._moves = []

    # Move from src to dst and return the cost of that move and the path taken.
    # Raises a ValueError, if the move is illegal.
    def move(self, src, dst):
        if self._cells[src] is None:
            raise ValueError('source cell is unoccupied')
        if self._cells[dst] is not None:
            raise ValueError('destination cell is occupied')
        a = self._cells[src]
        if dst not in homes[a]:
            if is_hallway(src):
                raise ValueError('amphipod stands in hallway and can only move home')
            elif not is_hallway(dst):
                raise ValueError('amphipod can only move home or into hallway')
        else:
            for c in homes[a]:
                if self._cells[c] != None and self._cells[c] != a:
                    raise ValueError('amphipods home is occupied by wrong amphipod type')
        p = self.path(src, dst)
        if p is None:
            raise ValueError('no path from source to destination cell')
        a = self._cells[src]
        c = cost[a] * sum(edges[e] for e in zip(p, p[1:]))
        self._cells[src] = None
        self._cells[dst] = a
        self._moves.append((src, dst, c))
        print(self)

    def undo(self):
        m = self._moves.pop()
        t = self._cells[m[0]]
        self._cells[m[0]] = self._cells[m[1]]
        self._cells[m[1]] = t
        print(self)

    def total_cost(self):
        return sum(m[2] for m in self._moves)

    def path(self, src, dst):
        q = PriorityQueue()
        visited = dict()
        q.put((0,None,src))
        while not q.empty():
            e = q.get()
            cost, prev, cur = e
            if cur in visited:
                continue
            visited[cur] = prev
            if cur == dst:
                break
            for i in range(15):
                if (cur, i) not in edges or self._cells[i] is not None:
                    continue
                if i in visited:
                    continue
                q.put((cost+edges[(cur,i)], cur, i))
        if dst not in visited:
            raise ValueError('no path from src to dst')
        path = []
        while dst is not None:
            path.append(dst)
            dst = visited[dst]
        path.reverse()
        return path

    def is_finished(self):
        return all(self._cells[i] == a for a, home in homes.items() for i in home)

    def __repr__(self):
        def char(x):
            return '.' if x is None else x

        s = '#############\n#'
        s += char(self._cells[0])
        s += char(self._cells[1])
        s += '.'
        s += char(self._cells[2])
        s += '.'
        s += char(self._cells[3])
        s += '.'
        s += char(self._cells[4])
        s += '.'
        s += char(self._cells[5])
        s += char(self._cells[6])
        s += '#\n###'
        s += char(self._cells[7])
        s += '#'
        s += char(self._cells[8])
        s += '#'
        s += char(self._cells[9])
        s += '#'
        s += char(self._cells[10])
        s += '###\n  #'
        s += char(self._cells[11])
        s += '#'
        s += char(self._cells[12])
        s += '#'
        s += char(self._cells[13])
        s += '#'
        s += char(self._cells[14])
        s += '#\n  #########\n'
        s += str(self.total_cost())
        return s

end = State({
    7:  'A', 8:  'B', 9:  'C', 10: 'D',
    11: 'A', 12: 'B', 13: 'C', 14: 'D',
})

# example:
#############
#...........#
###B#C#B#D###
  #A#D#C#A#
  #########
example = State({
    7:  'B', 8:  'C', 9:  'B', 10: 'D',
    11: 'A', 12: 'D', 13: 'C', 14: 'A',
})

# input:
#############
#...........#
###D#A#C#C###
  #D#A#B#B#
  #########
input = State({
    7:  'D', 8:  'A', 9:  'C', 10: 'C',
    11: 'D', 12: 'A', 13: 'B', 14: 'B',
})



