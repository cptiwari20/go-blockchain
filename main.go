package main
 
import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBlock).Methods("POST")

	log.Println("Listening on the port 3000")
	log.Fatal(http.listenAndServe(":3000", r))
}