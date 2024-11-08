# Loop : The Programming Language 
Made to combine the speed and low-level power of C with the ease of use and syntax sugar of modern Languages

Loop aims to be simple to use by absolute beginners without sacrificing on speed, thus it is a compliled language that trnaspiles to C before using the tiny C compiler `tcc` to compile.
This use of C not only allows for an easy access to cross-compilation, but also fast compilation as well.

For Now Loop is in its infancy, you can check out some example code for an idea of the language, I am currently working on a REPL and thus an interpreted version of the language before
diving into the compiler.

## Examples

### Hello World

```v
module main

fn main() {
    print("Hello World!)
}
```
### Fibbonaci Sequence

```v
module fibonnaci

fn fibbonaci(n: uint) -> if n <= 1 { 1 } else {
    fibbonaci(n-1) + fibbonaci(n-2)
}

fibbonaci(3)
```

### FizzBuzz

```v
numbers := 0..20

for number in numbers {
    mut s : str = ''
    if number % 3 == 0 {
        s += 'Fizz'
    }
    if number % 5 == 0 {
        s += 'Buzz'
    }
    print(s)
}
```

### Generics

```v
module generics
fn sum<T implements Addition> (items: List<T>) -> {
    s := T.zero 
    for item in items {
            s += item
    }
    s
}

sum([1, 2, 3])
sum(["hello", "world"])
sum([2.3, 1e-4, 82.3])
```

### Named Parameters

```v
fn greet({name: str?}) {
        print('Hello {name}!')
}

greet(name: 'Anas')
```

### Rust-like Enums (Algebraic Datatypes) & Matching

```v
enum Color {
    (u8, u8, u8)
    (u8, u8, u8, u8)
    Black
    White
    Red
    Blue
    Green
    (f32, f32, f32)
    (f32, f32, f32, f32)
}


fn toRGB(color: Color) -> {
    match color {
      (r: u8, g: u8, b: u8) => (r, g, b)
      (r: u8, g: u8, b: u8, a: u8) => (r, g, b)
      Black => (0, 0, 0)
      White => (255, 255, 255)
      Red => (255, 0, 0)
      Blue => (0, 0, 255)
      Green => (0, 255, 0)
      (r: f32, g: f32, b: f32) => (int(255 * r), int(255 * g), int(255 * b))
      (r: f32, g: f32, b: f32, a: f32) => (int(255 * r), int(255 * g), int(255 * b))
  }
}

```
             
