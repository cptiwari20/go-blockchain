package main
 
import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

// Structs
type Block struct {

}

type BookCheckout struct {

}

type Book struct {
	ID			string		`json:"id"`
	Author 		string		`json:"author"`
	PublishDate	string		`json:"publish_date"`
	CreatedAt	string		`json:"created_at"`
	ISBT		string		`json:"isbt"`

}

type Blockchain {
	blocks []*Block
}

func newBook (w http.ResponseWriter, r *http.Request) {
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book); err !=nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create :%v", err)
		w.Write([]byte("could not create a new book"))
		return
	}

	h := md5.New()
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on the port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}