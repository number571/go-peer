## HL

> Hidden Lake

<img src="_images/hl_logo.png" alt="hl_logo.png"/>

The `Hidden Lake` is an anonymous network built on a `micro-service` architecture. At the heart of HL is the core - `HLS` (service), which generates anonymizing traffic and combines many other services (for example, `HLT` and `HLM`). Thus, Hidden Lake is not a whole and monolithic solution, but a composition of several combined services.

## Build and run

```bash
$ cd ./cmd/hidden_lake
$ make docker-build
$ make docker-run

> hidden_lake-loader-1     | [INFO] 2023/11/22 18:59:27 HLL is running...
> hidden_lake-traffic-1    | [INFO] 2023/11/22 18:59:27 HLT is running...
> hidden_lake-messenger-1  | [INFO] 2023/11/22 18:59:27 HLM is running...
> hidden_lake-service-1    | [INFO] 2023/11/22 18:59:29 HLS is running...
...
> hidden_lake-service-1    | [INFO] 2023/11/22 18:59:36 service=HLS type=BRDCS hash=B04B315A...290EB85C addr=24DC908B...E8299D18 proof=0000365777 size=8192B conn=127.0.0.1:
> hidden_lake-traffic-1    | [INFO] 2023/11/22 18:59:36 service=HLT type=BRDCS hash=B04B315A...290EB85C addr=00000000...00000000 proof=0000365777 size=8240B conn=172.26.0.2:9571
> hidden_lake-service-1    | [INFO] 2023/11/22 18:59:41 service=HLS type=BRDCS hash=E1A8DDFC...674A9E06 addr=24DC908B...E8299D18 proof=0001019421 size=8192B conn=127.0.0.1:
> hidden_lake-traffic-1    | [INFO] 2023/11/22 18:59:41 service=HLT type=BRDCS hash=E1A8DDFC...674A9E06 addr=00000000...00000000 proof=0001019421 size=8240B conn=172.26.0.2:9571
...
```

## Settings

```yaml
# [HLS, HLT]
message_size_bytes: 8192
work_size_bits: 20
key_size_bits: 4096
queue_period_ms: 5000
limit_void_size_bytes: 4096

# [HLT]
hashes_window: 2048
## [62.233.46.109, 94.103.91.81]
messages_capacity: 1048576 # 2^20 msgs ~= 8GiB
## [185.43.4.253]
messages_capacity: 33554432 # 2^25 msgs ~= 256GiB
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
                <th>Version</th>
                <th>Host</th>
                <th>Port</th>
                <th>Network key</th>
                <th>Connections</th>
                <th>Provider</th>
                <th>Country</th>
                <th>City</th>
                <th>Characteristics</th>
                <th>Expired time</th>
            </tr>
            <tr>
                <td>1</td>
                <td>HLTr/HLTs</td>
                <td>v1.5.23</td>
                <td>94.103.91.81</td> 
                <td>9581/9582</td>
                <td>8Jkl93Mdk93md1bz</td>
                <td>[]</td>
                <td><a href="https://vdsina.ru">vdsina.ru</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x4.0GHz, 1.0GB RAM, 30GB HDD</td>
                <td>±eternal</td>
            </tr>
            <tr>
                <td>2</td>
                <td>HLTr</td>
                <td>v1.5.23</td>
                <td>195.133.1.126</td>
                <td>9581</td>
                <td>8Jkl93Mdk93md1bz</td>
                <td>[1]</td>
                <td><a href="https://ruvds.com">ruvds.ru</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x2.2GHz, 0.5GB RAM, 10GB HDD</td>
                <td>±28.07.2027</td>
            </tr>
            <tr>
                <td>3</td>
                <td>HLTr/HLTs</td>
                <td>v1.5.23</td>
                <td>62.233.46.109</td>
                <td>9581/9582</td>
                <td>oi4r9NW9Le7fKF9d</td>
                <td>[]</td>
                <td><a href="https://eternalhost.net">eternalhost.net</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x2.8GHz, 1.0GB RAM, 16GB HDD</td>
                <td>±eternal</td>
            </tr>
            <tr>
                <td>4</td>
                <td>HLTr</td>
                <td>v1.5.23</td>
                <td>193.233.18.245</td>
                <td>9581</td>
                <td>oi4r9NW9Le7fKF9d</td>
                <td>[3]</td>
                <td><a href="https://4vps.su">4vps.su</a></td>
                <td>Russia</td>
                <td>Novosibirsk</td>
                <td>1x2.5GHz, 1.0GB RAM, 5GB VNMe</td>
                <td>±07.08.2027</td>
            </tr>
            <tr>
                <td>5</td>
                <td>HLTs</td>
                <td>v1.5.23</td>
                <td>185.43.4.253</td>
                <td>9582</td>
                <td>j2BR39JfDf7Bajx3</td>
                <td>[]</td>
                <td><a href="https://firstvds.ru">firstvds.ru</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x3.1GHz, 2.0GB RAM, 300GB HDD</td>
                <td>±10.12.2024</td>
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
