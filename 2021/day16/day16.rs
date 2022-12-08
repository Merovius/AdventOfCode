use bitvec::prelude::*;
use hex;

type Error = Box<dyn std::error::Error>;

fn main() -> Result<(), Error> {
    let input = include_str!("input.txt");
    let data = hex::decode(input)?;
    let bits = BitSlice::<Msb0, _>::from_slice(&data)?;
    let p = Packet::parse(&mut bits.iter());

    println!("Version sum: {:?}", p.version_sum());
    println!("Evaluated: {:?}", p.eval());

    Ok(())
}

type Bits<'a> = bitvec::slice::Iter<'a, Msb0, u8>;

#[derive(Debug)]
enum Packet {
    Literal(u8, u64),
    Operator(u8, Op, Vec<Packet>),
}

#[derive(Clone, Copy, Debug)]
enum Op {
    Add,
    Mul,
    Max,
    Min,
    Gt,
    Lt,
    Eq,
}

impl Op {
    fn from_u8(v: u8) -> Op {
        match v {
            0 => Op::Add,
            1 => Op::Mul,
            2 => Op::Min,
            3 => Op::Max,
            5 => Op::Gt,
            6 => Op::Lt,
            7 => Op::Eq,
            _ => panic!("invalid op"),
        }
    }

    fn apply(self, vs: Vec<u64>) -> u64 {
        match self {
            Op::Add => vs.iter().sum(),
            Op::Mul => vs.iter().product(),
            Op::Min => *vs.iter().min().unwrap(),
            Op::Max => *vs.iter().max().unwrap(),
            Op::Gt => {
                if vs[0] > vs[1] {
                    1
                } else {
                    0
                }
            }
            Op::Lt => {
                if vs[0] < vs[1] {
                    1
                } else {
                    0
                }
            }
            Op::Eq => {
                if vs[0] == vs[1] {
                    1
                } else {
                    0
                }
            }
        }
    }
}

impl Packet {
    fn parse(bits: &mut Bits) -> Packet {
        let ver = Self::_int(bits, 3) as u8;
        let id = Self::_int(bits, 3);
        if id == 4 {
            return Self::_literal(bits, ver as u8);
        }
        return Self::_operator(bits, ver, Op::from_u8(id as u8));
    }

    fn _int(bits: &mut Bits, n: usize) -> u64 {
        let mut val = 0u64;
        for _ in 0..n {
            val <<= 1;
            if *bits.next().unwrap() {
                val |= 1;
            }
        }
        val
    }

    fn _literal(bits: &mut Bits, ver: u8) -> Packet {
        let mut val = 0u64;
        let mut cont = true;
        while cont {
            cont = *bits.next().unwrap();
            val <<= 4;
            val |= Self::_int(bits, 4);
        }
        Packet::Literal(ver, val)
    }

    fn _operator(bits: &mut Bits, ver: u8, id: Op) -> Packet {
        let mut sub = Vec::new();
        if *bits.next().unwrap() {
            let m = Self::_int(bits, 11);
            for _ in 0..m {
                sub.push(Self::parse(bits));
            }
        } else {
            let k = Self::_int(bits, 15) as usize;
            let m = bits.len() - k;
            while bits.len() > m {
                sub.push(Self::parse(bits));
            }
        }
        Packet::Operator(ver, id, sub)
    }

    fn version_sum(&self) -> u64 {
        match self {
            Packet::Literal(ver, _) => *ver as u64,
            Packet::Operator(ver, _, sub) => {
                sub.iter().fold(*ver as u64, |v, p| v + p.version_sum())
            }
        }
    }

    fn eval(&self) -> u64 {
        match self {
            Packet::Literal(_, val) => *val,
            Packet::Operator(_, op, sub) => op.apply(sub.iter().map(|p| p.eval()).collect()),
        }
    }
}
