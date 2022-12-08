use std::error::Error;
use std::fmt::Display;
use std::fs::File;
use std::io::BufRead;
use std::io::BufReader;

fn main() -> Result<(), Box<dyn Error>> {
    let f = File::open("input.txt")?;
    let mut reader = BufReader::new(f);
    let prog = parse(&mut reader)?;

    let p1 = run(&prog, apply1);
    println!("Final position, part 1: {}", p1);
    let p2 = run(&prog, apply2);
    println!("Final position: part 2: {}", p2);
    Ok(())
}

enum Op {
    Forward(i32),
    Up(i32),
    Down(i32),
}

impl Display for Op {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::result::Result<(), ::std::fmt::Error> {
        match *self {
            Op::Forward(n) => f.write_fmt(format_args!("Forward({})", n)),
            Op::Up(n) => f.write_fmt(format_args!("Up({})", n)),
            Op::Down(n) => f.write_fmt(format_args!("Down({})", n)),
        }
    }
}

fn parse<R: BufRead>(r: &mut R) -> Result<Vec<Op>, Box<dyn Error>> {
    let mut prog = Vec::new();
    for line in r.lines() {
        let l = line?;
        let i = l.find(' ').ok_or("no space in line")?;
        let count = l[i + 1..].parse::<i32>()?;
        let cmd = match &l[..i] {
            "forward" => Ok(Op::Forward(count)),
            "up" => Ok(Op::Up(count)),
            "down" => Ok(Op::Down(count)),
            _ => Err("invalid command"),
        }?;
        prog.push(cmd);
    }
    Ok(prog)
}

struct Pos {
    horizontal: i32,
    depth: i32,
    aim: i32,
}

impl Pos {
    fn new() -> Pos {
        Pos {
            horizontal: 0,
            depth: 0,
            aim: 0,
        }
    }
}

impl Display for Pos {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::result::Result<(), ::std::fmt::Error> {
        f.write_fmt(format_args!(
            "horizontal={}, depth={}, product={}",
            self.horizontal,
            self.depth,
            self.horizontal * self.depth
        ))
    }
}

fn run(prog: &Vec<Op>, apply: fn(Pos, &Op) -> Pos) -> Pos {
    prog.iter().fold(Pos::new(), apply)
}

fn apply1(p: Pos, cmd: &Op) -> Pos {
    match cmd {
        Op::Forward(n) => Pos {
            horizontal: p.horizontal + n,
            depth: p.depth,
            aim: p.aim,
        },
        Op::Down(n) => Pos {
            horizontal: p.horizontal,
            depth: p.depth + n,
            aim: p.aim,
        },
        Op::Up(n) => Pos {
            horizontal: p.horizontal,
            depth: p.depth - n,
            aim: p.aim,
        },
    }
}

fn apply2(p: Pos, cmd: &Op) -> Pos {
    match cmd {
        Op::Forward(n) => Pos {
            horizontal: p.horizontal + n,
            depth: p.depth + n * p.aim,
            aim: p.aim,
        },
        Op::Down(n) => Pos {
            horizontal: p.horizontal,
            depth: p.depth,
            aim: p.aim + n,
        },
        Op::Up(n) => Pos {
            horizontal: p.horizontal,
            depth: p.depth,
            aim: p.aim - n,
        },
    }
}
