#!anko

func fib(n) {
  a, b = 1, 1
  f = []
  for i in range(n) {
    f += a
    b += a
    a = b - a
  }
  return f
}


println(fib(20))
