mod introduction


-- These are comments

-- Lets begin with a hello-world:

println("Hello, world!")

-- simple enough right? Ok lets make it slightly more complex

fn greet(name: str?) -> none { -- we define a function "greet that takes a single optional string argument "name" and has no return value
    println(
        "Hello, " + 
        str(if name is none { "world"} else name) -- as you can see ifs can also be passed into expressions
        + "!"
    )
}

-- case: name is none
greet() -- "Hello, world!"

-- case: name is str
greet("Bob") -- "Hello, Bob!"

-- this is called method-notation, it is used to use OOP paradigms
"Bob".greet() -- "Hello, Bob!"

-- as you can see argument names can be written for more clarity
-- when there is only one argument it isn't obligatory
-- but for more than one it is obligatory
greet(name: "Bob") -- "Hello, Bob!"


-- Lets do some math

@cache
fn fibonnaci(n: u64) -> u64 {
    ret if n == 0 or n == 1 {
        1
    } else {
        fibbonaci(n-1) + fibbonaci(n-2)
    }
}

-- Here we defined a recursive function to get the fibbonaci numbers
-- as you can see there is a "@cache" marker above it, this means that during execution,
-- each call will store the value in a hashmap. Thus instead of redoing the calculation
-- again each successive call with the same value will simply look in the hashmap
-- this implements caching for us automatically

mut first_100_fibbonaci_numbers : []u64 = 0..100.map(fibbonaci) 

-- here we created the list of the first 100 fibbonaci numbers
-- notice the "mut" declaring the variable as mutable which means we can change it

-- here we import the math library's square root function
import math::sqrt

fn is_prime(n: u64) -> bool {
    if n in [0, 1] {
        ret false
    }
    for i in 2..u64(math::sqrt(n)) { -- we only need to check the divisors up to sqrt(n)
        if n % i == 0 {
            ret false
        }
    }
    ret true
}

first_100_numbers : []u64 = 0..100

-- lets get the set of the prime-indexed fibbonaci numbers for a prime less than 100:

my_set := first_100_numbers.filter(is_prime).map(fibbonaci)

-- lets optimize out is_prime function

@cache
fn optimized_is_prime(n: u64) -> bool {
    if n == 0 or n == 1 {
        ret false
    } else {
        for i in 2..u64(math::sqrt(n)).filter(optimized_is_prime) { -- we only need to check for prime divisors
            if n % i == 0 {
                ret false
            }
        }
        ret true
    }
}

-- While it will be slower on the first few runs, after a run for a large number, it will calculate very quickly for others

-- Lets work with enums and structs a little bit

-- say we retake our greet function, but want to call it on a Person struct:

enum Sex {
    Male -- sex takes one of these two values
    Female
}

struct Person {
    name: str
    age: u8
    profession: str? -- maybe has no job, the "?" denotes the optionality of the argument
    sex: Sex
}

bob := Person(name: "Bob", age: 23, profession: "Programmer", sex: Male)

bob.greet() -- "Hello, Bob!"

-- the greet function automatically extracted the correct information as the arguments coincided!

-- Lets say we define structs and want methods to run on them, but haven't implemented them yet:


enum Terrain {
    Land
    Sea
    Air
}

struct Car {
    model: str
    top_speed: f32
}

-- there is no need for interfaces in loop, as long as you have the attributes, you can run the function:

struct Boat {
    model: str
    name: str
}

fn print_model(model: str) -> none {
    println(model)
}

ferrari := Car(model: "La Ferrari", top_speed: 300)
cruise := Boat(model: "Cruise Ship", name: "Titanic")

cruise.print_model()
ferrari.print_model()

fn terrain(a: Any) -> Terrain {
    ret switch a {
        a is Car => Land
        a is Boat => Sea
        a is Plane => Air
        def => Land
    }
}