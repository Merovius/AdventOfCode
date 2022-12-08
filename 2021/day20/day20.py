from functools import reduce
import math

def read(name):
    with open(name, 'r') as file:
        lines = file.readlines()
        alg = lines[0].strip()
        pic = [l.strip() for l in lines[2:]]
        return Pic(alg, pic)

class Pic:
    def __init__(self, alg, pic):
        def bools(s):
            return list(c == '#' for c in s)
        self.alg = bools(alg)
        self.pic = []
        for r in range(len(pic)):
            self.pic.append(bools(pic[r]))
        self.oob = False

    def __getitem__(self, t):
        r, c = t
        if r < 0 or r >= len(self.pic):
            return self.oob
        if c < 0 or c >= len(self.pic[0]):
            return self.oob
        return self.pic[r][c]

    def env(self, r, c):
        s = []
        for ri in range(r-1, r+2):
            for ci in range(c-1, c+2):
                if self[ri, ci]:
                    s.append('1')
                else:
                    s.append('0')
        n = int(''.join(s), 2)
        return n

    def step(self):
        out = []
        for r in range(-1, len(self.pic)+2):
            row = []
            for c in range(-1, len(self.pic[0])+2):
                row.append(self.alg[self.env(r, c)])
            out.append(row)
        self.oob = self.alg[self.env(-10, -10)]
        self.pic = out

    def dump(self):
        for r in range(-1, len(self.pic)+1):
            for c in range(-1, len(self.pic[0])+1):
                if self[r, c]:
                    print('#', end='')
                else:
                    print('.', end='')
            print()
        print()

    def alive(self):
        n = 0
        for r in range(len(self.pic)):
            for c in range(len(self.pic[0])):
                if self[r, c]:
                    n += 1
        return n

pic = read('input.txt')
pic.step()
pic.step()
print(f"{pic.alive()} alive cells after 2 steps")
for i in range(48):
    pic.step()
print(f"{pic.alive()} alive cells after 50 steps")
