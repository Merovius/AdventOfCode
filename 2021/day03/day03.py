from collections import defaultdict
from itertools import filterfalse

def readfile(name):
    with open(name, 'r') as file:
        return file.readlines()

data = [l.strip() for l in readfile('input.txt')]
#data = ['00100','11110','10110','10111','10101','01111','00111','11100','10000','11001','00010','01010']

total = 0
ones = defaultdict(lambda: 0)
for line in data:
    total += 1
    for i in range(len(line)):
        ones[i] += int(line[i])
gamma = 0
epsilon = 0
for i in range(len(ones)):
    gamma = gamma << 1
    epsilon = epsilon << 1
    if ones[i] > total/2:
        gamma += 1
    else:
        epsilon += 1
print(f'γ: {gamma:b}={gamma}, ε: {epsilon:b}={epsilon}, γ•ε: {gamma*epsilon}')

def most_common_bit(lines, i):
    n = 0
    for l in lines:
        if l[i] == '1':
            n += 1
    if n >= len(lines)/2:
        return '1'
    return '0'

filtered0 = data
filtered1 = data
for i in range(len(ones)):
    bit0 = most_common_bit(filtered0, i)
    bit1 = most_common_bit(filtered1, i)
    if len(filtered0) > 1:
        filtered0 = [l for l in filtered0 if l[i] == bit0]
    if len(filtered1) > 1:
        filtered1 = [l for l in filtered1 if l[i] != bit1]
    if len(filtered0) == 1 and len(filtered1) == 1:
        break

o2 = int(filtered0[0], 2)
co2 = int(filtered1[0], 2)
print(f'O₂: {o2:b}={o2}, CO₂: {co2:b}={co2}, product: {o2*co2}')
