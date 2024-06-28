package apis

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/ONG-Mais/website-back-end/entities"
	"github.com/ONG-Mais/website-back-end/firebase"
	"github.com/ONG-Mais/website-back-end/utils" // Adjust the import path as necessary
	"github.com/gorilla/mux"
)

var client *firestore.Client

func init() {
	client = firebase.InitFirebase()
}

func CreateVoluntary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var voluntary entities.Voluntary
	if err := json.NewDecoder(r.Body).Decode(&voluntary); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Validate email and phone number
	emailValidator := utils.EmailValidator{Email: voluntary.Email}
	if !emailValidator.IsValid() {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	phoneValidator := utils.PhoneValidator{Phone: voluntary.Phone}
	if !phoneValidator.IsValid() {
		http.Error(w, "Invalid phone format", http.StatusBadRequest)
		return
	}
	docRef, _, err := client.Collection("voluntaries").Add(ctx, voluntary)
	if err != nil {
		http.Error(w, "Failed to add voluntary: "+err.Error(), http.StatusInternalServerError)
		return
	}
	voluntary.ID = docRef.ID // Optionally update the ID with the Firestore generated ID
	json.NewEncoder(w).Encode(voluntary)
	w.WriteHeader(http.StatusCreated)
}

// GetVoluntary retrieves a voluntary by ID
func GetVoluntary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"] // Using mux to extract variables from route
	doc, err := client.Collection("voluntaries").Doc(id).Get(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch voluntary: "+err.Error(), http.StatusNotFound)
		return
	}
	var voluntary entities.Voluntary
	doc.DataTo(&voluntary)
	json.NewEncoder(w).Encode(voluntary)
}

// UpdateVoluntary updates an existing voluntary
func UpdateVoluntary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Validate email and phone number
	if email, ok := updates["email"]; ok {
		emailValidator := utils.EmailValidator{Email: email.(string)}
		if !emailValidator.IsValid() {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}
	}
	if phone, ok := updates["phone"]; ok {
		phoneValidator := utils.PhoneValidator{Phone: phone.(string)}
		if !phoneValidator.IsValid() {
			http.Error(w, "Invalid phone format", http.StatusBadRequest)
			return
		}
	}

	_, err := client.Collection("voluntaries").Doc(id).Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		http.Error(w, "Failed to update voluntary: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteVoluntary removes a voluntary from the database
func DeleteVoluntary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	_, err := client.Collection("voluntaries").Doc(id).Delete(ctx)
	if err != nil {
		http.Error(w, "Failed to delete voluntary: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
