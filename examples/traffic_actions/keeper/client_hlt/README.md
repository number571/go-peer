## Example

```bash
go run ./main.go [prod] w 'hello, world!'
> cb3c6558fe0cb64d0d2bad42dffc0f0d9b0f144bc24bb8f2ba06313af9297be4 # stdout
go run ./main.go [prod] h
go run ./main.go [prod] r cb3c6558fe0cb64d0d2bad42dffc0f0d9b0f144bc24bb8f2ba06313af9297be4
```

## Options

The example is default launched for a local environment. The application can be launch with the `prod` flag. Then all requests will be redirected to HLTs server `6a20015eacd8.vps.myjino.ru:49191`.
