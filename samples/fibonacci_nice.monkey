let wrap = fn(self) {
    fn(argOne, argTwo, argThree, argFour) {
        self(self, argOne, argTwo, argThree, argFour)
    }
}

let loop = wrap(fn(self, start, end, func, param) {
    if (start < end) {
        let result = func(param)
        return self(self, start + 1, end, func, result)
    }
    return param;
})

let n = 36

let res = loop(2, n+1, fn(param) {
    let pOne = param[0]
    let pTwo = param[1]
    let res = [pTwo, pOne + pTwo]
}, [0, 1])

puts(res[1])
