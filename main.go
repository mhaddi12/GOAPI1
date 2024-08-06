package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Profile struct {
	Department  string `json:"department"`
	Designation string `json:"designation"`
	Employee    User   `json:"employee"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var profiles []Profile = []Profile{}

func emailExists(email string) bool {
	for _, profile := range profiles {
		if profile.Employee.Email == email {
			return true
		}
	}
	return false
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var profile Profile

	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if emailExists(profile.Employee.Email) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Email already exists",
		})
		return
	}

	profiles = append(profiles, profile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)

	fmt.Printf("Added profile: %+v\n", profile)
}

func getAllProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	var deleteEmployee User

	err := json.NewDecoder(r.Body).Decode(&deleteEmployee)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, val := range profiles {
		if deleteEmployee.Email == val.Employee.Email {
			profiles = append(profiles[:i], profiles[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(profiles)
			fmt.Printf("Deleted profile: %+v\n", deleteEmployee)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "Profile not found"})
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	var updateEmployee Profile

	err := json.NewDecoder(r.Body).Decode(&updateEmployee)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, update := range profiles {
		if update.Employee.Email == updateEmployee.Employee.Email {
			profiles[i] = updateEmployee
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(profiles)
			fmt.Printf("Updated profile: %+v\n", updateEmployee)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "Profile not found"})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/profile", addItem).Methods("POST")

	router.HandleFunc("/profile", getAllProfile).Methods("GET")

	router.HandleFunc("/profile", deleteEmployee).Methods("DELETE")

	router.HandleFunc("/profile", updateEmployee).Methods("PUT")

	fmt.Println("Server is starting at port 5000...")

	err := http.ListenAndServe(":5000", router)

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
