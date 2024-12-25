let fib = fn(x) { 
    if (x < 1) {
        1
    } else { 
        fib(x - 1) + fib(x - 2)
    } 
}

print(fib(5))