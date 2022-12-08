// Solution is not complete

use simple_error::SimpleError;

type Error = Box<dyn std::error::Error>;

#[derive(Debug)]
enum Num {
    Reg(u64),
    Pair(Box<Num>, Box<Num>),
}

impl Num {
    fn from_str(s: &str) -> Result<Num, Error> {
        Num::from_json(&serde_json::from_str(s).unwrap())
    }

    fn from_json(v: &serde_json::Value) -> Result<Num, Error> {
        match v {
            serde_json::Value::Number(n) => match n.as_u64() {
                Some(n) => Ok(Num::Reg(n)),
                None => Err(Box::new(SimpleError::new("number can't be represented as u64"))),
            },
            serde_json::Value::Array(v) => match &v[..] {
                [l,r] => Ok(Num::Pair(Box::new(Num::from_json(l)?), Box::new(Num::from_json(r)?))),
                _ => Err(Box::new(SimpleError::new("array has wrong number of elements"))),
            }
            _ => Err(Box::new(SimpleError::new("invalid json value"))),
        }
    }

    fn magnitude(&self) -> u64 {
        match self {
            Num::Reg(n) => *n,
            Num::Pair(l, r) => 3*l.magnitude()+2*r.magnitude(),
        }
    }

    fn depth(&self) -> u64 {
        match self {
            Num::Reg(_) => 0,
            Num::Pair(l, r) => u64::max(l.depth(), r.depth())+1,
        }
    }

    fn reduce(self) -> Num {
        let mut n = self;
        loop {
            if n.depth() > 4 {
                n = n.splode();
                continue
            }
            if n.splits() {
                n = n.split();
                continue
            }
            return n
        }
    }

    fn splits(&self) -> bool {
        match self {
            Num::Reg(n) => *n >= 10,
            Num::Pair(l, r) => l.splits() || r.splits(),
        }
    }

    fn split(self) -> Num {
        match self {
            Num::Reg(n) => {
                let l = n/2;
                let r = n-l;
                Num::Pair(Box::new(Num::Reg(l)), Box::new(Num::Reg(r)))
            },
            Num::Pair(l, r) => {
                if l.splits() {
                    return Num::Pair(Box::new(l.split()), r);
                }
                if r.splits() {
                    return Num::Pair(l, Box::new(r.split()));
                }
                panic!("split on non-splitting number");
            },
        }
    }

    fn splode(self) -> Num {
        let (n, _, _) = self.splode_rec(0);
        n
    }

    fn splode_rec(self, lvl: u64) -> (Num, Option<u64>, Option<u64>) {
        match self {
            Num::Reg(n) => (Num::Reg(n), None, None),
            Num::Pair(l, r) => {
                if lvl == 4 {
                    return (Num::Reg(0), l.reg(), r.reg());
                }
                if lvl+l.depth() >= 4 {
                    return match l.splode_rec(lvl+1) {
                        (n, lv, Some(rv)) => (Num::Pair(Box::new(n), Box::new(r.add_left(rv))), lv, None),
                        (n, lv, None) => (Num::Pair(Box::new(n), r), lv, None),
                    }
                }
                if lvl+r.depth() >= 4 {
                    return match r.splode_rec(lvl+1) {
                        (n, Some(lv), rv) => (Num::Pair(Box::new(l.add_right(lv)), Box::new(n)), None, rv),
                        (n, None, rv) => (Num::Pair(l, Box::new(n)), None, rv),
                    }
                }
                panic!("splode_rec on non-sploding number");
            }
        }
    }

    fn reg(&self) -> Option<u64> {
        match self {
            Num::Reg(n) => Some(*n),
            _ => None,
        }
    }

    fn add_left(self, v: u64) -> Num {
        match self {
            Num::Reg(n) => Num::Reg(n+v),
            Num::Pair(l, r) => Num::Pair(Box::new(l.add_left(v)), r),
        }
    }

    fn add_right(self, v: u64) -> Num {
        match self {
            Num::Reg(n) => Num::Reg(n+v),
            Num::Pair(l, r) => Num::Pair(l, Box::new(r.add_right(v))),
        }
    }
}

fn main() -> Result<(), Error> {
    Ok(())
}
