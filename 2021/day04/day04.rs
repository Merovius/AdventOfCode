use std::error::Error;
use std::fmt::{self, Debug, Display};
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::str::FromStr;

const ROWS: usize = 5;
const COLS: usize = 5;

fn main() -> Result<(), Box<dyn Error>> {
    let f = File::open("input.txt")?;
    let lines: Vec<String> = BufReader::new(f).lines().collect::<Result<_, _>>().unwrap();
    let data = parse(&lines);

    let winners = play(&data);

    println!(
        "First winner's score: {}",
        winners.first().unwrap().score.unwrap()
    );
    println!(
        "Last winner's score: {}",
        winners.last().unwrap().score.unwrap()
    );

    Ok(())
}

fn parse(lines: &Vec<String>) -> Data {
    if lines.len() % (ROWS + 1) != 1 {
        panic!("wrong number of lines")
    }
    let nums = lines[0]
        .split(",")
        .map(|v| u8::from_str(v).unwrap())
        .collect();

    let mut boards = Vec::new();
    for i in (2..lines.len()).step_by(ROWS + 1) {
        let b = parse_board(&lines[i..i + ROWS]);
        boards.push(b);
    }

    Data {
        nums: nums,
        boards: boards,
    }
}

fn parse_board(lines: &[String]) -> Board<u8> {
    let mut board = Board::new(0u8);
    for r in 0..ROWS {
        let vals: Vec<u8> = lines[r]
            .split_whitespace()
            .map(|v| u8::from_str(v).unwrap())
            .collect();
        if vals.len() != COLS {
            panic!("wrong number of columns");
        }
        for c in 0..COLS {
            board.els[r][c] = vals[c];
        }
    }
    board
}

#[derive(Debug)]
struct Board<T> {
    els: [[T; COLS]; ROWS],
}

impl<T: Display> Display for Board<T> {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        for row in &self.els {
            for v in row {
                f.write_fmt(format_args!("{:>2} ", v))?;
            }
            f.write_str("\n")?;
        }
        Ok(())
    }
}

impl<T: Copy> Board<T> {
    fn new(val: T) -> Board<T> {
        return Board {
            els: [[val; COLS]; ROWS],
        };
    }
}

#[derive(Debug)]
struct ScoredBoard<'a> {
    board: &'a Board<u8>,
    marked: Board<bool>,
    score: Option<u64>,
}

impl<'a> Display for ScoredBoard<'a> {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        for r in 0..ROWS {
            for c in 0..COLS {
                if self.marked.els[r][c] {
                    f.write_fmt(format_args!("{:>2}* ", self.board.els[r][c]))?;
                } else {
                    f.write_fmt(format_args!("{:>2}  ", self.board.els[r][c]))?;
                }
            }
            f.write_str("\n")?;
        }
        match self.score {
            Some(score) => write!(f, "Score: {}", score),
            None => Ok(()),
        }
    }
}

impl<'a> ScoredBoard<'a> {
    fn new(board: &Board<u8>) -> ScoredBoard {
        ScoredBoard {
            board: board,
            marked: Board::new(false),
            score: None,
        }
    }

    // mark marks any instance of v on the board.
    fn mark(&mut self, v: u8) {
        match self.score {
            Some(_) => return,
            None => (),
        }
        for r in 0..ROWS {
            for c in 0..COLS {
                if self.board.els[r][c] != v {
                    continue;
                }
                self.marked.els[r][c] = true;
                self.calculate_score(r, c);
            }
        }
    }

    // calculate_score checks if marking [row,col] caused the board to win and calculate its score.
    fn calculate_score(&mut self, row: usize, col: usize) {
        let mut won_row = true;
        let mut won_col = true;
        for r in 0..ROWS {
            won_row &= self.marked.els[r][col];
        }
        for c in 0..COLS {
            won_col &= self.marked.els[row][c];
        }
        if !(won_row || won_col) {
            return;
        }
        let mut score = 0;
        for r in 0..ROWS {
            for c in 0..COLS {
                if !self.marked.els[r][c] {
                    score += self.board.els[r][c] as u64;
                }
            }
        }
        self.score = Some(score * self.board.els[row][col] as u64)
    }
}

struct Data {
    nums: Vec<u8>,
    boards: Vec<Board<u8>>,
}

fn play(data: &Data) -> Vec<ScoredBoard> {
    let mut winners = Vec::new();
    let mut boards: Vec<ScoredBoard> = data.boards.iter().map(ScoredBoard::new).collect();

    for v in &data.nums {
        let mut i = 0;
        while i < boards.len() {
            {
                let b = boards.get_mut(i).unwrap();
                b.mark(*v);
                if b.score.is_none() {
                    i += 1;
                    continue;
                }
            }
            winners.push(boards.swap_remove(i));
        }
    }

    winners
}
