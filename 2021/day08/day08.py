from itertools import permutations

def parse_line(line):
    w = line.split()
    return w[:10]+w[-4:]


def parse(name):
    def parse_line(line):
        w = line.split()
        return w[:10]+w[-4:]

    with open(name, 'r') as file:
        return [parse_line(l) for l in file.readlines()]

data = parse("input.txt")

# The segments turned on for each digit
segments = [
    set("abcefg"),
    set("cf"),
    set("acdeg"),
    set("acdfg"),
    set("bcdf"),
    set("abdfg"),
    set("abdefg"),
    set("acf"),
    set("abcdefg"),
    set("abcdfg"),
]
digits = {
    frozenset("abcefg"): 0,
    frozenset("cf"): 1,
    frozenset("acdeg"): 2,
    frozenset("acdfg"): 3,
    frozenset("bcdf"): 4,
    frozenset("abdfg"): 5,
    frozenset("abdefg"): 6,
    frozenset("acf"): 7,
    frozenset("abcdefg"): 8,
    frozenset("abcdfg"): 9,
}
digit_by_len = {
    2: 1,
    4: 4,
    3: 7,
    8: 8,
}

simple_digits = [s for x in data for s in x[-4:] if len(s) in {2,4,3,7}]
print(f"Simple digits in output: {len(simple_digits)}")

def possible_mappings(entry):
    # m maps the possible segments that could be connected to each wire
    m = {}
    for x in "abcdefg":
        m[x] = set("abcdefg")

    for s in entry:
        if not len(s) in digit_by_len:
            continue
        d = digit_by_len[len(s)]
        for c in s:
            m[c] = m[c].intersection(segments[d])

    assigned = set()
    for a in m['a']:
        assigned.add(a)
        for b in m['b'].difference(assigned):
            assigned.add(b)
            for c in m['c'].difference(assigned):
                assigned.add(c)
                for d in m['d'].difference(assigned):
                    assigned.add(d)
                    for e in m['e'].difference(assigned):
                        assigned.add(e)
                        for f in m['f'].difference(assigned):
                            assigned.add(f)
                            for g in m['g'].difference(assigned):
                                yield a+b+c+d+e+f+g
                            assigned.remove(f)
                        assigned.remove(e)
                    assigned.remove(d)
                assigned.remove(c)
            assigned.remove(b)
        assigned.remove(a)

def apply_mapping(entry, m):
    out = []
    for s in entry:
        out.append(''.join(m[ord(e)-ord('a')] if e in "abcdefg" else e for e in s))
    return out

def is_valid(entry):
    for e in entry:
        if not frozenset(e) in digits:
            return False
    return True

def output_value(entry):
    for m in possible_mappings(entry):
        mapped = apply_mapping(entry, m)
        if not is_valid(mapped):
            continue
        s = ''.join(str(digits[frozenset(s)]) for s in mapped[-4:])
        return int(s)

entry = parse_line('acedgfb cdfbe gcdfa fbcad dab cefabd cdfgeb eafb cagedb ab | cdfeb fcadb cdfeb cdbaf')
total = sum(output_value(entry) for entry in data)
print(f"Sum of output values: {total}")
