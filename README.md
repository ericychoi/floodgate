# floodgate
A testing tool for unix sockets. Acts as a proxy between 2 sockets.  Payload delimiter is by default a newline.

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
You can use this inside your testing library to listen on what gets passed betwen 2 sockets.
