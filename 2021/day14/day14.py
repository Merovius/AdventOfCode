from collections import Counter
import math

def read(name):
    def parse_rule(line):
        src, dst = line.strip().split(" -> ")
        return (src, dst)

    with open(name, 'r') as file:
        lines = file.readlines()
        tpl = lines[0].strip()
        if lines[1] != "\n":
            raise ValueError
        rules = dict(parse_rule(l) for l in lines[2:])
        return (tpl, rules)

tpl, rules = read('input.txt')

class Counts:
    def __init__(self, tpl, rules):
        self._last_letter = tpl[-1]
        self._rules = rules
        self._counts = Counter(p[0]+p[1] for p in zip(tpl, tpl[1:]))

    def score(self):
        counts = self.letter_counts()
        return (max(counts.values())-min(counts.values()))

    def letter_counts(self):
        counts = Counter()
        for k, v in self._counts.items():
            counts.update({k[0]: v})
        counts.update(self._last_letter)
        return counts

    def step(self):
        new = Counter()
        for k, v in self._counts.items():
            if k in self._rules:
                new.update({k[0]+self._rules[k]: v})
                new.update({self._rules[k]+k[1]: v})
            else:
                new.update({k: v})
        self._counts = new

counts = Counts(tpl, rules)
for i in range(10):
    counts.step()
print(f"After 10 steps, the difference between most and least common element is {counts.score()}")
for i in range(30):
    counts.step()
print(f"After 40 steps, the difference between most and least common element is {counts.score()}")
