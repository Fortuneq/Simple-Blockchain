package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"
)

type Block struct {
	Index int 
	Timestamp string 
	BPM int
	Hash string
	PreHash string 
}

var Blockchain []Block


//func calculate hash for block 

func calculateHash(block Block)string{
	record := string(block.Index) + block.Timestamp+ string(block.BPM) + block.PreHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// generating block of puls

func GenerateBlock(oldBlock Block, BPM int)(Block,	error){
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index+1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PreHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	
	return newBlock,nil
}

// func for check is old block valid  ?

func isBlockValid(newBlock,oldBlock Block) bool {
	if oldBlock.Index +1 != newBlock.Index{
		return false
	}
	if oldBlock.Hash != newBlock.PreHash{
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash{
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
    if len(newBlocks) > len(Blockchain) {
        Blockchain = newBlocks
    }
}
//простейшая фунция сервера
func run() error {
    mux := makeMuxRouter()
    httpAddr := os.Getenv("ADDR")
    log.Println("Listening on ", os.Getenv("ADDR"))
    s := &http.Server{
        Addr:           ":" + httpAddr,
        Handler:        mux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    if err := s.ListenAndServe(); err != nil {
        return err
    }

    return nil
}

// фунция для обработки блоков в бразуере 
func makeMuxRouter() http.Handler {
    muxRouter := mux.NewRouter()
    muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
    muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
    return muxRouter
}