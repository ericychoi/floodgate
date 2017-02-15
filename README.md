# floodgate
A testing tool for unix sockets. Acts as a proxy between 2 sockets.

## Command Line
floodgate

## Install
```bash
go get github.com/ericychoi/floodgate
```

## Usage
```bash
% nc -l -U ~/tmp/dst

# on a separate terminal
% floodgate -src ./src -dst ~/tmp/dst

# on a separate terminal
% echo "{}\n" | nc -U ~/tmp/dst

# should see ~/tmp/dst getting the traffic, but floodgate will also report traffic
```

## Library
You can use this inside your testing library to listen on what gets passed between 2 sockets.

`floodgate.Proxy.Out` gives an `io.ReadCloser` through a channel.  You can block the proxying until you read from this `ReadCloser`.

`Close()` close both source and destination connections.

```go
import "github.com/ericychoi/floodgate"

...

outChan := make(chan io.ReadCloser)
defer close(outChan)
fg := &floodgate.Proxy{
  Type:    "unix",
  SrcAddr: src,
  DstAddr: dst,
  Out:     outChan,
  Timeout: time.Duration(timeoutMS) * time.Millisecond,
}
err := fg.StartServer()
if err != nil {
  log.Fatal(err)
}
defer fg.Close()
go fg.Start()

r := <-fg.Out
br := bufio.NewReader(r)
line, err := br.ReadBytes('\n')
if err != nil {
  // on EOF, or other errors
  r.Close()
}
fmt.Print(string(line))
```
