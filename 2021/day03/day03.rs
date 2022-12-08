use std::cmp::min;
use std::error::Error;
use std::fs::File;
use std::io::BufRead;
use std::io::BufReader;

fn main() -> Result<(), Box<dyn Error>> {
    let f = File::open("input.txt")?;
    let mut reader = BufReader::new(f);
    let data = parse(&mut reader)?;
    let (gamma, epsilon) = gamma_epsilon(&data);
    println!(
        "γ: {:012b}={}, ε: {:012b}={}, γ•ε: {}",
        gamma,
        gamma,
        epsilon,
        epsilon,
        gamma * epsilon
    );

    let o2 = filter(&data, true);
    let co2 = filter(&data, false);
    println!(
        "O₂: {:12b}={}, CO₂: {:12b}={}, product: {}",
        o2,
        o2,
        co2,
        co2,
        o2 * co2
    );

    Ok(())
}

fn parse<R: BufRead>(r: &mut R) -> Result<Data, Box<dyn Error>> {
    let mut n = 16;
    let mut result = Vec::new();
    for line in r.lines() {
        let v = u16::from_str_radix(&line?, 2)?;
        n = min(n, v.leading_zeros());
        result.push(v);
    }
    Ok(Data {
        n: 16 - n as usize,
        nums: result,
    })
}

struct Data {
    n: usize,
    nums: Vec<u16>,
}

fn bit_is_commonly_set(nums: &Vec<u16>, i: usize) -> bool {
    let mut ones = 0;
    let mut zeros = 0;
    for v in nums {
        if (v >> i) & 1 == 0 {
            zeros += 1;
        } else {
            ones += 1;
        }
    }
    return ones >= zeros;
}

fn gamma_epsilon(data: &Data) -> (u64, u64) {
    let mut gamma = 0;
    let mut epsilon = 0;
    for i in (0..data.n).rev() {
        gamma <<= 1;
        epsilon <<= 1;
        if bit_is_commonly_set(&data.nums, i) {
            gamma += 1;
        } else {
            epsilon += 1;
        }
    }
    (gamma, epsilon)
}

fn filter(data: &Data, common: bool) -> u64 {
    let mut filtered = data.nums.to_vec();
    for i in (0..data.n).rev() {
        let want_common = common == bit_is_commonly_set(&filtered, i);
        let mut j = 0;
        for k in 0..filtered.len() {
            let v = filtered[k];
            let is_set = (v >> i) & 1 != 0;
            if is_set == want_common {
                filtered[j] = v;
                j += 1;
            }
        }
        filtered.truncate(j);
        if filtered.len() == 1 {
            break;
        }
    }
    if filtered.len() > 1 {
        panic!("filtered set has more than one element");
    }
    return filtered[0] as u64;
}
