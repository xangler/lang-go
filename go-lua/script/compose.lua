local gfib = require("gfib")

function fib(n)
    print("lua fib start", n)
    local obj = Obj.new(n)
    print(obj:age())
    local x = gfib.fib(n)
    obj:age(x)
    print(obj:age())
    print("lua fib end", x)
    return x
end