package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type Candidate struct {
	Id        int
	Name      string
	CreatedAt time.Time
}

type CreateCandidate struct {
	Name string
}

func connectDB() (*sql.DB, error) {

	host := getEnv("DEMO_HOST", "localhost")
	dbPort := getEnv("DEMO_DBPORT", "5432")
	dbUser := getEnv("DEMO_USER", "postgres")
	password := getEnv("DEMO_DATABASE_PASSWORD", "mysecretpassword")
	dbName := getEnv("DEMO_DBNAME", "postgres")
	sslMode := getEnv("DEMO_SSLMODE", "disable")
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, dbPort, dbUser, password, dbName, sslMode)
	db, err := sql.Open("postgres", connectString)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, err
}

type Handlers struct {
	Db *sql.DB
}

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if body != nil {
		payload, _ := json.Marshal(body)
		w.Write(payload)
	}
}

func (h *Handlers) DeleteCandidateHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	sql := `DELETE FROM Candidates WHERE id = $1`
	res, err := h.Db.Exec(sql, id)
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	contagem, err := res.RowsAffected()
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	if contagem == 0 {
		WriteResponse(http.StatusNotFound, map[string]string{"Erro": "Record not found"}, w)
		return
	}

	WriteResponse(http.StatusNoContent, nil, w)
}

func (h *Handlers) UpdateCandidateHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var newCandidate CreateCandidate
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&newCandidate); err != nil {
		WriteResponse(http.StatusBadRequest, map[string]string{"error": "invalid"}, w)
		return
	}
	sql := `UPDATE Candidates SET Name = $2 WHERE id = $1`
	res, err := h.Db.Exec(sql, id, newCandidate.Name)
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	contagem, err := res.RowsAffected()
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	if contagem == 0 {
		WriteResponse(http.StatusNotFound, map[string]string{"Erro": "Record not found"}, w)
		return
	}

	WriteResponse(http.StatusOK, map[string]string{"alterado": newCandidate.Name}, w)
}

func (h *Handlers) CreateCandidateHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var newCandidate CreateCandidate
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&newCandidate); err != nil {
		WriteResponse(http.StatusBadRequest, map[string]string{"error": "invalid"}, w)
		return
	}
	sql := `INSERT INTO Candidates (Name, created_at) VALUES($1,$2)`
	if _, err := h.Db.Exec(sql, newCandidate.Name, time.Now()); err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	WriteResponse(http.StatusCreated, nil, w)
}

func (h *Handlers) CandidatesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	Candidates, err := h.Db.Query(`SELECT * FROM Candidates`)
	responseCode := http.StatusOK
	if err != nil {
		fmt.Println(err)
		message := map[string]string{"error": "error connecting to database"}
		switch err {
		case sql.ErrNoRows:
			responseCode = http.StatusNotFound
			message["cause"] = "no cadidates"
		default:
			responseCode = http.StatusInternalServerError
			message["cause"] = "database error"
		}
		WriteResponse(responseCode, message, w)
		return
	}
	var lista []Candidate
	for Candidates.Next() {
		var Candidate Candidate
		if err := Candidates.Scan(&Candidate.Id, &Candidate.Name, &Candidate.CreatedAt); err != nil {
			panic(err)
		}
		lista = append(lista, Candidate)
	}

	WriteResponse(responseCode, lista, w)
}

func main() {
	porta := getEnv("DEMO_PORT", "8080")
	db, err := connectDB()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	h := Handlers{Db: db}
	router := mux.NewRouter()
	router.HandleFunc("/candidates", h.CandidatesHandlerFunc).Methods("GET")
	router.HandleFunc("/candidate", h.CreateCandidateHandlerFunc).Methods("POST")
	router.HandleFunc("/candidate/{id}", h.UpdateCandidateHandlerFunc).Methods("PUT")
	router.HandleFunc("/candidate/{id}", h.DeleteCandidateHandlerFunc).Methods("DELETE")
	err = http.ListenAndServe(fmt.Sprintf(":%s", porta), router)
	fmt.Println(err)
}
