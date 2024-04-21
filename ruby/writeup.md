# Adventure

For this assignment, I reimplemented my JavaScript interpreter, originally written in Go, in Ruby.

### Difficulty

How difficult was this task? Honestly, not as bad as I thought it would be. While I had never touched Ruby in my life, I
chose it because it is one of the most widely known languages that I don't know. I feel that there are a lot of aspects
of the language that I didn't touch in this assignment, but the language itself was easy enough to learn in a short
time.

### Contributing Factors

Some factors that contributed to the difficulty of learning Ruby (or lack thereof) were lack of C-style syntax, unique
and interesting keywords, and prior knowledge and experience in object-oriented languages. Using completely different
syntactical structure than what I'm used to from Java, Go, and Python took some time to get used to. However, Ruby is a
very widely used language, so there is a plethora of help available through sites like Stack Overflow. Probably the most
difficult aspect of the project was translating an object-oriented, statically-typed language into a dynamically-typed
language where you can just do what you want and errors won't be caught until runtime.

### Degree of Expression as a Native Ruby Developer

While I feel that there are many aspects of Ruby that I didn't even come close to using, I felt that as I was coding,
I used Ruby how it was intended. For example, there are many things that are 'optional' in Ruby style-wise, that are
only optional because other languages use them. One instance of this is the `return` keyword. At the end of a function,
you don't actually need to specify if you're returning a value, you can just put that value standalone. Another example
is when calling functions, you can put a space and then the argument instead of putting it in parentheses, like `f x`
instead of `f(x)`. The Ruby convention is typically to do the thing that Ruby introduces rather than sticking to what
the developer may be used to from other languages, and I felt that I did a good job of doing that. Of course, it was
weird translating Go into Ruby but I did my best to stay away from Go principles as I did so.