I use a 15/23 element integer array to represent the burrow, numbering like this:

```
•••••••••••••••••••
•0 1. 2. 3. 4. 5 6•
••••7 •8 •9 •10••••
   •11•12•13•14•
   •••••••••••••
 ```

I started writing writing a Dijkstra-based Python solution. It turned out to be
extremely slow, so I abandoned it. A couple of friends pointed out that it's
possible to just manually solve Part 1 so I ended up doing that, re-using much
of my Python code to check moves, visualize the state and keeping track of
moves. I used an ipython shell, imported the code and experimented, e.g. the
example starts as:

```
> import part1
> s = part1.input
> s
#############
#...........#
###B#C#B#D###
  #A#D#C#A#
  #########
0

> s.move(9,2)
#############
#...B.......#
###B#C#.#D###
  #A#D#C#A#
  #########
40

> s.move(8,9)
#############
#...B.......#
###B#.#C#D###
  #A#D#C#A#
  #########
440

> s.move(12,3)
#############
#...B.D.....#
###B#.#C#D###
  #A#.#C#A#
  #########
3440
```

and so on.

For part 2, I started by extending the code to work on the larger burrow and
tried playing around as well, but doing it manually turned out to not work very
well. So I went back to writing code, but I used Go, to make it reasonably fast
(and eliminate a significant source of bugs, because I don't know Python very
well).

It's probably possible to optimize this a lot, but the Go code ended up being
reasonably fast.
