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
)

var user string
var pass string

func main() {
	host := flag.String("h", "0.0.0.0", "interface to serve on")
	port := flag.String("p", "8000", "port to serve on")
	insecure := flag.Bool("k", false, "don't use TLS")
	userpass := flag.String("u", "", "user:pass for basic auth")
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
		http.HandleFunc("/serve/", requireAuth(fsHandler("./", "/serve/")))
		http.HandleFunc("/upload", requireAuth(uploadHandler))
	} else {
		log.Printf("Starting without authentication..")
		http.HandleFunc("/", mainHandler)
		http.HandleFunc("/serve/", fsHandler("./", "/serve/"))
		http.HandleFunc("/upload", uploadHandler)
	}

	if *insecure {
		log.Printf("Server started listening on http://%s:%s", *host, *port)
		log.Fatal(http.ListenAndServe(*host+":"+*port, nil))
	} else {
		log.Printf("Server started listening on https://%s:%s", *host, *port)
		log.Fatal(http.ListenAndServeTLS(*host+":"+*port, "server.pem", "server.key", nil))
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
		return true
		log.Printf("Successful auth for user %s", u)
	}
	log.Printf("Got auth request for user %s", u)
	return false
}

func fsHandler(dir, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	h := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		logRequest(w, req)
		log.Printf("Serving file %s", req.URL)
		h(w, req)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got upload request...")

	if err := r.ParseMultipartForm(5 << (10 * 3)); // 5GB limit
	err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := handler.Filename
	log.Printf("Uploading file %s", filename)
	for k, v := range handler.Header {
		log.Printf("%v:%v", k, v)
	}

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	log.Printf("File successfully uploaded...")
	//fmt.Fprintf(w, "Successfully uploaded file!") //can't redirect if you display this

	http.Redirect(w, r, "/serve", 302)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(w, r)

	t, err := template.New("index").Parse(`
<!DOCTYPE html>
<head>
        <title>Simple Go Serve...</title>
<style>
* {
        font-family: "Helvetica", sans-serif;
}
.form{
        border-radius:2px;
        padding:5px;
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
<h1>Welcome</h1>
<h3>Upload Files</h3>
<form enctype="multipart/form-data" action="/upload" method="POST" class="form">
                <p label for="avatar">Choose a file to upload: </label><br/>
                <input type="file" name="file" />
                <input type="submit" value="Upload" class="button">
</form>

<hr style="width:400px text-align:left">
<h3>Download Files</h3>
<div><a href="/serve">Browse Files</a></div>

</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	log.Print("")
	fmt.Printf("\t%s %s %s\n", r.Method, r.URL, r.Proto)
	fmt.Printf("\tRemoteAddr = %q\n", r.RemoteAddr)
	fmt.Printf("\tUser-Agent = %s\n", r.Header.Get("User-Agent"))
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		log.Printf("Form[%q] = %q\n", k, v)
	}

}
