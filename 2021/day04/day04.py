from collections import namedtuple

def readfile(name):
    with open(name, 'r') as file:
        return file.readlines()

class Board:
    def __init__(self, board):
        self._board = board
        self._marked = []
        for i in range(5):
            self._marked.append([False]*5)
        self.won = False
        self.score = 0

    def mark(self, n):
        for row in range(5):
            for col in range(5):
                if self._board[row][col] != n:
                    continue
                self._marked[row][col] = True
                self._check_win(row, col)
                return

    def _check_win(self, row, col):
        if self.won:
            return

        row_win, col_win = True, True
        for i in range(5):
            row_win = row_win and self._marked[row][i]
            col_win = col_win and self._marked[i][col]
        self.won = row_win or col_win

        if not self.won:
            return

        for r in range(5):
            for c in range(5):
                if not self._marked[r][c]:
                    self.score += self._board[r][c]
        self.score *= self._board[row][col]

    def __str__(self):
        return "\n".join([" ".join([f"{n: >2}" for n in row]) for row in self._board])

Data = namedtuple('Data', ['draw', 'boards'])

def read_data(name):
    lines = readfile(name)
    draw = [int(l) for l in lines[0].split(',')]
    boards = []

    lines = lines[1:]
    for i in range(0, len(lines), 6):
        block = lines[i+1:i+6]
        boards.append(Board([[int(x) for x in l.split()] for l in block]))

    return Data(draw, boards)

def find_winners(data):
    winners = []
    for n in data.draw:
        for b in data.boards:
            if b.won:
                continue
            b.mark(n)
            if b.won:
                winners.append(b)
    return winners


data = read_data('input.txt')

winners = find_winners(data)
print(f"Winning board's score: {winners[0].score}")
print(f"Last winner's score: {winners[-1].score}")
