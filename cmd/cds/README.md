# CDS

> Cryptographic Data Storage

```bash
usage: 
    ./main (get|put|del|new) [storage-path] [data-key]
stdin:
    [storage-password]~[data-value]EOF
```

## Example

```bash
$ echo $(cat password.txt;cat data.txt) | ./cds put storage.stg example.com 
$ echo $(cat password.txt) | ./cds get storage.stg example.com 
hello, world!
$ echo $(cat password.txt) | ./cds del storage.stg example.com 
$ echo $(cat password.txt) | ./cds new storage.stg example.com 
$ echo $(cat password.txt) | ./cds get storage.stg example.com 
d1c169963dc69b0d73ffb4a16f821640
```
