# Storage

> Cryptographic Data Storage

```bash
usage: 
    ./main (get|put|del) [storage-path] [data-key]
stdin:
    [password]EOL
    [data-value]EOF
```

EOL - End of Line (Enter)
EOF - End of File (Ctrl+D)

## Example

```bash
$ ./main put storage.stg data-key
> [password]EOL
> [data-value]EOF
ok
$ ./main get storage.stg data-key
> [password]EOL
[data-value]
$ ./main del storage.stg data-key
> [password]EOL
ok
```
