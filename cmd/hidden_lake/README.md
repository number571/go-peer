## HL

> Hidden Lake

<img src="_images/hl_logo.png" alt="hl_logo.png"/>

The `Hidden Lake` is an anonymous network built on a `micro-service` architecture. At the heart of HL is the core - `HLS` (service), which generates anonymizing traffic and combines many other services (for example, `HLT` and `HLM`). Thus, Hidden Lake is not a whole and monolithic solution, but a composition of several combined services.

## Build and run

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

## Connections

<table style="width: 100%">
  <tr>
    <th>Available network</th>
    <th>Types of services</th>
  </tr>
  <tr>
    <td>
        <table style="width: 100%">
            <tr>
                <th>ID</th>
                <th>Type</th>
                <th>Host</th>
                <th>Port</th>
                <th>Connections</th>
                <th>Provider</th>
                <th>Characteristics</th>
                <th>Expired time</th>
            </tr>
            <tr>
                <td>1</td>
                <td>HLTr/HLTs</td>
                <td>v1989393.hosted-by-vdsina.ru</td>
                <td>9581/9582</td>
                <td>[]</td>
                <td><a href="https://vdsina.ru">vdsina.ru</a></td>
                <td>1x3.4GHz, 1.0GB RAM, 30GB HDD</td>
                <td>±02.08.2025</td>
            </tr>
            <tr>
                <td>2</td>
                <td>HLTr</td>
                <td>6a20015eacd8.vps.myjino.ru</td>
                <td>49356</td>
                <td>[1]</td>
                <td><a href="https://jino.ru">jino.ru</a></td>
                <td>1x2.0GHz, 1.5GB RAM, 10GB HDD</td>
                <td>±28.07.2026</td>
            </tr>
            <tr>
                <td>3</td>
                <td>HLTr</td>
                <td>195.133.1.126</td>
                <td>9581</td>
                <td>[1,2]</td>
                <td><a href="https://ruvds.com">ruvds.ru</a></td>
                <td>1x2.2GHz, 0.5GB RAM, 10GB HDD</td>
                <td>±28.07.2027</td>
            </tr>
        </table>
    </td>
    <td>
        <table style="width: 100%">
            <tr>
                <th>Type</th>
                <th>Name</th>
                <th>Default port</th>
            </tr>
            <tr>
                <td>HLS</td>
                <td>node</td>
                <td>9571</td>
            </tr>
            <tr>
                <td>HLTr</td>
                <td>relayer</td>
                <td>9581</td>
            </tr>
            <tr>
                <td>HLTs</td>
                <td>storage</td>
                <td>9582</td>
            </tr>
        </table>
    </td>
  </tr>
</table>
