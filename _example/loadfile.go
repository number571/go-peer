package main

import (
    "os"
    "strings"
    "./gopeer"
)

type FileTransfer struct{
    Id uint
    Name string
    Data string
    IsNull bool
}

const (
    BUFFSIZE = 8 // bytes
    ADDRESS = ":8080"
    TITLE = "TITLE"
)

var (
    anotherClient = new(gopeer.Client)
    IsEOF = false
)

func main() {
    listener := gopeer.NewListener(ADDRESS)
    listener.Open().Run(handleServer)
    defer listener.Close()

    client := listener.NewClient(gopeer.GeneratePrivate(1024))
    anotherClient = listener.NewClient(gopeer.GeneratePrivate(1024))

    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := &gopeer.Destination{
        Address: ADDRESS,
        Public: anotherClient.Keys.Public,
    }
    client.Connect(dest)
    loadFile(client, dest, "file.txt")
    client.Disconnect(dest)
}

func loadFile(client *gopeer.Client, dest *gopeer.Destination, filename string) {
    id := uint(0)
    for !IsEOF {
        client.SendTo(dest, &gopeer.Package{
            Head: gopeer.Head{
                Title: TITLE,
                Option: gopeer.Get("OPTION_GET").(string),
            },
            Body: gopeer.Body{
                Data: string(gopeer.PackJSON(FileTransfer{
                    Id: id,
                    Name: filename,
                })),
            },
        })
        id++
    }
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(TITLE, pack, 
        func(client *gopeer.Client, pack *gopeer.Package) (set string) {
            var read = new(FileTransfer)
            gopeer.UnpackJSON([]byte(pack.Body.Data), read)
            name := strings.Replace(read.Name, "..", "", -1)
            data := readFile(name, read.Id)
            return string(gopeer.PackJSON(FileTransfer{
                Id: read.Id,
                Name: name,
                Data: gopeer.Base64Encode(data),
                IsNull: data == nil,
            }))
        },
        func(client *gopeer.Client, pack *gopeer.Package) {
            var read = new(FileTransfer)
            gopeer.UnpackJSON([]byte(pack.Body.Data), read)
            if read.IsNull {
                IsEOF = true
                return
            }
            newFile := "new" + read.Name
            if read.Id == 0 && fileIsExist(newFile) {
                IsEOF = true
                return
            }
            writeFile(newFile, gopeer.Base64Decode(read.Data))
        },
    )
}

func writeFile(filename string, data []byte) error {
    if !fileIsExist(filename) {
        _, err := os.Create(filename)
        if err != nil {
            return err
        }
    }
    file, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    file.Write(data)
    return nil
}

func readFile(filename string, id uint) []byte {
    file, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer file.Close()

    _, err = file.Seek(int64(id*BUFFSIZE), 0) // Beggining of file
    if err != nil {
        return nil
    }

    var buffer = make([]byte, BUFFSIZE)
    length, err := file.Read(buffer)
    if err != nil {
        return nil
    }

    return buffer[:length]
}

func fileIsExist(filename string) bool {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return false
    }
    return true
}
