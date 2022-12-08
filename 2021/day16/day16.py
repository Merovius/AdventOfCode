from functools import reduce
from operator import mul,gt,lt,eq

def decode_hex(s):
    d = {
        '0': '0000', '1': '0001', '2': '0010', '3': '0011',
        '4': '0100', '5': '0101', '6': '0110', '7': '0111',
        '8': '1000', '9': '1001', 'A': '1010', 'B': '1011',
        'C': '1100', 'D': '1101', 'E': '1110', 'F': '1111',
    }
    return reduce(lambda a, b: a+b, (d[c] for c in s))

class Packet:
    def __init__(self, s, bits=False):
        if not bits:
            s = decode_hex(s)
        self._bits = s
        self._len = 0
        self.version = int(self._consume(3), 2)
        self.id = int(self._consume(3), 2)
        if self.id == 4:
            self.val = 0
            cont = True
            while cont:
                g = self._consume(5)
                cont = int(g[0])
                self.val = (self.val<<4) + int(g[1:], 2)
            delattr(self, "_bits")
            return
        mode = int(self._consume(1)[0])
        self.sub = []
        if mode:
            m = int(self._consume(11), 2)
            for _ in range(m):
                p = Packet(self._bits, bits=True)
                self._consume(p._len)
                delattr(p, "_len")
                self.sub.append(p)
        else:
            m = int(self._consume(15), 2)
            while m > 0:
                p = Packet(self._bits, bits=True)
                self._consume(p._len)
                m -= p._len
                delattr(p, "_len")
                self.sub.append(p)
        delattr(self, "_bits")

    def _consume(self, n):
        v = self._bits[:n]
        self._len += n
        self._bits = self._bits[n:]
        return v

    def __repr__(self):
        if self.id == 4:
            return f"Literal(ver={self.version}, val={self.val})"
        sub = ', '.join(p.__repr__() for p in self.sub)
        return f"Operator(ver={self.version}, id={self.id}, sub=({sub})"

    def version_sum(self):
        if self.id == 4:
            return self.version
        return self.version + sum(p.version_sum() for p in self.sub)

    def eval(self):
        if self.id == 4:
            return self.val
        ops = {
            0: sum,
            1: lambda l: reduce(mul, l, 1),
            2: min,
            3: max,
            5: lambda l: int(l[0] > l[1]),
            6: lambda l: int(l[0] < l[1]),
            7: lambda l: int(l[0] == l[1]),
        }
        return ops[self.id]([p.eval() for p in self.sub])

with open('input.txt') as file:
    packet = Packet(file.read().strip())
print(f"Version sum: {packet.version_sum()}")
print(f"Evaluated: {packet.eval()}")
