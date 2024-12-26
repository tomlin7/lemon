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
  5 iter(arr, []);
};
let numbers = [ 1, 1 + 1, 4 - 1, 2 * 2, 2 + 3, 12 / 2 ];
map(numbers, fibonacci);
// => returns: [1, 1, 2, 3, 5, 8]