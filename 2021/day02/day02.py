from functools import reduce

def readfile(name):
    with open(name, 'r') as file:
        return file.readlines()

data = readfile('input.txt')
#data = ["forward 5", "down 5", "forward 8", "up 3", "down 8", "forward 2"]

def apply_commands_part1(cmds):
    def apply_command(pos, cmd):
        [op, count] = cmd.split()
        if op == 'forward':
            return (pos[0]+int(count), pos[1])
        if op == 'down':
            return (pos[0], pos[1]+int(count))
        if op == 'up':
            return (pos[0], pos[1]-int(count))
        raise ValueError(f'invalid op {op}')
    return reduce(apply_command, cmds, (0, 0))

def apply_commands_part2(cmds):
    def apply_command(pos, cmd):
        [op, count] = cmd.split()
        if op == 'forward':
            return (pos[0]+int(count), pos[1]+pos[2]*int(count), pos[2])
        if op == 'down':
            return (pos[0], pos[1], pos[2]+int(count))
        if op == 'up':
            return (pos[0], pos[1], pos[2]-int(count))
        raise ValueError(f'invalid op {op}')
    return reduce(apply_command, cmds, (0, 0, 0))

pos = apply_commands_part1(data)
print(f'Final position, part 1: horizontal={pos[0]}, depth={pos[1]}, product={pos[0]*pos[1]}')
pos = apply_commands_part2(data)
print(f'Final position, part 2: horizontal={pos[0]}, depth={pos[1]}, product={pos[0]*pos[1]}')
