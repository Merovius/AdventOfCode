from statistics import median, mean
import math

with open('input.txt', 'r') as file:
    data = list(map(int, file.read().strip().split(',')))
#data = [16,1,2,0,4,2,7,1,2,14]

def fuel_cost1(pos):
    return round(sum(map(lambda x: abs(x-pos), data)))

def fuel_cost2(pos):
    return round(sum(map(lambda n: n*(n+1)//2, map(lambda x: abs(x-pos), data))))

_median, _mean = round(median(data)), mean(data)
print(f"Minimum fuel cost part 1: {fuel_cost1(_median)} at {_median}")
fc1 = fuel_cost2(math.floor(_mean))
fc2 = fuel_cost2(math.ceil(_mean))
if fc1 < fc2:
    print(f"Minimum fuel cost part 2: {fc1} at {math.floor(_mean)}")
else:
    print(f"Minimum fuel cost part 2: {fc2} at {math.ceil(_mean)}")

