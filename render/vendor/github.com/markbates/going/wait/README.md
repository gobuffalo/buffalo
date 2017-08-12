# github.com/markbates/going/wait

Wait cleans up the pattern around using `sync.WaitGroup`.

## Installation

```bash
$ go get https://github.com/markbates/going/wait
```

## Usage

Before:

```go
users := []User{User{"A"}, User{"B"}}
var w sync.WaitGroup
length := len(users)
w.Add(length)
for i := 0; i < length; i++ {
  go func(w *sync.WaitGroup, index int) {
    user := users[index]
    user.Name = fmt.Sprintf("User: %d", index)
    users[i] = user
    w.Done()
  }(&w, i)
}
w.Wait()
```

After:

```go
users := []User{User{"A"}, User{"B"}}
Wait(len(users), func(index int) {
  user := users[index]
  user.Name = fmt.Sprintf("User: %d", index)
  users[index] = user
})
```
