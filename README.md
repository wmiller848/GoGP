# GoGP #
### Go Genetic Programing ###

## Usage ##

`go get github.com/wmiller848/GoGP`

`GoGP learn -c 4 -p 20 -g 10 < some.dat`

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
will output `1000`.

This will output the contents of a [CoffeScript](https://github.com/jashkenas/coffeescript) program that evolved to match the input and desired output.

## About ##

GoGP is a DNA inspired genetic programmer.

What does that mean? Lets break it down:

### Initialization ###

GoGP is heavly inspired by nature, namely natural
selection mixed with random mutation. At the start
of a new learning session the specifed population
size is generated. Each new program in this population
is defined by its `DNA`. `DNA` is defined as `[]byte`,
the first generation programs just have random bytes
for their `DNA`.

### DNA ###

So `DNA` is a `[]byte`, but how does it work?

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
