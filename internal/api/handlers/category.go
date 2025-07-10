package handlers

import (
	"log"
	"net/http"
)

func ListCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("List Categories")); err != nil {
		log.Println("Write error: ", err)
	}

}
func GetCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Get Category")); err != nil {
		log.Println("Write error: ", err)
	}
}
func CreateCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Create Category")); err != nil {
		log.Println("Write error: ", err)
	}
}
func UpdateCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Update Category")); err != nil {
		log.Println("Write error: ", err)
	}
}
func DeleteCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Delete Category")); err != nil {
		log.Println("Write error: ", err)
	}
}
