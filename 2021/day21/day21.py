from collections import namedtuple

# example
#pos = (4, 8)
# input
pos = (4, 7)

class DeterministicDie:
    def __init__(self):
        self._roll = 0

    def roll(self):
        self._roll += 1
        if self._roll > 100:
            self._roll = 1
        return self._roll

Result = namedtuple("Result", ["pos", "score", "rolls"])

def play_deterministic(pos):
    die = DeterministicDie()
    n = 0
    score = [0,0]
    pos = list(pos)
    while score[0]<1000 and score[1]<1000:
        r = (die.roll(),die.roll(),die.roll())
        p = (n//3)%2
        pos[p] += r[0]+r[1]+r[2]
        while pos[p] > 10:
            pos[p] -= 10
        score[p] += pos[p]
        n += 3
    return Result(pos=tuple(pos), score=tuple(score), rolls=n)

res = play_deterministic(pos)
print(f"Scores: {res.score}")
print(f"Die rolls: {res.rolls}")
print(f"Product: {res.rolls*min(res.score[0], res.score[1])}")

State = namedtuple("State", ["pos", "score", "player"])
mem = {}
def play_dirac(pos):
    mem = {}
    def play_rec(state):
        if state in mem:
            return mem[state]

        nwin = [0,0]
        # iterate through all possible totals of three Dirac Dice rolls
        for i in range(27):
            total = (i%3)+(i//3)%3+(i//9)%3+3
            p = (state.player + 1)%2

            pos = list(state.pos)
            pos[p] += total
            while pos[p] > 10:
                pos[p] -= 10
            pos = tuple(pos)

            score = list(state.score)
            score[p] += pos[p]
            score = tuple(score)

            if score[p] >= 21:
                nwin[p] += 1
                continue
            nw = play_rec(State(pos, score, p))
            nwin[0] += nw[0]
            nwin[1] += nw[1]
        mem[state] = nwin
        return nwin
    return play_rec(State(pos, (0,0), 1))

print(f"Player wins in Dirac Dice: {play_dirac(pos)}")
