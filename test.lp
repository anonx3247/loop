comp Person {
    name: string
    age: u8
}

fn main() {
    me := Person(
        name: "Anas",
        age: 22
    )

    me.greet()
}

fn greet(name: string) {
    print("Hello, {name}!")
}