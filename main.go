package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/deevanshu-k/URL-Shortener/libs"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var db = make(map[string]string)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Println(color.GreenString("%v", r.Method), r.Host, r.RequestURI, r.ContentLength)
		next.ServeHTTP(w, r)
		fmt.Println(color.GreenString("%v", r.Method), r.Host, r.RequestURI, r.ContentLength, color.RedString("Time Taken: %v", time.Since(start)))
	})
}

func main() {
	host := flag.String("host", "127.0.0.1", "Default 127.0.0.1")
	port := flag.Int("port", 8000, "Default 8000")
	baseurl := flag.String("baseurl", "http://"+*host+":"+strconv.Itoa(*port), "Default 'http://127.0.0.1:8000'")
	flag.Parse()

	// Templates
	rootTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	shortUrlTemplate, err := template.ParseFiles("templates/urlgenerated.html")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// Serve public folder
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	// Middlewares
	router.Use(loggingMiddleware)

	// Routes
	/*
	 * POST Request, form-data: { url }
	 * Return html page with shortUrl
	 */
	router.HandleFunc("/api/shorturl", func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		url := r.FormValue("url")
		hash := libs.ComputeShortHash(url, &db)

		db[hash] = url

		shortUrlTemplate.Execute(w, struct {
			Url      string
			ShortUrl string
		}{Url: url, ShortUrl: *baseurl + "/" + hash})
	}).Methods(http.MethodPost)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootTemplate.Execute(w, struct{}{})
	})

	router.HandleFunc("/{hash}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]
		url, ok := db[hash]
		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.Redirect(w, r, url+"?ref=shortner", http.StatusTemporaryRedirect)
	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Start Server
	fmt.Println(color.GreenString("Listening On:"), *host+":"+strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(*host+":"+strconv.Itoa(*port), router))
}
