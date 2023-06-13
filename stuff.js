

`

struct name,
    home: int,
    that: this,

fn add [x y] ,
    first: x + y,
    then: first + 1,
    to: then - 2,
    return: (if to
                (+ x y)
                (add y x))

`

const add = (x, y) => {
    var first = x + y;
    var then = first + 1;
    var to = then - 2;
    return (to ? (x + y)
                : add(y, x));
};



var x = 1;

console.log(x++);

console.log(x)