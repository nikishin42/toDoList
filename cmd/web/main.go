package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

// структура для хранения всех зависимостей веб-приложения
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "Сетевой адрес HTTP")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := connDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	infoLog.Println("connected to snippetbox database")

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Эта структура позволяет записывать ошибки, которые происходят на сервере в errorLog
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start serve %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

func connDB() (*sql.DB, error) {
	user := os.Getenv("snippetbox_user")
	pass := os.Getenv("snippetbox_pass")
	connSTR := fmt.Sprintf("user=%s password=%s dbname=snippetbox sslmode=disable", user, pass)
	db, err := sql.Open("postgres", connSTR)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
