from collections import namedtuple
from math import inf

Range = namedtuple("Range", ["min", "max"])

def read(name):
    def parse_range(s):
        s = s[2:]
        start, end = s.split("..")
        return Range(int(start), int(end))

    with open(name, 'r') as file:
        data = file.read().strip()
        data = data[len("target area: "):]
        xrange, yrange = map(parse_range, data.split(", "))
        return (xrange, yrange)

xtarget, ytarget = read('input.txt')

def sign(x):
    if x > 0:
        return 1
    if x < 0:
        return -1
    return 0

class Probe:
    def __init__(self, xtarget, ytarget, vx, vy):
        self.target = (xtarget, ytarget)
        self.pos = (0,0)
        self.v = (vx, vy)
        self.maxy = -inf
        self.hit = False
        self.overshot = False

    def step(self):
        self.pos = (self.pos[0]+self.v[0], self.pos[1]+self.v[1])
        self.v = (max(0, self.v[0]-1), self.v[1]-1)
        self.maxy = max(self.maxy, self.pos[1])

        x, y = self.target
        if self.pos[0] >= x.min and self.pos[0] <= x.max and self.pos[1] >= y.min and self.pos[1] <= y.max:
            self.hit = True
        if self.pos[0] > x.max or self.pos[1] < y.min:
            self.overshot = True

    def shoot(self):
        while (not self.hit) and (not self.overshot):
            self.step()
        if self.hit:
            return self.maxy
        return None

maxy = -inf
best = (0,0)
ngood = 0
for vx in range(1,1000):
    for vy in range(-1000,1000):
        my = Probe(xtarget, ytarget, vx, vy).shoot()
        if my is not None:
            ngood += 1
            if my > maxy:
                maxy = my
                best = (vx, vy)
print(f"Highest is {maxy} towards {best}")
print(f"There are {ngood} good vectors")
