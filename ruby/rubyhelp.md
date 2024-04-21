# Ruby Quick Reference

## Functions
### Define a function
```
def foo
  do stuff
end
```
```
def foo2(x)
  do stuff
end
```
```
def foo3(x = "Initial Value")
  do stuff
end
```
### Call a function
`foo` or `foo()`
`foo2(10)` or `foo2 10`

## Classes
Class instance variables are created with `@`. They are all private, but adding `attr_accessor :<name>` allows access and mutation.
### Define a class
```
class Bar
  def initialize(x = "Initial Value")
    @x = x
  end
  def print_x
    puts x
  end
```
### Instantiate a class
`bar = Bar.new("hi")`



## Conditionals
### Basic if / else if / else
```
if condition
  puts "True"
elsif condition2
  puts "Maybe"
else
  puts "False"
end
```
### Switch statement
```
case x
when 1..5
  "It's between 1 and 5"
when 6
  "It's 6"
when "foo", "bar"
  "It's either foo or bar"
when String
  "You passed a string"
else
  "You gave me #{x} -- I have no idea what to do with that."
end
```



## Equality
### `==`
This is what `eql?` uses under the hood for `Object`s. Does what you would think for strings and numbers.
### `===`
Case equality operator. Implemented differently in different classes, but often used as instanceof with the type on the left, for example:
```
String === "test" # true
Range === (1..2) # true
```
It can also be used to determine if a value lies within a range:
```
(1..4) === 3 # true
("a".."d") === "c" # true
```
This is the equality operator used under the hood of `case/when` statements. Might be considered the "user-friendly" equals.
### `equal?`
Checks if they are literally the same object, the strictest equality check in Ruby.



## Loops
### While
```
while true
  puts x
end
```
### For
```
for i in 1..5
  puts i
end
```
### Iterate over a list without index
```
l.each do |item|
  puts item
end
```



## Miscellaneous
### Comments
When a line begins with `#`, it is a comment. Like Python.
### Print something
`puts "Hello World!"`
### Stringify variable
`puts "Value of x: #{x}"`
### Check nullity
`x.nil?` returns boolean
### Verify method existence on object
`x.respond_to?("method_name")` returns a boolean
### Check if this file is the program entry point
In Python, we do something like `if __name__ == '__main__'`. The Ruby version is `if __FILE__ == $0`.