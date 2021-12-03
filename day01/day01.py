data = [int(l) for l in open('input.txt', 'r').readlines()]
#data = [199,200,208,210,200,207,240,269,260,263]

count = sum(w[0] < w[1] for w in zip(data, data[1:]))
print(f'Number of increases: {count}')
count = sum(w[0] < w[1] for w in zip(data, data[3:]))
print(f'Number of windowed increases: {count}')
