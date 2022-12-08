use std::collections::HashMap;
use std::fs::File;
use std::hash::Hash;
use std::io::{BufRead,BufReader};

use clap::{App,Arg};
use simple_error::SimpleError;

type Error = Box<dyn std::error::Error>;

fn main() -> Result<(), Error> {
    let matches = App::new("day14")
        .arg(Arg::with_name("file").index(1).required(true))
        .get_matches();
    let (tpl, rules) = parse(matches.value_of("file").unwrap())?;

    let mut counts = Counts::new(&tpl, &rules);
    for _ in 0..10 {
        counts.step();
    }
    println!("After 10 steps, the difference between most and least common element is {}", counts.score());
    for _ in 0..30 {
        counts.step();
    }
    println!("After 40 steps, the difference between most and least common element is {}", counts.score());

    Ok(())
}

fn parse(name: &str) -> Result<(String, Rules), Error> {
    let mut lines = BufReader::new(File::open(name)?).lines();
    let tpl = match lines.next() {
        None => return Err(Box::new(SimpleError::new("empty file"))),
        Some(l) => l?,
    };
    match lines.next() {
        None => return Err(Box::new(SimpleError::new("no rules in file"))),
        Some(l) => if !l?.is_empty() { return Err(Box::new(SimpleError::new("missing newline"))) },
    };
    let mut rules = Rules::new();
    for l in lines {
        let l = l?;
        let sp: Vec<_> = l.split(" -> ").collect();
        if sp.len() != 2 {
            return Err(Box::new(SimpleError::new("invalid rule")));
        }
        let pair: Vec<_> = sp[0].chars().collect();
        let repl: Vec<_> = sp[1].chars().collect();
        if pair.len() != 2 || repl.len() != 1 {
            return Err(Box::new(SimpleError::new("invalid rule")));
        }
        match rules.insert((pair[0], pair[1]), repl[0]) {
            None => (),
            Some(r) => return Err(Box::new(SimpleError::new(format!("duplicate rule {} maps to {} and {}", sp[0], r, sp[1])))),
        }
    }
    Ok((tpl, rules))
}

type Rules = HashMap<(char, char), char>;

struct Counts<'a> {
    last: char,
    counts: HashMap<(char, char), usize>,
    rules: &'a Rules,
}

fn incr<K: Eq+Hash+Copy>(m: &mut HashMap<K, usize>, k: &K, v: usize) {
    match m.get_mut(k) {
        Some(x) => *x += v,
        None => { m.insert(*k, v); () },
    };
}

impl<'a> Counts<'a> {
    fn new(tpl: &str, rules: &'a Rules) -> Counts<'a> {
        let last = tpl.chars().last().unwrap();
        let counts = tpl.chars().zip(tpl.chars().skip(1)).fold(HashMap::new(), |mut counts, pair| {
            match counts.get_mut(&pair) {
                None => {
                    counts.insert(pair, 1);
                    ()
                },
                Some(v) => *v += 1,
            };
            counts
        });
        Counts{
            last: last,
            counts: counts,
            rules: rules,
        }
    }

    fn step(&mut self) {
        let mut next = HashMap::new();
        for (k, v) in self.counts.iter() {
            match self.rules.get(&k) {
                Some(r) => {
                    incr(&mut next, &(k.0, *r), *v);
                    incr(&mut next, &(*r, k.1), *v);
                }
                None => incr(&mut next, &k, *v),
            }
        }
        self.counts = next;
    }

    fn score(&self) -> usize {
        let mut letter_counts = HashMap::new();
        for (k, v) in self.counts.iter() {
            incr(&mut letter_counts, &k.0, *v);
        }
        incr(&mut letter_counts, &self.last, 1);
        let max = letter_counts.values().max().unwrap();
        let min = letter_counts.values().min().unwrap();
        return max-min;
    }
}
