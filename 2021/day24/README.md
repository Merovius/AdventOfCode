The code reads input.txt and transforms it into a graph.
It then applies a bunch of graph transformations, to simplify it.
At this point, you can look at the graph:

![graph.svg](https://raw.githubusercontent.com/Merovius/AdventOfCode/main/2021/day24/graph.svg)

With a bit of squinting, you can see that the code simulates a stack machine,
by encoding the stack into a base 26 number. `Push`ing a number (less than 26)
unto the stack is done by multiplying with 26 and then adding it to the result.
For example, the first two nodes (at the very bottom-left) push `input[0]+4` and
`input[1]+10` on the stack.

This is followed by a block which does a `conditional push`:

1. It compares `input[2]+6` with `input[3]` (I introduced a `!=` operator to
   make this more visible).
2. The left part of that block multiplies the stack with 26, if they are unequal.
   If they are equal, it is multiplied with 1.
3. The right part of the block adds `input[3]+14` to the stack, if they are unequal.
   If they are qeual, it adds 0.
4. As a result, the block in total pushes `input[3]+14` to the stack, if and
   only if `input[2]+6 == input[3]`.

This is repeated a couple times with other constants.

Roughly halfway through the code (around the middle of the graph) a new kind of
block appears. This block `pop`s a value of the stack (by dividing by 26 and
taking the remainder with 26), subtracts 9 from the result and compares it to `input[10]`.
If the two are unequal, `input[10]+8` is pushed on the stack, otherwise nothing happens.

At the very end, `z` must be 0, which is equivalent to the stack being empty.

This is the entire decoded program in pseudocode:

```
input: abcdefghijklmn

v = []
v.push(a+4)
v.push(b+10)
if c+6 != d {
	v.push(d+14)
}
v.push(e+6)
if f+7 != g {
	v.push(g+1)
}
v.push(h+7)
if i+3 != j {
	v.push(j+11)
}
if v.pop() != k+9 {
	v.push(k+8)
}
if v.pop() != l+5 {
	v.push(l+3)
}
if v.pop() != m+2 {
	v.push(m+1)
}
if v.pop() != n+7 {
	v.push(n+8)
}
```

To get the solution, we observe that there are four unconditional `pop`s and
four unconditional `push`es. To end with an empty stack, we must prevent all
the conditional pushes. By lining up the pushes and pops, we can determine a
set of simple equations, which we can use to minimize or maximize the
respective digits.
