func a() { return 2 }
is(2, a(), "func a() { return 2 }")

func b(x) { return x + 1 }
is(3, b(2), "func b(x) { return x + 1 }")

func c(x) { return x, x + 1 }
is([2,3], c(2), "func c(x) { return x, x + 1 }")

func d(x) { return func() { return x + 1 } }
is(3, d(2)(), "func d(x) { return func() { return x + 1 } }")

var x = func(x) {
  return func(y) {
    x(y)
  }
}(func(z) {
  return "Yay! " + z
})("hello world")

is("Yay! hello world", x, "...")

# vim: set ft=anko:
