use std::cmp::Eq;
use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::iter::{IntoIterator, Iterator};
use std::ops::{Add, AddAssign, Index, IndexMut, Sub, SubAssign};
use std::str::FromStr;

use array2d::Array2D;
use simple_error::SimpleError;

fn main() -> Result<(), Box<dyn Error>> {
    let lines = BufReader::new(File::open("day05/example.txt")?)
        .lines()
        .map(|s| s?.parse::<Line>())
        .collect::<Result<Vec<_>, _>>()?;
    let max = (&lines).into_iter().flatten().reduce(Point::max).unwrap() + Point(1, 1);
    let mut grid = Array2D::filled_with(0, max.0 as usize, max.1 as usize);

    for l in (&lines).into_iter().filter(|l| l.is_aligned()) {
        for p in l {
            grid[p] += 1;
        }
    }

    let mut total = 0u64;
    for row in grid.rows_iter() {
        for v in row {
            if *v > 1 {
                total += 1;
            }
        }
    }
    println!("Total number of intersections: {}", total);

    for l in (&lines).into_iter().filter(|l| !l.is_aligned()) {
        for p in l {
            grid[p] += 1;
        }
    }
    total = 0;
    for row in grid.rows_iter() {
        for v in row {
            if *v > 1 {
                total += 1;
            }
        }
    }
    println!("Total number of intersections: {}", total);

    Ok(())
}

fn print_grid(grid: &Array2D<i32>) {
    for row in 0..grid.num_rows() {
        for col in 0..grid.num_columns() {
            match grid.get(row, col).unwrap() {
                0 => print!("."),
                n => print!("{}", n),
            }
        }
        println!();
    }
}

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
struct Point(isize, isize);

impl Point {
    fn normalize(self) -> Point {
        Point(self.0.signum(), self.1.signum())
    }

    fn max(self, other: Point) -> Point {
        Point(isize::max(self.0, other.0), isize::max(self.1, other.1))
    }
}

impl FromStr for Point {
    type Err = Box<dyn Error>;
    fn from_str(s: &str) -> Result<Point, Self::Err> {
        let coords: Result<Vec<_>, _> = s.split(",").map(|s| s.parse::<usize>()).collect();
        match coords {
            Err(e) => Err(Box::new(SimpleError::from(e))),
            Ok(v) => match v[..] {
                [a, b] => Ok(Point(a as isize, b as isize)),
                _ => Err(Box::new(SimpleError::new("wrong number of coordinates"))),
            },
        }
    }
}

impl Add for Point {
    type Output = Point;

    fn add(self, other: Point) -> Point {
        Point(self.0 + other.0, self.1 + other.1)
    }
}

impl AddAssign for Point {
    fn add_assign(&mut self, rhs: Point) {
        *self = *self + rhs;
    }
}

impl Sub for Point {
    type Output = Point;

    fn sub(self, other: Point) -> Point {
        Point(self.0 - other.0, self.1 - other.1)
    }
}

impl SubAssign for Point {
    fn sub_assign(&mut self, rhs: Point) {
        *self = *self - rhs;
    }
}

impl<T: Clone> Index<Point> for Array2D<T> {
    type Output = T;

    fn index(&self, p: Point) -> &T {
        self.get(p.0 as usize, p.1 as usize).unwrap()
    }
}

impl<T: Clone> IndexMut<Point> for Array2D<T> {
    fn index_mut(&mut self, p: Point) -> &mut T {
        self.get_mut(p.0 as usize, p.1 as usize).unwrap()
    }
}

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
struct Line {
    from: Point,
    to: Point,
}

impl Line {
    fn make_safely(from: Point, to: Point) -> Option<Line> {
        let delta = to - from;
        if delta.0 != 0 && delta.1 != 0 && delta.0.abs() != delta.1.abs() {
            None
        } else {
            Some(Line { from: from, to: to })
        }
    }

    fn is_aligned(self) -> bool {
        self.from.0 == self.to.0 || self.from.1 == self.to.1
    }
}

impl FromStr for Line {
    type Err = Box<dyn Error>;
    fn from_str(s: &str) -> Result<Line, Self::Err> {
        let points: Result<Vec<_>, _> = s.split(" -> ").map(|s| s.parse::<Point>()).collect();
        match points {
            Err(e) => Err(e),
            Ok(v) => match v[..] {
                [from, to] => match Line::make_safely(from, to) {
                    None => Err(Box::new(SimpleError::new("line is not straight"))),
                    Some(l) => Ok(l),
                },
                _ => Err(Box::new(SimpleError::new("wrong number of coordinates"))),
            },
        }
    }
}

impl IntoIterator for Line {
    type Item = Point;
    type IntoIter = LineIter;

    fn into_iter(self) -> LineIter {
        let step = (self.to - self.from).normalize();
        LineIter {
            pt: self.from,
            end: self.to + step,
            step: step,
        }
    }
}

impl IntoIterator for &Line {
    type Item = Point;
    type IntoIter = LineIter;

    fn into_iter(self) -> LineIter {
        let step = (self.to - self.from).normalize();
        LineIter {
            pt: self.from,
            end: self.to + step,
            step: step,
        }
    }
}

struct LineIter {
    pt: Point,
    step: Point,
    end: Point,
}

impl Iterator for LineIter {
    type Item = Point;

    fn next(&mut self) -> Option<Point> {
        if self.pt == self.end {
            None
        } else {
            let p = self.pt;
            self.pt += self.step;
            Some(p)
        }
    }
}
