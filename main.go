package main

import (
    "fmt"
    "time"
    "crypto/sha256"
    "encoding/hex"
    "strconv"
    "html/template"
    "net/http"
    "timezone"
)

type Block struct{
    Index int
    Timestamp time.Time
    Data string
    PreviousHash string
    Hash string
}

func createHash(b *Block) {
    s := strconv.Itoa(b.Index) + b.PreviousHash + b.Timestamp.String() + b.Data 
    h := sha256.New()
    h.Write([]byte(s))
    hash := hex.EncodeToString(h.Sum(nil))

    b.Hash = hash
}

type Blockchain struct{
    chain [10]Block
    len int
}

func createGenesis (bChain *Blockchain) Block{
    b := Block{
        Index: 0, 
        Timestamp: timezone.TimeIn(time.Now(), "Asia/Dhaka"), 
        Data: "First Block", 
        PreviousHash: "NULL", 
    }
    createHash(&b)
    bChain.chain[bChain.len] = b
    bChain.len++

    return b
}

func addBlock (bChain *Blockchain, i int, data string, previousHash string) Block{
    b := Block{
        Index: i, 
        Timestamp: timezone.TimeIn(time.Now(), "Asia/Dhaka"), 
        Data: data, 
        PreviousHash: previousHash, 
    }
    createHash(&b)
    bChain.chain[bChain.len] = b
    bChain.len++

    return b
}

func latestBlock (bChain *Blockchain) Block{
    return bChain.chain[bChain.len - 1]
}

func showAllBlocks (bChain *Blockchain){
    for i := 0; i < bChain.len; i++{
        fmt.Println("Block", bChain.chain[i].Index)
        fmt.Println("Block Data", bChain.chain[i].Data)
        fmt.Println("Block Created at", bChain.chain[i].Timestamp)
        fmt.Println("Block Hash", bChain.chain[i].Hash)
        fmt.Println("Block Previous Hash", bChain.chain[i].PreviousHash)
        fmt.Println("")
    }
}

func main() {
    bChain := Blockchain{len: 0}
    
    // First Block
    createGenesis(&bChain)

    // Second Block 
    addBlock(&bChain, bChain.len, "Second Block", bChain.chain[bChain.len-1].Hash)

    // Third Block 
    addBlock(&bChain, bChain.len, "Third Block", bChain.chain[bChain.len-1].Hash)

    // Fourth Block 
    addBlock(&bChain, bChain.len, "Four Block", bChain.chain[bChain.len-1].Hash)

    // Displaying Blocks
    showAllBlocks(&bChain)

    tmpl := template.Must(template.ParseFiles("index.html"))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")

        data := bChain.chain
        err := tmpl.Execute(w, data)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    http.ListenAndServe(":3030", nil)
}
