use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};

fn main() -> Result<(), Box<dyn Error>> {
    let file = "input.txt";

    let mut corrupt_score = 0u64;
    let mut incomplete_scores = Vec::new();
    BufReader::new(File::open(file)?)
        .lines()
        .for_each(|l| match score(&l.unwrap()) {
            Score::Corrupt(s) => corrupt_score += s,
            Score::Incomplete(s) => incomplete_scores.push(s),
        });
    incomplete_scores.sort();
    println!("Score for corrupted lines: {}", corrupt_score);
    println!("Score for incomplete lines: {}", incomplete_scores[incomplete_scores.len()/2]);
    Ok(())
}

fn score(line: &str) -> Score {
    let mut stack = Vec::new();
    for c in line.chars() {
        match c {
            '(' | '[' | '{' | '<' => {
                stack.push(c);
                continue;
            }
            _ => (),
        }
        if c == closing_char(stack.pop().unwrap()) {
            continue;
        }
        match c {
            ')' => return Score::Corrupt(3),
            ']' => return Score::Corrupt(57),
            '}' => return Score::Corrupt(1197),
            '>' => return Score::Corrupt(25137),
            _ => panic!("invalid character {}", c),
        }
    }
    let mut score = 0u64;
    while stack.len() > 0 {
        let c = stack.pop().unwrap();
        score = score*5 + match c {
            '(' => 1,
            '[' => 2,
            '{' => 3,
            '<' => 4,
            _ => panic!("invalid character {}", c),
        }
    }
    Score::Incomplete(score)
}

fn closing_char(c: char) -> char {
    match c {
        '(' => ')',
        '[' => ']',
        '{' => '}',
        '<' => '>',
        _ => panic!("invalid character {}", c),
    }
}

#[derive(Debug)]
enum Score {
    Corrupt(u64),
    Incomplete(u64),
}
