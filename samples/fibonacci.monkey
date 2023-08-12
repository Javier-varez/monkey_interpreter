let wrap = fn(self) {
    fn(arg) {
        self(self, arg)
    }
}

let fib = wrap(fn(self, x) {
    if (x < 2) {
        return x;
    }
    return self(self, x - 1) + self(self, x - 2);
});

puts(fib(26))
