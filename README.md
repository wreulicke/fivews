# Five Ws

You MUST always write error messages as literal everywhere.
You SHOULD write error message contains ["Five Ws"](https://en.wikipedia.org/wiki/Five_Ws) to clarify the condition.

## Usage

```go
fivews.New("unable to write the contents to foo by service unavaliable")

err := DoSomething() // err.Error() returns "unable to do something"

fivews.Wrap("unable to write the contents", err) // this.Error() returns "unable to write the contents: unable to do something"

fivews.Join("unable to write the contents", cause1, cause2)
```

