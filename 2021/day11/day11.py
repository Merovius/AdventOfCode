def read(name):
    with open(name, 'r') as file:
        return [list(map(int, l.strip())) for l in file.readlines()]

data = read('example.txt')

class Point(tuple):
    def __new__(self, r, c):
        return tuple.__new__(Point, (r, c))

    def __add__(self, other):
        return Point(self[0]+other[0], self[1]+other[1])

    @property
    def valid(self):
        return self[0] >= 0 and self[0] < len(data) and self[1] >= 0 and self[1] < len(data[self[0]])

    @property
    def val(self):
        return data[self[0]][self[1]]

    @val.setter
    def val(self, value):
        data[self[0]][self[1]] = value

    @property
    def neighbors(self):
        for n in [Point(-1,-1),Point(-1,0),Point(-1,1),Point(0,-1),Point(0,1),Point(1,-1),Point(1,0),Point(1,1)]:
            p = self+n
            if p.valid:
                yield p

def points():
    for r in range(len(data)):
        for c in range(len(data[r])):
            yield Point(r, c)

total = sum(len(r) for r in data)
flashes = 0

def step():
    global flashes
    flashed = []
    for p in points():
        p.val += 1
        if p.val == 10:
            flashed.append(p)
            flashes += 1
    while len(flashed) > 0:
        p = flashed.pop()
        for q in p.neighbors:
            q.val += 1
            if q.val == 10:
                flashed.append(q)
                flashes += 1
    for p in points():
        if p.val > 9:
            p.val = 0

def print_data():
    for r in range(len(data)):
        for c in range(len(data[r])):
            print(f"{data[r][c]}", end='')
        print('')
    print('')

simuflash = -1
for i in range(100):
    before = flashes
    step()
    if flashes-before == total:
        simuflash = i+1
print(f"After 100 steps there where {flashes} flashes")
if simuflash > 0:
    print(f"Simultaneous flash at step {simuflash}")
else:
    for i in range(100, 1000):
        before = flashes
        step()
        if flashes-before == total:
            print(f"Simultaneous flash at step {i+1}")
            break
