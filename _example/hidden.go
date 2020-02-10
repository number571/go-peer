package main

import (
    "encoding/json"
    "fmt"
    "github.com/number571/gopeer"
)

var (
    ADDRESS1 = gopeer.Get("IS_CLIENT").(string)
    ADDRESS3 = gopeer.Get("IS_CLIENT").(string)
)

const (
    ADDRESS2 = ":8080"
    TITLE    = "TITLE"
)

var (
    another1Client = new(gopeer.Client)
    another2Client = new(gopeer.Client)
)

func main() {
    listener1 := gopeer.NewListener(ADDRESS1)
    listener1.Open().Run(handleServer)
    defer listener1.Close()

    // s7sCvRP0q03zNfXlkepjAvqPYhdj7Uz/Jo9cPSHSdQw=
    client := listener1.NewClient(gopeer.ParsePrivate(privateKey1))

    listener2 := gopeer.NewListener(ADDRESS2)
    listener2.Open().Run(handleServer)
    defer listener2.Close()

    // zXRzW0xlgdNhK5hn3LUGTTJJWm+RE179xcFhWFgnnGg=
    another1Client = listener2.NewClient(gopeer.ParsePrivate(privateKey2))

    listener3 := gopeer.NewListener(ADDRESS3)
    listener3.Open().Run(handleServer)
    defer listener3.Close()

    // ZKWBosdx1i9/RkFJesJtHkpl/+NAkOoK5yN5KzAME3Y=
    another2Client = listener3.NewClient(gopeer.ParsePrivate(privateKey3))
    another2Client.Sharing.Perm = true
    another2Client.Sharing.Path = "./"

    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS2,
        Public:  another1Client.Keys.Public,
    })
    client.Connect(dest)
    another2Client.Connect(dest)

    dest2 := gopeer.NewDestination(&gopeer.Destination{
        Address:  ADDRESS2,
        Public:   another1Client.Keys.Public,
        Receiver: another2Client.Keys.Public,
    })
    client.Connect(dest2) // Hidden connection
    // client.LoadFile(dest2, "archive.zip", "output.zip")
    client.SendTo(dest2, &gopeer.Package{
        Head: gopeer.Head{
            Title:  TITLE,
            Option: gopeer.Get("OPTION_GET").(string),
        },
        Body: gopeer.Body{
            Data: "hello, world!",
        },
    })
    client.Disconnect(dest2)
    client.Disconnect(dest)
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(TITLE, pack,
        func(client *gopeer.Client, pack *gopeer.Package) (set string) {
            fmt.Printf("[%s]: '%s'\n", pack.From.Sender.Hashname, pack.Body.Data)
            return set
        },
        func(client *gopeer.Client, pack *gopeer.Package) {
            // after receive result package
        },
    )
}

func printJSON(data interface{}) {
    jsonData, _ := json.MarshalIndent(data, "", "\t")
    fmt.Println(string(jsonData))
}

var (
    privateKey1 = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAxSdUMJ6+8M4yPirk+4MTCQING6zrYAiME0HdMx+qqwQLdkA3
Y5TVzCVZ2MgrixtEmgEPvDnc9vjjmgfjUhluT3u4wkof/kxyOhaG3hN/j5ucAs7X
soDgmq60tSC5msx9xMpYLft7Ez/94zuxW+Ygqc4wfnp8SgjGY7vFaBfaacDH7nF/
oMxYI0YoZ5yOPBZJMwdPqmkO7sw4YuB9tGgY+2KQRFbLlTrlybrtVGo5EKAkm42D
uz0J2/K+ncACwuFdmvGYweIEhkVZlD+QR6lluiFoV5bEFEcMpTSmaSNILaSQ1h4A
xYpuDzjzXjjtOAap7HpqJ42HO5jHuZBFysfjdQIDAQABAoIBAQC7AnxMhjgOSTjF
WYDMxl9G+ygd6V93P4RHPAGrXc1Q3MxWhcFEd0h5lbBs/iq3j8z53Cnl3Gkp55pV
YEgTd0X4pR3zRcalPDRZv0Z83rfwK6XH0BYwHylt8Gw/J2SHXpOqGFmefF4ZO2kD
o3qv9lFjYM8FGgBNZZdxwQoWnBG0ntuEnlNwDMjoh8QRVaSp0Wr7Jil3ZxfCv3HD
slQXHfewWEsMSsP9SQrKGwjEapXJdHa0AcoS5TCxOnHsed+6oEjrd+UMOf3tf/3P
ROJCfyNe3yQr5keYgqW7Vl1zfUUKv7XSDw0ORDJ+1e5t54Owbv9KDg0Qhh62xQQy
utBceWOhAoGBAN0bxeH08CIEMlDhCEYGeS+gCLN443HvrP0m9e3acATM/n+i9w3N
SE/UPVZSbAKkm9HDptYdOVlfklaSchwSG29HzHotqi4VLWQnwSJX5IWVhlhN4ts/
BVfe1wnytoQ0WkXeoHdB+bylIEnEt97xE/EN2siyRR1pHE4yX3O07AX5AoGBAORD
1vusS7mRqPzgTyGjPZT/rZv149JeEUPPuMReDUiZUAZZaEXXiuireij4lO7nOGmW
8/NN+GhW8EmCS/5eDJx+3goVNEh7RheZPwidLOLitdeg3O6ET5IZyxY/OSbKflKF
L4DaMSNj5IwU+wwP1DqdCJqftvhk+9FptZlrZXhdAoGAVED6BaE9Q+kPd0xYx749
vY0g46rEGK144LpQ6kLbfqjSrbZep+66iFjayqL7r4IkMil40IwwR0Mo0z5YpyOr
OptEaqYt/ANr2YdgjAKr/M8+czWypVL9aT9r98l6DSSZ5Zfw06DbViwiApooapa4
v5lE7kcoGQ3tkEqXntKpQZkCgYBOcandz3YujYofbQ6Eps6w97S31ia20rDoNuhu
Q0wZWOaRaARXjB0mnFdc4SB9gWR3lPK8+FyXrtjgyjBHeBapaUfw+xx+lC6gSX/J
/AnC5tpLAfMq3Lljog2S1aNUW15SCYcrptAgM2IFaehkWsQ6aGDekmiUsE+Bxews
jyXN2QKBgHbDjtoZo+NTRqZ5DbkA0mQKrn94xm8oDwnGgxhZS1YzMPmh7HWxPo6V
E2sbHJznkIB5bbFgqKOeFdp1OO0lFBHS4sPPJewa7ZhLm53nKt7r4K6gDPB+0tGp
Uw9T6QJ/EZF6+ou0Nf5cRFeXCZUcjZLvNw7S/z7b5ZQElREyeH+f
-----END RSA PRIVATE KEY-----`
    privateKey2 = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA1Z6npUgaPTmJUJN9JcnAxEeKdfYkHtZ4HJyO0kzjFyCVnQZH
Dinnit7IPdPOdcHKpD0upMyIbbYt2Agg+u+8KOuHzP2DMs3TKstHWKpU/+byZpRA
8DH/1QvJQXpKjV06Qxd3/JyDND8K0ZExf01lvb3vVvyDDmLK0Qiay0et+hHXWqa8
EK5KWMdepIVRoHVXRjMGPqE9a23mB9dFYNmUzxcMti4IL5Y6rUvRmbso0tpWQQSy
I31mwlll1I4ojxscD1ypyvNVawKoVQlfuVIeNev+dAQaT1VTNK041FK8sqHToMTN
rWQZlrFYQOLRdJUVtRJO6xXowrh3ZCNYxlVsHwIDAQABAoIBAFM46B8gI/jOPYzC
qPLb0tmk9XBXYGMTMuASriGICsCr1R3DoFMISEh12pUbu0dtJEEwBMf3Vv9HBj0v
jYm1dByNBe76pO5Z+XamkzkbwtmfY7hK8bGiCQU6/kEgH4NLWrNgpUIox4THOrPC
WQI7aPOu11uQLI6iNlmRfJzNZB4Ttw8VLjZfkNv8A8o0vS6VIBlvwDat//KZsOA7
xo+oSFc5HIwvrVIqRLfe/+9UDaFFF/c0ZYRA8SOIa1j3J8zY7771SbgL7UZNy2ld
djmGF5+on7D/f7zDlUAuf0IHEZ9IwjakG/pkjsvUdcVkkyxW/5/iFGmrXUtpJBCO
PGPeE4ECgYEA7pCojBa4kyk8GIyPLm0kDrB+D59R8/gWPMIEg8abVbUBQpDf300V
jmlGZO3IguIN3RLN1AHb9+TRu9tlbC78JqUn0LBJUc1sgVqMr3hB2FTQP46V7YO8
AApf+ON6+VE3Ya8F7MnG1fXMBz5lejo5xMT62MxKuvG3DDSCNtz2MvECgYEA5TtK
nhLEl78oIfm3QbAvpboBOr0quTiws/Y+0Lp11YvklMqAenFoBMJmUuZiGDshAWtf
h0tRFfHEoxnmrDZnGrcKScpGLPuvkAahbiDvX8VFn5jLftnDaqFCoI1a9IcNTX5V
+T+xi9D9r2bkkbv1kPWHb5XfSGEVeis8GDmPcA8CgYBpi+C2EftZSF4JMm7KiIjy
Ys1zFfbJLJKSEPi0YHMrCSjkjXourkkCN7toPfd/SIn/rCkaSjRKyZatOVT29xah
9mHWJ5hYs7z0wd4KZ/chwexcojXc3nKXxf9N+z7V/UO2WRwS9fadhODet5Fn8UjL
sKaWslPBv91Pbg/KPBpE0QKBgGxGsVwxKUM0O9Swi2svuZHiZipEqCWNLYoTyl/1
cytHRNUzQbSUVLnKyWJnB/bCFzkAasMRRF/FL0iLN3YozFLGGsn4DWW6DJdPSnkm
fWsTV2unVJe6bJ+1RC9qFYhjMllkT1/IQij0sp9jTpu32Kp8D1kZxbn+gZKPUXdv
2NzpAoGBAKi/UrsQzrbytDeNs/o3s9eGnzsbNNoFWqHe1yVmD043dcQ742NpEDrs
9nqiS9tDqKk4Re/FYclhGjzg/0dnVIdHZjITEhW4kfii8fYPCPedvkCc56psjMBj
UgayZ0CUFTOmyA3RkiniAbVJvHE3i3xGxAv/4Re1sl2LjuMNpVfB
-----END RSA PRIVATE KEY-----`
    privateKey3 = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAylTftJ50lH7PloqWPSDc+6fRGf6d5cNzSX79LB697fIlinre
GrquDdkYNYh91Zp0EJd5RwAeynjXQS7iKWgZx71C3rASQYPbUWAP66kNb4WTY+JL
PdqWKZ2pWU7+sndOEGcyB832ZExTBC4WXWaYLnlyR4n72W03jnH2rJIi4VdLyAMa
vXEMyA53AaZJpeZ5Mg2qIqnE503lrlJ1DNUGrwQm3VGD1te7/oYBWgpr8VU33pfD
y7YhKmdM8wyuLJPDzCpgASMKI7K3YTmzGDDFCcUmgZsVDdQbITboBvuaRlMgEfK6
IOnQoiPZW56qgVHnWZNc/Vw23dGYZygKnb5vkwIDAQABAoIBAQCk1iCqda1kjStx
6dYcVvmxzDZ+hwD8fw8dgWeg2irB/9S1zQPFovDKN2ORuXFK5FpKah1TyrVLHse/
QwLd2QGnyHkCE0/MMDAtS6WsyD057gj2BxZlff94SAn/yGuX2bqvgmMwGIvzinrX
nPR7g2nX8vk3byLPMDtiwVXFogjoq1KksoKN/UT1AXw72h5k+angM7nkKpzIvdyl
TYE8F/dzj1az1ykoE8H2ImP64iQHAUp0/Q3v9yWkB+klBkvTgS3HmPIG/36YXvQB
Bky3tudSpXZrTtZhbrVgtm2D3lOy8AXFwZtblY24B+OL6ePSWw1dH8C0rnVlaBlB
3vEgM3KBAoGBAPAZak62OIF6JFj+6Xa2KPsChB2TSA7K2MiXCEGDGSMOQBimGciL
/NjG+3hpW//ORVdqxzV9IDvXwF4b2wmB3wiyn/EP1hsy+ZyCIoxxdgONRElevFOE
MDXxZs6cOYTKlTCpSdlhIJ6AZt3kts28jrS9LAsIS/Oicw+kBtF1lvhNAoGBANe7
J0zBLixlb3Fp0Btv5LM5wIF5pkpi6HfUk0erD6VHlvlpF9USj4E0JRGBzEUfH7/4
uhBuUkhPIokNvcxva07P3Iy1YNpEH+vYSkEgsKJLRHH191+lsZg+aLH0e4vF1NTW
Z4LJN7LLouuQ7LJc5Rt6+Ctsvm96cQdy2I8ovvdfAoGAGz0W5WUg1feQZhRNUi1q
SsfHSz+pPhxfKaqQwjXoRSTZurIlXK4c+k7gupFhYYz6KuevP+85F/DrHwIUAGke
b3MsWAHO7XkD/nB5EOvSUqbVJ2m6/dKSUZxYaHoqwFjnQgUCnsm5FKJGiUfoQUDy
A6kudPX0/+ffG9gk+eBYR0UCgYAAuxORAP6FC/rqqW8ZCLH/oWxzg9P6YIdlEIVH
Mt8ksi9ivOZlxGBUEbcmbgghG8/huJf4wkbpE8uMJ03DSYVViQK4P00KsDxjciIe
QlwW0KZ0tF6YJlmJqHx2Tdu1R4BHEErdeI0FwAbXQXBr0kC8bRg2HXIvsnx7h/oP
0hWDhQKBgQCovRSD6l6IKYRANwePAe/0CH/VN4o8rXzv2S7u0kXp+Q1QhF78dp1g
HOeW5a8ir+F79uzUC0NTdNoXBxNuyZ0Q5UZ/tmBxtMVNoUREi4Aje4F6Mn+k7IR3
9+C8QWvdEiGmUMuuxBgPgDyxhBrDGii2w77EOuonEJbgvqobpFxr6Q==
-----END RSA PRIVATE KEY-----`
)
