from collections import namedtuple

Range = namedtuple('Range', ['min', 'max'])
Step = namedtuple('Step', ['state', 'cuboid'])

def read(name):
    def parse_range(s):
        if len(s) < 2:
            raise ValueError(f"range string too short")
        if s[0] not in {'x', 'y', 'z'}:
            raise ValueError(f"invalid range axis {s[0]}")
        if s[1] != '=':
            raise ValueError(f"invalid range")
        s = s[2:]
        min, max = s.split('..')
        # We use ranges to be exclusive, as that makes the math easier
        return Range(min=int(min), max=int(max)+1)

    def parse_cuboid(s):
        xr, yr, zr = s.split(",")
        return Cuboid(parse_range(xr), parse_range(yr), parse_range(zr))

    def parse_step(l):
        l = l.strip()
        if l.startswith('on '):
            state = True
            l = l[3:]
        elif l.startswith('off '):
            state = False
            l = l[4:]
        else:
            raise ValueError(f"invalid line {l}")
        return Step(state=state, cuboid=parse_cuboid(l))

    with open(name, 'r') as file:
        return [parse_step(l) for l in file.readlines()]

class Cuboid:
    def __init__(self, x, y, z):
        self.x, self.y, self.z = x, y, z

    # Returns the (possibly empty) intersection between self and other.
    def intersect(self, other):
        x = Range(max(self.x.min, other.x.min), min(self.x.max, other.x.max))
        y = Range(max(self.y.min, other.y.min), min(self.y.max, other.y.max))
        z = Range(max(self.z.min, other.z.min), min(self.z.max, other.z.max))
        return Cuboid(x,y,z)

    # Returns an iterator over cuboids which are in self, but not in other.
    def difference(self, other):
        c = self.intersect(other)
        if c.empty():
            yield self
            return
        def maybe(xmin, xmax, ymin, ymax, zmin, zmax):
            c = Cuboid(Range(xmin, xmax), Range(ymin, ymax), Range(zmin, zmax))
            if not c.empty():
                yield c

        # See difference.svg/png
        yield from maybe(self.x.min, self.x.max, self.y.min, self.y.max, self.z.min, c.z.min)
        yield from maybe(self.x.min, self.x.max, self.y.min, c.y.min, c.z.min, c.z.max)
        yield from maybe(self.x.min, c.x.min, c.y.min, c.y.max, c.z.min, c.z.max)
        yield from maybe(self.x.min, self.x.max, self.y.min, self.y.max, c.z.max, self.z.max)
        yield from maybe(self.x.min, self.x.max, c.y.max, self.y.max, c.z.min, c.z.max)
        yield from maybe(c.x.max, self.x.max, c.y.min, c.y.max, c.z.min, c.z.max)

    def empty(self):
        return self.x.min >= self.x.max or self.y.min >= self.y.max or self.z.min >= self.z.max

    def size(self):
        if self.empty():
            return 0
        return (self.x.max-self.x.min)*(self.y.max-self.y.min)*(self.z.max-self.z.min)

    def __repr__(self):
        return f"[{self.x.min},{self.x.max})×[{self.y.min},{self.y.max})×[{self.z.min},{self.z.max})"


init_region = Cuboid(Range(-50,51), Range(-50,51), Range(-50,51))
steps = read('input.txt')

def execute(steps):
    on = set()
    for s in steps:
        new = set()
        for c in on:
            for cd in c.difference(s.cuboid):
                new.add(cd)
        if s.state:
            new.add(s.cuboid)
        on = new
    return on

on = execute(s for s in steps if not s.cuboid.intersect(init_region).empty())
N = sum(c.size() for c in on)
print(f"After initialization {N} cubes are on")

on = execute(steps)
N = sum(c.size() for c in on)
print(f"After full reboot {N} cubes are on")
