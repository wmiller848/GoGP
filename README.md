# GoGP #
### Go Genetic Programing ###

A Go implementation of a DNA based Gentic Programmer.

## Usage ##

`go get github.com/wmiller848/GoGP`

`GoGP learn -c 4 -p 20 -g 10 < some.dat > MyProgram.coffee`

`-c` is the number of columns in the input data
that should be learned.

`-p` is the running population to keep around.

`-g` is the number of generations to iterate.


The contents of `some.dat` looks like this:

```
1.0 3.1 5.2 1.0 1000
3.1 5.2 1.8 2.3 2000
```

The last column in each row is the desired output, meaning
GoGP will evolve a program that given the input `1.0 3.1 5.2 1.0`
will output a value as close as possible to `1000`.

Invoking the program is easy `MyProgram.coffee < some.dat`, you
could then pipe that output into your tooling.

This will output the contents of a [CoffeScript](https://github.com/jashkenas/coffeescript) program that evolved to match the input and desired output.

## About ##

GoGP is a DNA inspired genetic programmer.

What does that mean? Lets break it down:

#### Initialization ####

GoGP is heavly inspired by nature, namely natural
selection mixed with random mutation. At the start
of a new learning session the specifed population
size is generated. Each new program in this population
is defined by its `DNA`. `DNA` is defined as `[]byte`,
the first generation programs just have random bytes
for their `DNA`.

#### DNA ####

So `DNA` is a `[]byte`, but how does it work? There are several
steps involved in converting that `[]byte` into a working math
equation to operate on input.

The first step is to sequence the `DNA`, like real `DNA` our
digital `DNA` contains 4 bases and can be read in three frames.
Giving us 4^3=64 possible encodings to work with. Each program
contains two sequences of `DNA`, a `ying` and `yang` strand. Like
real `DNA` each strand is sequenced togeather and produces a single
reading of the genes encoded. The process looks at each of the three
readings in each strand, and produces the gene sequence in order of
index of each gene. For example, lets look at these two stands:

Bases are defined as:
\* Note this may become configurable in the future

```
A = [0x00 to 0x40]
B = [0x40 to 0x80]
C = [0x80 to 0xc0]
D = [0xc0 to 0xff]
```

`ying = [0,  24, 200, 241, 3,  12,  33, 4,  132]`
`yang = [20, 51, 127, 9,   15, 198, 18, 10, 215]`

We could read each strand in three ways:

Strand `ying`:

* [0, 24, 200], [241, 3, 12], [33, 4, 132]
* [24, 200, 241], [3, 12, 33], [4, 132, 0]
* [200, 241, 3], [12, 33, 4], [132, 0, 24]

Strand `yang`

* [20, 51, 127], [9, 15, 198], [18, 10, 215]
* [51, 127, 9], [15, 198, 18], [10, 215, 20]
* [127, 9, 15], [198, 18, 10], [215, 20, 51]

Between those three readings we look for the first one that
contains the start block `AAA`. If no start block is found that
means that strand doesnt sequence to any genes. If we do find
a start block we start reading from the strand from that frame,
we read until we hit and end block `DDD` or the strand ends. We
do the same thing for both strands and then sequence the output
together respecting index.

For example, say strand `ying` has a gene that starts at index 2
and ends at index 20. Say strand `yang` has two genes, one that
starts at index 5 and goes to 15, and the other that starts at index
25 and goes to index 30. We would produce the following gene sequence:

`gene = ying[2:20] + yang[25:30]`

Notice that `yang[5:15]` was not included, that is because stand `ying`
gene sequence ran through those indexes. This works almost the same way
in real `DNA`.

Remembering that we have 64 total encoding, and that the start and stop
blocks take up 2 of those, this means we can define 62 additional encodings.
Currently we define the following blocks:

`+ - * / 0 1 2 3 4 5 6 7 8 9 ,` and one block for each input column as
additional encodings.

The net result is that after we sequence all the genes we could get something
like this `Program Expression`:

`programExpression = '+12,5,$a-10*2000,7/12,$b'`, where `$a` and `$b` refer
to an input column.

Before we move onto `Program Expression`'s lets talk about what this `DNA`
gets us. First, several bases may encode to the same value, `AAA` and `AAB`
may encode to same value for example. Second we can arbitrarily add, remove,
and mutate (XOR) bytes in our `DNA` strands. Mutations, even a single bit,
may lead to no change or completly change the gene sequence, this is
important because of the way we sequence this gives a lot of power to easy
to perform mutations.

#### Program Expressions and Trees ####

Now we enter the the standard genetic programing space. We have a program
expression we can validate and convert into a `Program Tree`. These expressions
are read from right to left and produce a program tree. For example:

`+12,5,$a-10,2`

would convert to the following tree

```
+
12  5  $a  -
           10  2
```

Given a tree, we can then covert that into a valid program. The previous tree
would convert to:

`(12 + 5 + $a) + (10 - 2)`

Using these simple rules, we can evolve arbitrary program trees to estimate
our desired output. Because the mutations occures on the meta level above
the `Program Expression` we gain a lot of resilence to getting stuck in
bad tree evolvotion cycles. Because the dual `DNA` strands encode based
off the index we can preserve 'dominate genes' aka those with a lower index.

#### Output ####

## Legal Foo and Licences ##

This code is licenced under a general BSD Licences, basicaly use
as you see fit commercialy or not, don't take credit for work you
didn't do. Don't use me, anyone or any group affiliated with this
software names, companies or orginizations to advertise for your
company, software or projects.

```
Copyright (c) 2016, William C. Miller
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
3. Neither the names, companies, or orginizations of the copyright holders nor
   the names, companies, or orginizations of its contributors may be used to
   endorse or promote products derived from this software without specific
   prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
