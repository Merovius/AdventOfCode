use std::collections::HashSet;
use std::fs::File;
use std::io::{BufRead, BufReader};

use clap::{App, Arg};
use simple_error::SimpleError;

type Error = Box<dyn std::error::Error>;

fn main() -> Result<(), Error> {
    let matches = App::new("day13")
        .arg(Arg::with_name("file").index(1).required(true))
        .get_matches();

    let (mut points, folds) = parse(matches.value_of("file").unwrap())?;
    println!(
        "After first fold, there are {} points.",
        set_apply(&points, folds[0]).len()
    );
    for f in folds {
        points = set_apply(&points, f);
    }
    println!("After all folds, the paper looks like:");
    print_set(&points);
    Ok(())
}

fn set_apply(points: &HashSet<Point>, f: Fold) -> HashSet<Point> {
    points.into_iter().map(|p| f.apply(*p)).collect()
}

fn print_set(points: &HashSet<Point>) {
    let minx = points.into_iter().map(|p| p.0).min().unwrap();
    let maxx = points.into_iter().map(|p| p.0).max().unwrap();
    let miny = points.into_iter().map(|p| p.1).min().unwrap();
    let maxy = points.into_iter().map(|p| p.1).max().unwrap();
    for y in miny..=maxy {
        for x in minx..=maxx {
            if points.contains(&Point(x, y)) {
                print!("•");
            } else {
                print!(" ");
            }
        }
        println!();
    }
}

#[derive(Clone, Copy, Debug)]
enum Fold {
    X(i64),
    Y(i64),
}

impl std::str::FromStr for Fold {
    type Err = Error;
    fn from_str(s: &str) -> Result<Fold, Error> {
        match s.strip_prefix("fold along ") {
            None => Err(Box::new(SimpleError::new("not a valid fold instruction"))),
            Some(s) => {
                let result: Vec<_> = s.split("=").collect();
                match result[..] {
                    [s, v] => match s {
                        "x" => Ok(Fold::X(i64::from_str(v)?)),
                        "y" => Ok(Fold::Y(i64::from_str(v)?)),
                        _ => Err(Box::new(SimpleError::new("invalid fold instruction"))),
                    },
                    _ => Err(Box::new(SimpleError::new("invalid fold instruction"))),
                }
            }
        }
    }
}

impl Fold {
    fn apply(&self, p: Point) -> Point {
        match *self {
            Fold::X(v) => Point(if p.0 < v { p.0 } else { 2 * v - p.0 }, p.1),
            Fold::Y(v) => Point(p.0, if p.1 < v { p.1 } else { 2 * v - p.1 }),
        }
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
struct Point(i64, i64);

impl std::str::FromStr for Point {
    type Err = Error;
    fn from_str(s: &str) -> Result<Point, Error> {
        let result: Result<Vec<_>, _> = s.split(",").map(i64::from_str).collect();
        match result {
            Err(e) => Err(Box::new(e)),
            Ok(vec) => match vec[..] {
                [x, y] => Ok(Point(x, y)),
                _ => Err(Box::new(SimpleError::new("wrong number of coordinates"))),
            },
        }
    }
}

fn parse(name: &str) -> Result<(HashSet<Point>, Vec<Fold>), Error> {
    let mut points = HashSet::new();
    let mut folds = Vec::new();
    let mut parse_folds = false;
    for l in BufReader::new(File::open(name)?).lines() {
        let l = l?;
        if l.is_empty() {
            if parse_folds {
                return Err(Box::new(SimpleError::new("extra empty line in input")));
            }
            parse_folds = true;
            continue;
        }
        if parse_folds {
            folds.push(l.parse::<Fold>()?);
        } else {
            points.insert(l.parse::<Point>()?);
        }
    }
    Ok((points, folds))
}
