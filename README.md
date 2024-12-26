# The Lemon Programming Language

Lemon is a tiny, fast, interpreted language. This repository contains the source code for the Lemon interpreter, written in Go. The language is still in development, and the interpreter is not yet feature-complete. A peek at the REPL atm:

![repl](.github/image.png)

## Hello, Lemon!

```js
let message = "Hello, Lemon!";
print(message);
```

For more examples, check out the [examples](./examples) directory.

```js
let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      return 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), push(accumulated, f(first(arr))));
    }
  };
  iter(arr, []);
};
let numbers = [ 1, 1 + 1, 4 - 1, 2 * 2, 2 + 3, 12 / 2 ];
map(numbers, fibonacci);
// => returns: [1, 1, 2, 3, 5, 8]
```

## Features

- [x] Data types: boolean, integer and string
- [x] Variables `let name = value;`
- [x] Arithmetic operations (`+`, `-`, `*`, `/`)
- [x] Logical operations (`!`, `&&`, `||`)
- [x] Functions `fn (args) { body }`
- [x] First-class functions (closures)
- [x] Comparison operations (`==`, `!=`, `>`, `>=`, `<`, `<=`)
- [ ] Control structures
  - [x] if, else `if (condition) { body } else { body }`
  - [ ] else if branch
  - [ ] switch, case
  - [ ] match, when
  - [ ] while, for, loop
  - [ ] break, continue
- [x] Garbage collection
- [x] Strings `let name = "value";`
- [x] String concatenation `"value" + "value";`
- [x] Arrays `[1, 2, 3]`
- [x] Hash maps `{ "key": "value" }`
- [ ] Comments
- [ ] Error handling
- [ ] Standard library
- [ ] Modules
- [ ] Classes
- [ ] Generics
- [ ] Multithreading

## Usage

To run the REPL, use the following command:

```bash
go run main.go
```

To build the interpreter, use the following command:

```bash
go build
```

Running [lemon files](./example.mm) (after building):

```bash
lemon example.mm
```

To run tests, use the following command:

```bash
go test <directory>
```
