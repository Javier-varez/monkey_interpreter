let wrap = fn(self) {
    fn(...) {
        self(self, ...);
    };
};

let reduce = fn(arr, initial, callable) {
    let iter = wrap(fn(self, arr, acc) {
        if (len(arr) == 0) {
            return acc;
        };

        let obj = first(arr);
        let newAcc = callable(obj, acc);
        return self(self, rest(arr), newAcc);
    });

    return iter(arr, initial);
};

let fib = fn(n) {
    if (n < 2) {
        return n;
    }

    let state = reduce(2..n+1, [0, 1], fn(idx, state) {
        let f = state[0];
        let s = state[1];
        let next = f + s;
        return [s, next];
    });

    return state[1];
}

puts(fib(36))
