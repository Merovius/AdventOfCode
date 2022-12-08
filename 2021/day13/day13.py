def read(name):
    def parse_point(l):
        x, y = l.strip().split(",")
        return (int(x), int(y))

    def parse_fold(l):
        if not l.startswith("fold along "):
            raise ValueError
        axis, val = l[len("fold along "):].strip().split("=")
        return (axis, int(val))

    with open(name, 'r') as file:
        lines = file.readlines()
        i = lines.index('\n')
        points = set(parse_point(l) for l in lines[:i])
        folds = [parse_fold(l) for l in lines[i+1:]]
        return (points, folds)

def apply_fold(points, fold):
    axis, v = fold
    applied = set()
    for p in points:
        if axis == 'x' and p[0] > v:
            applied.add((2*v-p[0], p[1]))
        elif axis == 'y' and p[1] > v:
            applied.add((p[0], 2*v-p[1]))
        else:
            applied.add(p)
    return applied

def print_points(points):
    mx, my = max(p[0] for p in points), max(p[1] for p in points)
    for y in range(my+1):
        for x in range(mx+1):
            if (x,y) in points:
                print("•", end='')
            else:
                print(" ", end='')
        print("")

points, folds = read('input.txt')
points = apply_fold(points, folds[0])
print(f"After first fold there are {len(points)} points")
for f in folds[1:]:
    points = apply_fold(points, f)
print(f"Paper after all folds:")
print_points(points)
