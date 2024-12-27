# ping
A ping tool by go.

## Build
```shell
go build -o ping main.go
```

## Usage
ping [options] &lt;destination&gt;


Options:

| args | value           | desc                                                |
|------|-----------------|-----------------------------------------------------|
| -h   |                 | print help and exit                                 |
| -w   | &lt;timeout&gt; | time to wait for response                           |
| -s   | &lt;size&gt;	   | use &lt;size&gt; as number of data bytes to be sent |
| -c   | &lt;count&gt;	  | stop after &lt;count&gt; packets sent               |

