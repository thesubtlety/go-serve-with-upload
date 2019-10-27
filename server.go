package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"path/filepath"
)

var user string
var pass string
var dir = flag.String("d", ".", "directory to serve up")

func main() {
	host := flag.String("h", "0.0.0.0", "interface to serve on")
	port := flag.String("p", "8000", "port to serve on")
	insecure := flag.Bool("k", false, "don't use TLS")
	userpass := flag.String("u", "", "user:pass for basic auth")
	certfile := flag.String("kc", "server.pem", "path to cert file")
	keyfile := flag.String("kf", "server.key", "path to key file")
	flag.Usage = usage
	flag.Parse()

	if *userpass != "" {
		if !strings.Contains(*userpass, ":") {
			fmt.Println("Err: user pass format is user:pass")
			usage()
		}
		log.Printf("Starting with authentication..")
		user = strings.SplitN(*userpass, ":", 2)[0]
		pass = strings.SplitN(*userpass, ":", 2)[1]
		http.HandleFunc("/", requireAuth(mainHandler))
		http.HandleFunc("/serve/", requireAuth(fsHandler(*dir, "/serve/")))
		http.HandleFunc("/upload", requireAuth(uploadHandler))
	} else {
		log.Printf("Starting without authentication..")
		http.HandleFunc("/", mainHandler)
		http.HandleFunc("/serve/", fsHandler(*dir, "/serve/"))
		http.HandleFunc("/upload", uploadHandler)
	}

	if *insecure {
		log.Printf("Server started listening on http://%s:%s", *host, *port)
		log.Fatal(http.ListenAndServe(*host+":"+*port, nil))
	} else {
		log.Printf("Server started listening on https://%s:%s", *host, *port)
		log.Fatal(http.ListenAndServeTLS(*host+":"+*port, *certfile, *keyfile, nil))
	}
}

func usage() {
	fmt.Printf("Usage: ./%s\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func requireAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if !checkCreds(user, pass) {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Can I help you?\"")
			http.Error(w, "Unauthorized", 401)
			return
		}
		fn(w, r)
	}
}

func checkCreds(u, p string) bool {
	if u == user && p == pass {
		log.Printf("Successful auth for user %s", u)
		return true
	}
	log.Printf("Got auth request for user %s", u)
	return false
}

func fsHandler(dir, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	serve := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		logRequest(w, req)
		serve(w, req)
	}
}

//https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(5 << (10 * 3)); // 5GB limit
	err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := filepath.Base(handler.Filename)
	log.Printf("Uploading file %s", filename)
	for k, v := range handler.Header {
		log.Printf("\t%v:%v", k, v)
	}

	f, err := os.OpenFile(*dir + "/" + filename, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("File successfully uploaded...")

	http.Redirect(w, r, "/serve", 302)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(w, r)

	t, err := template.New("index").Parse(`
<!DOCTYPE html>
<head>
        <title>Simple Transfer</title>
<style>
.content {
	font-family: "Helvetica", sans-serif;
    vertical-align: top;
	margin: auto;
	width: 600px;
	padding: 20px;
	border: 1px solid grey;
}
.button{
        background: #2472FE;
        border: none;
        padding: 7px;
        border-radius: 4px;
        color: #D2E2FF;
}
</style>
</head>

<body>
<div class="content">
<h3>Upload Files</h3>
<form enctype="multipart/form-data" action="/upload" method="POST">
                <p label for="avatar">Choose a file to upload: </label><br/>
                <input type="file" name="file" />
                <input type="submit" value="Upload" class="button">
</form>

<hr style="width:400px text-align:left">
<h3>Download Files</h3>
<div><a href="/serve">Browse Files</a></div>

</div>
</body>
</html>`)
	if err != nil {
		http.Error(w, "It dun broke...", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	t.Execute(w, r.RemoteAddr)
	if err != nil {
		http.Error(w, "It dun broke...", http.StatusInternalServerError)
		log.Print(err)
	}
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s from %s, User-Agent = %s\n", r.URL, r.RemoteAddr, r.Header.Get("User-Agent"))
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		log.Printf("Form[%q] = %q\n", k, v)
	}
}
