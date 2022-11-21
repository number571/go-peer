# CDS

> Cryptographic Data Storage

```bash
usage: 
    ./main (get|put|del) [path] [storage-password] [data-password]
```

## Example

```bash
$ echo hello, world! > hello.txt 
$ go run main.go put file.stg some-password#1 some-password#2 < hello.txt
$ go run main.go get file.stg some-password#1 some-password#2
> hello, world!
$ go run main.go del file.stg some-password#1 some-password#2
```
