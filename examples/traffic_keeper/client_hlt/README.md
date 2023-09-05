## Example

```bash
go run ./main.go w 'hello, world!'
go run ./main.go h
go run ./main.go r cb3c6558fe0cb64d0d2bad42dffc0f0d9b0f144bc24bb8f2ba06313af9297be4 # hash get by 'h' option
```

## Options

The example is compiled for a local executable file. The application can be created with the `-tags=prod` flag. Then all requests will be redirected to HLTs server `6a20015eacd8.vps.myjino.ru:49191`.
