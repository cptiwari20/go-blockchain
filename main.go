package main
 
import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
)

// Structs
type Block struct {
	PreviousHash 		string
	Position			int
	Data				BookCheckout
	Hash				string
	CreatedAt			string
}

type BookCheckout struct {
	BookId		string 		`json:book_id`
	CreatedAt	string		`json:created_at`
	User		string 		`json:user`
	IsGenesis	bool		`json:is_genesis`
}

type Book struct {
	ID			string		`json:"id"`
	Author 		string		`json:"author"`
	PublishDate	string		`json:"publish_date"`
	CreatedAt	string		`json:"created_at"`
	ISBN		string		`json:"isbn"`

}

type Blockchain struct{
	blocks []*Block
}

var blockchain Blockchain
func (b *Block) generateHash()  {
	block := &Block{}

	blockData, _ := json.Marshal(b.Data)
	// each block should be unique in position, time, previous block history and the data
	dataToHash := string(block.Position) + string(blockData) + block.CreatedAt + block.PreviousHash
	// encrypt this data
	hash := sha256.New()
	hash.Write([]byte(dataToHash))

	block.Hash = hex.EncodeToString(hash.Sum(nil))
	fmt.Printf("this is block %v", block)
	fmt.Printf("this is Hash %v", block.Hash)
	fmt.Printf("this is position %v", block.Position)
	fmt.Printf("this is data %v", block.Data)
}

func createBlock(prevBlock *Block, checkOutInfo BookCheckout) *Block  {
	block := &Block{}
	block.PreviousHash = prevBlock.Hash
	block.Data = checkOutInfo
	block.Position = prevBlock.Position + 1
	block.CreatedAt = time.Now().String()
	// generate a new hash and add hash to the new block and with data is checkoutInfo
	block.generateHash()
	return block
}

func isValid(block *Block, prevBlock *Block) bool {
	// check previosHash and the old hash
	if block.PreviousHash == prevBlock.Hash {
		return true
	}
		// check positing
	if block.Position == prevBlock.Position +1 {
		return true
	}
		// check the hash is right hash.
	if block.validateHash(block.Hash) {
		return true
	}
	return false
}

func (block Block) validateHash(hash string) bool  {
	block.generateHash()

	if block.Hash == hash {
		return true
	}
	return false
}

func (bchain Blockchain) AddBlock(checkOutInfo BookCheckout){
	// create a new block using previos Hash
	prevBlock := bchain.blocks[len(bchain.blocks) - 1]
	newBlock := createBlock(prevBlock, checkOutInfo)
		// validate the block
	if isValid(newBlock, prevBlock) {
		// if data is validated properly add to the blockchain 
		bchain.blocks = append(bchain.blocks, newBlock)
	}

}
// new Blockchan
func createGenesisBlock() *Block {
	return createBlock(&Block{}, BookCheckout{IsGenesis: true})
}
func NewBlockChain() *Blockchain {
	return &Blockchain{
		[]*Block{createGenesisBlock()},
	}
}

// handle routes
func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err !=nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create :%v", err)
		w.Write([]byte("could not create a new book"))
		return
	}

	hash := md5.New()
	io.WriteString(hash, book.ISBN + book.PublishDate)
	book.ID = fmt.Sprintf("%x", hash.Sum(nil))
	
	resp, err := json.MarshalIndent(book, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to marshal the json :%v", err)
		w.Write([]byte("Could not marshal the json"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(blockchain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)

	// w.Write([]byte("This is response"))
	io.WriteString(w, string(jbytes))
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var bookCheckout BookCheckout

	if err := json.NewDecoder(r.Body).Decode(&bookCheckout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Something went wrong please check : %v", err)
		w.Write([]byte("Something went wrong please check"))
		return
	}

	blockchain.AddBlock(bookCheckout)

	resp, err := json.MarshalIndent(bookCheckout, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Something went wrong after creating a block %v", err)
		w.Write([]byte("Something went wrong afater creating block, the checkout item not found!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main(){

	// create a new genesis block on the function start.

	BlockChain := NewBlockChain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Prev. hash: %x\n", block.PreviousHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data: %v\n", string(bytes))
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("Listening on the port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}