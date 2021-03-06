package restapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

func (c *InitAPI) initDb() {
	dbHost := "127.0.0.1"
	dbPass := "Secret!23"
	dbName := "koto"
	dbPort := "5432"
	dbUser := "postgres"

	port, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Println(err)
		return
	}

	dbConfig := &pgx.ConnConfig{
		Port:     uint16(port),
		Host:     dbHost,
		User:     dbUser,
		Password: dbPass,
		Database: dbName,
	}

	connection := pgx.ConnPoolConfig{
		ConnConfig:     *dbConfig,
		MaxConnections: 5,
	}

	c.Db, err = pgx.NewConnPool(connection)
	if err != nil {
		log.Println(err)
		return
	}
}

// HandleListUser GetUsers from database
func (c *InitAPI) HandleListUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p GetUsers
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "failed-to-convert-json", http.StatusBadRequest)
	}

	resp, err := c.ListUser(ctx, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed-convert-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// HandleCreateUser GetUsers from database
func (c *InitAPI) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p User
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "failed-to-convert-json", http.StatusBadRequest)
		return
	}

	roleID := r.Header.Get("role_id")
	resp, err := c.CreateUser(ctx, &p, roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed-to-convert-data-into-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// HandleUpdateUser update users data from database
func (c *InitAPI) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p User
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "failed-to-convert-json", http.StatusBadRequest)
		return
	}

	roleID := r.Header.Get("role_id")
	resp, err := c.UpdateUser(ctx, &p, roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed-to-convert-data-into-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// HandleDeleteUser delete users data from database
func (c *InitAPI) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	param := mux.Vars(r)
	id := param["userId"]

	roleID := r.Header.Get("role_id")
	resp, err := c.DeleteUser(ctx, &UserID{ID: id}, roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed-to-convert-data-into-json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// StartHTTP function is used to start the app
func StartHTTP() http.Handler {
	api := createAPI()
	api.initDb()

	r := mux.NewRouter()
	r.HandleFunc("/api/user/list", api.HandleListUser).Methods("GET")
	r.HandleFunc("/api/user/create", api.HandleCreateUser).Methods("POST")
	r.HandleFunc("/api/user/update", api.HandleUpdateUser).Methods("PUT")
	r.HandleFunc("/api/user/delete/{userId}", api.HandleDeleteUser).Methods("DELETE")

	return r
}
