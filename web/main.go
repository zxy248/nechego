package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/zxy248/nechego/data"
)

var templates = template.Must(template.ParseGlob("templates/*"))

type Chat struct {
	Queries *data.Queries
}

func (h *Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()
	if _, err = h.Queries.GetChat(ctx, chatID); err != nil {
		http.NotFound(w, r)
		return
	}

	var data struct {
		MessageCount string
		CommandCount string
		TopCommands  string
		TopUsers     string
	}
	data.MessageCount, err = h.Queries.MessageCount(ctx, chatID)
	if err != nil {
		serverError(w, err)
		return
	}
	data.CommandCount, err = h.Queries.CommandCount(ctx, chatID)
	if err != nil {
		serverError(w, err)
		return
	}
	data.TopCommands, err = h.Queries.TopCommands(ctx, chatID)
	if err != nil {
		serverError(w, err)
		return
	}
	data.TopUsers, err = h.Queries.TopUsers(ctx, chatID)
	if err != nil {
		serverError(w, err)
		return
	}

	if err := templates.ExecuteTemplate(w, "chat", data); err != nil {
		serverError(w, err)
		return
	}
}

var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func serverError(w http.ResponseWriter, err error) {
	const code = http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
	errorLog.Printf("%s\n%s\n", err.Error(), debug.Stack())
}

func getenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s not set", key))
	}
	return value
}

func main() {
	addr := getenv("NECHEGO_ADDR")
	db := getenv("NECHEGO_DATABASE")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	mux.Handle("/chat/{id}", &Chat{data.New(conn)})

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
