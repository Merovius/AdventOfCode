from collections import defaultdict

def to_dict(fish):
    d = defaultdict(lambda: 0)
    for f in fish:
        d[f] += 1
    return d

def readfile(name):
    with open(name, 'r') as file:
        return to_dict(list(map(int, file.readlines()[0].strip().split(","))))

fish = readfile('input.txt')

def step(fish):
    new = defaultdict(lambda: 0)
    for days, n in fish.items():
        if days == 0:
            new[8] += n
            new[6] += n
        else:
            new[days-1] += n
    return new

for i in range(256):
    if i == 80:
        print(f"After 80 days there are {sum(fish.values())} fish")
    fish = step(fish)
print(f"After 256 days there are {sum(fish.values())} fish")
