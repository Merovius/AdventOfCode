use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};

fn main() -> Result<(), Box<dyn Error>> {
    let mut fish = parse("input.txt")?;
    for i in 0..256 {
        if i == 80 {
            println!("After 80 days there are {} fish", fish.iter().sum::<u64>());
        }
        fish = step(&fish);
    }
    println!("After 256 days there are {} fish", fish.iter().sum::<u64>());

    Ok(())
}

fn parse(name: &str) -> Result<Vec<u64>, Box<dyn Error>> {
    let mut result = vec![0u64; 9];

    let mut line = String::new();
    BufReader::new(File::open(name)?).read_line(&mut line)?;
    line.trim()
        .split(",")
        .filter_map(|s| s.parse::<u64>().ok())
        .for_each(|v| {
            result[v as usize] += 1;
        });
    Ok(result)
}

fn step(fish: &Vec<u64>) -> Vec<u64> {
    let mut new = vec![0u64; 9];
    for (d, n) in fish.iter().enumerate() {
        match d {
            0 => {
                new[8] += n;
                new[6] += n;
            }
            _ => {
                new[d - 1] += n;
            }
        }
    }
    new
}
