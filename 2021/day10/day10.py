from statistics import median

def read(name):
    with open(name, 'r') as file:
        return [l.strip() for l in file.readlines()]

match = {"(": ")", "[": "]", "{": "}", "<": ">"}
scores = {")": 3, "]": 57, "}": 1197, ">": 25137, "(": 1, "[": 2, "{": 3, "<": 4}

def score(line):
    stack = []
    for c in line:
        if c in "([{<":
            stack.append(c)
            continue
        top = stack.pop()
        if c == match[top]:
            continue
        return (scores[c], 0)
    score = 0
    while len(stack) > 0:
        score = score*5+scores[stack.pop()]
    return (0, score)

data = read('input.txt')
scores = [score(l) for l in data]
print(f"Score of corrupted lines: {sum(s[0] for s in scores)}")
print(f"Score of incomplete lines: {median([s[1] for s in scores if s[1] != 0])}")
