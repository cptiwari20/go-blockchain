package main
 
import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"encoding/json"
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

func (bchain Blockchain) AddBlock(checkOutInfo BookCheckout){
	// create a new block using previos Hash
	prevBlock := bchain.blocks[len(bchain.blocks) - 1]
	newBlock := createBlock(prevBlock, checkOutInfo)
		// validate the block
		// check previosHash and the old hash
		// check positing
		// check the hash is right hash.
	// if data is validated properly add to the blockchain 

	bchain.blocks = append(bchain.blocks, newBlock)

}

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

	w.Write([]byte("This is response"))
	io.WriteString(w, string(jbytes))
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var bookCheckout BookCheckout

	if err := json.NewDecoder(r.Body).Decode(&bookCheckout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Something went wrong please check : %v", err)
		w.Write([]byte("Something went wrong please check"))
	}

	blockchain.AddBlock(bookCheckout)


}

func main(){

	// TODO:: create a new genesis block on the function start.
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on the port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}