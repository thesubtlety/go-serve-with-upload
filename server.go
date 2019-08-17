package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var user = "admin"
var pass = "password"

func main() {
	port := flag.String("p", "8000", "port to serve on")
	insecure := flag.Bool("k", false, "don't use TLS")
	userpass := flag.String("u", "admin:pass", "user:pass for basic auth")
	flag.Parse()

	http.HandleFunc("/", requireAuth(mainHandler))
	http.HandleFunc("/serve/", requireAuth(fsHandler("./", "/serve/")))
	http.HandleFunc("/upload", requireAuth(uploadHandler))

	if *insecure {
		log.Printf("Server started listening on http://host:%s", *port)
		log.Fatal(http.ListenAndServe("localhost:"+*port, nil))
	} else {
		log.Printf("Server started listening on https://host:%s", *port)
		log.Fatal(http.ListenAndServeTLS("localhost:"+*port, "server.pem", "server.key", nil))
	}
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

func logRequest(w http.ResponseWriter, r *http.Request) {
	log.Print("Got new request...")

	fmt.Printf("\tURL Path = %q\n", r.URL.Path)
	fmt.Printf("\t%s %s %s\n", r.Method, r.URL, r.Proto)
	fmt.Printf("\tHost = %q\n", r.Host)
	fmt.Printf("\tRemoteAddr = %q\n", r.RemoteAddr)
	fmt.Printf("\t%s\n", "Headers")
	for k, v := range r.Header {
		fmt.Printf("\t\t%q:%q\n", k, v)
	}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		log.Printf("Form[%q] = %q\n", k, v)
	}

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(w, r)

	t, err := template.New("index").Parse(`
<!DOCTYPE html>
<h3>Welcome</h3>
<body>
<form enctype="multipart/form-data" action="/upload" method="POST">
        <div>
                <label for="avatar">Choose a file to upload: </label><br/>
		<input type="file" name="file" />
                <input type="submit" value="Upload">
	</div>
</form>
<hr>

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
