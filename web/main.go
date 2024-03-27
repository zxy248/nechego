package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/zxy248/nechego/data"
)

//go:embed static
var static embed.FS

//go:embed templates
var resources embed.FS

var templates = template.Must(template.ParseFS(resources, "templates/*"))

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

	if err := templates.ExecuteTemplate(w, "chat.html.tmpl", data); err != nil {
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

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, config.Database)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServerFS(static))
	mux.Handle("/chat/{id}", &Chat{data.New(conn)})

	log.Println("Listening on", config.Address)
	log.Fatal(http.ListenAndServe(config.Address, mux))
}
