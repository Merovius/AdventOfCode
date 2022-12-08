def read(name):
    with open(name, 'r') as file:
        return Grid([list(l.strip()) for l in file.readlines()])

class Grid:
    def __init__(self, lines):
        self._lines = list(lines)
        self._rows = len(self._lines)
        self._cols = len(self._lines[0])

    def __getitem__(self, p):
        r, c = p
        return self._lines[r % self._rows][c % self._cols]

    def __setitem__(self, p, v):
        r, c = p
        self._lines[r % self._rows][c % self._cols] = v

    def __repr__(self):
        return '\n'.join(''.join(l) for l in self._lines)

    def step(self):
        changed = False
        new = Grid(['.'] * self._cols for _ in range(self._rows))
        for r in range(self._rows):
            for c in range(self._cols):
                if self[r, c] == '>' and self[r, c+1] == '.':
                    new[r, c+1] = '>'
                    changed = True
                    continue
                if new[r, c] == '.':
                    new[r, c] = self[r, c]
        self._lines = new._lines

        new = Grid(['.'] * self._cols for _ in range(self._rows))
        for r in range(self._rows):
            for c in range(self._cols):
                if self[r, c] == 'v' and self[r+1, c] == '.':
                    new[r+1, c] = 'v'
                    changed = True
                    continue
                if new[r, c] == '.':
                    new[r, c] = self[r, c]
        self._lines = new._lines

        return changed

    def stabilize(self):
        n = 1
        while self.step():
            n += 1
        return n

if __name__ == "__main__":
    g = read('input.txt')
    print(f"{g.stabilize()} steps needed to stabilize")
