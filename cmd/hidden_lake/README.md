## HL

> Hidden Lake

The `Hidden Lake` is an anonymous network built on a `micro-service` architecture. At the heart of HL is the core - `HLS` (service), which generates anonymizing traffic and combines many other services (for example, `HLT` and `HLM`). Thus, Hidden Lake is not a whole and monolithic solution, but a composition of several combined services.

### Build and run

```bash
$ cd ./cmd/hidden_lake
$ make docker-build
$ make docker-run

> hidden_lake-traffic-1    | [INFO] 2023/06/03 16:45:46 HLT is running...
> hidden_lake-messenger-1  | [INFO] 2023/06/03 16:45:46 HLM is running...
> hidden_lake-service-1    | [INFO] 2023/06/03 16:45:50 HLS is running...
...
> hidden_lake-service-1    | [INFO] 2023/06/03 16:45:57 service=HLS type=BRDCS hash=D81414C4...F703F591 addr=C8F29854...E443A75C proof=0000000001006473 conn=127.0.0.1:
> hidden_lake-traffic-1    | [INFO] 2023/06/03 16:45:57 service=HLT type=UNDEC hash=D81414C4...F703F591 addr=00000000...00000000 proof=0000000001006473 conn=172.20.0.3:9571
> hidden_lake-service-1    | [INFO] 2023/06/03 16:46:02 service=HLS type=BRDCS hash=0615BD44...5DD1B0DB addr=C8F29854...E443A75C proof=0000000000495814 conn=127.0.0.1:
> hidden_lake-traffic-1    | [INFO] 2023/06/03 16:46:02 service=HLT type=UNDEC hash=0615BD44...5DD1B0DB addr=00000000...00000000 proof=0000000000495814 conn=172.20.0.3:9571
...
```
