package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var user string
var pass string
var dir = flag.String("d", ".", "directory to serve up")
var key = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ3QzUXlpeTRGWkhrWVYKd002ZTVCcjdEODZtWDRuV0xMYlk4d2ZQdlRVbHg0U2tqbkxveVdWS3dXRnBmZlZiUzh2WnBWNm5wbHZEMjM2VgpYTzdST0wvK1U5T282RzM4cENrdWIvRllZOWNJV3FxMWd2WmdzYTJHNzNwWXNxblg0MVJJZVYvb0JnM2QzQVE1CitkTnZteXB4MmZrR3IxbzJFb2piMTFpUUZTZnphYzVJellFVTBLSjI5QXBXVFQ3OTFxQVphdXczdXUwdmRGN3UKaSs2clRVS2lYLzlKK3A4Y0FBTlh1dkdzUmllY0tiRGVsT1M2Q1paaWZ3YVhSc2ZPMjZFOHAxN1ZLYzNZZXBYVQozZGlVMDMvMmtjM3J1QlI5bFJBUmVkZ0x4WkhLeEZ4bUFjYmxGLzJhUkx1emUxQkVobUd1UE9sYUtyZ0Nkc0c0CkNnN2tUM3hQQWdNQkFBRUNnZ0VBSmxwaDVET1NTQTkybEd2ZzZJb1hMWlR5R0I5eEw0N1RrdzRoaGdFT0RWUnkKL1Qzek9VamNFRjZTVjR3U3FONFNqT04rK3VxbXlaRE0zclFPZHBiWE80cFFZYTFNUGZRVXBZcktLWjUwbkFJcwpNRGhBazFuK2xvcVRhYWVYOGVqUytkM1VlMEdDbzVOYVYxTzJBYU92L0VlQ09LaEw5U3VuaHg5OUNPT3gyVzc3CmJhMDlHbzJxSW9pVlVNZWFEZThVZTZxNmdjNjhNci9jZThMWG0zMDRPKzU2QVN0VW9Jc0VJeCtiY1A2blluNDAKOTY1citGaW1MUkpmMmJSakx3Zjc2NVdrdEhMUUtLc2hPSDBGZG9ZRWxIT1pXTk1MTDh3NnRZM09RNFdxWjcyVApUdzZmOXBrWjdVaERZTWtTaVRZUXhkL0VtcDZsamJkcEFLRGtVeTZWU1FLQmdRRGhCRHhUNXZUU2NKbmcvOThOCmxvMnowYjArcnVnQ05KRE5KM1NlYUppUVMreUhuNzRGNHhLNlkvU2laS1Q2K3F6UnIvcWVwYkdGdE1KT3RjcEMKdllGalIwQzArSXlVRmQyYXh4VExmTlJOYTlsYUpKSWJhNFVNMXBMSkZreW9Rcnh2NVRpVHNZZTZCbk51bGRlMwpRdlUxM3A3TWhOQmxZNFFmRHFmL2JOa0FBd0tCZ1FERnphOEFCemw5NUovQlhvOUtwU2N6UFRvbVhMbzd0KzlrCjRMRE1GYVNSNEhoSStlRkVFWCtKMUtxeGd1bkpiZ1RPT0ZYQVRWRDlrN2NBQlVCd2lZT09EdDlXRmNoNWlPWVUKZUp5TmR2MTEwcURWY2k0UXJYNUU5d0NQaTlqam9LVWROS2xxRXRhdG1xOU8rV2F4L0pTZjJieWgyTCtOOC9JQQp2ZTBEY1J0K3hRS0JnQUsrU1hvQVk5VzQ5N2ROaDB1a0hVQW0rM2FyTFRyeHB4NUpMOXZLakttZHMxbUg4Z29pClZaVWVLTnBkL2NEdGszUFBBSEEwdHZCWlh0RVUyRTF1QUFqVTBvNGlSWng4azhJU1VVZVYwd1RLbnREQmgySjgKTWhnUSthTW4rWEZIdHdKcU9nRmE5YnVuM25wbnEwU1p0V0dkd0RQZ0hxWk55MHVSb3l5ekNBWS9Bb0dBQzFWVQpvSkRKWjRBdzh4aGk1Mmo5RFArR0ZHcWR0UXc5NkM3RGtuM3U5dmpBaTVYZHBWUEhWZk5jY0YxSzNlS3kzY24yCmg5VW1QZEUzM0FWeEFzR3VTdlpwTDNxQ0NReWgraXhLOUFRTVU3TGt5allIazZjTkpCQnU5TXFUZTc5Wmxvbk0KNXluN0tPbERBQ2hrRFBDbTUxM0haQktTTHlUNkNiYllIS2xmWk8wQ2dZRUFnQlQ5Z0phL3EzU2FrNGFlUGJHYwovOEZoT0s2bkxjRFhpdVphckk2Y3ZKNndBRVZENC9lTnBuMFZrcUhUK3dDd2J2c2x4Y1RoVHlhRzloM21WUmYxCnNsam5vQmh6UHE4Q2Z1WmxVamhSY05DbDdUM1BaakFCUXo5S05vNWFJR1dkVGM1ek11VnFSalZmRXBMa1NjbVkKcUxsVE1MZ3UzVVRXZDRhcXN6TExjWW89Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K"
var pem = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURjVENDQWxtZ0F3SUJBZ0lVVUFFTzdld0hEdWJOMGZnaWd2R2xvNk5rUXlJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1NERUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2xkaGMyaHBibWQwYjI0eEVEQU9CZ05WQkFjTQpCMU5sWVhSMGJHVXhFakFRQmdOVkJBb01DVTFwWTNKdmMyOW1kREFlRncweU1qQXlNVGN3TURBMU5EQmFGdzB5Ck16QXlNVGN3TURBMU5EQmFNRWd4Q3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSURBcFhZWE5vYVc1bmRHOXUKTVJBd0RnWURWUVFIREFkVFpXRjBkR3hsTVJJd0VBWURWUVFLREFsTmFXTnliM052Wm5Rd2dnRWlNQTBHQ1NxRwpTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFDdDNReWl5NEZaSGtZVndNNmU1QnI3RDg2bVg0bldMTGJZCjh3ZlB2VFVseDRTa2puTG95V1ZLd1dGcGZmVmJTOHZacFY2bnBsdkQyMzZWWE83Uk9MLytVOU9vNkczOHBDa3UKYi9GWVk5Y0lXcXExZ3ZaZ3NhMkc3M3BZc3FuWDQxUkllVi9vQmczZDNBUTUrZE52bXlweDJma0dyMW8yRW9qYgoxMWlRRlNmemFjNUl6WUVVMEtKMjlBcFdUVDc5MXFBWmF1dzN1dTB2ZEY3dWkrNnJUVUtpWC85SitwOGNBQU5YCnV2R3NSaWVjS2JEZWxPUzZDWlppZndhWFJzZk8yNkU4cDE3VktjM1llcFhVM2RpVTAzLzJrYzNydUJSOWxSQVIKZWRnTHhaSEt4RnhtQWNibEYvMmFSTHV6ZTFCRWhtR3VQT2xhS3JnQ2RzRzRDZzdrVDN4UEFnTUJBQUdqVXpCUgpNQjBHQTFVZERnUVdCQlFjdU1hNmx6d1kyWmN5N3RTazlFTjNhMENyM2pBZkJnTlZIU01FR0RBV2dCUWN1TWE2Cmx6d1kyWmN5N3RTazlFTjNhMENyM2pBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUEwR0NTcUdTSWIzRFFFQkN3VUEKQTRJQkFRQmg0MzhoblVyRk1uZ2laWDUwNnF0c29FZ0h6WGkveG14d2ZSOFlvaTR1NG95TFFrYkl3S3JrTXoxTwpFY1J2ajJpcmhmdWg4ZDR6ejg4OHQ0OFV6ODJMSTBYQzd4WW9BS001WkdQV2N4Y1o0aUx6bktJL04xR1NmSC9SCnlNSyt0VW1jcSsvY3NobTcvWlFnM1dENjVRckE2ODhOQ3Q5bU5IZGJSeThxRElnOFo1K1R0cHY3cUFuekIwRjIKbE1sWkxtQ2VOd0w3dDNTOHFlQ2xidS9aTWI4bUNRU05pbm1JeDY3dlQwdDR6MFNvcUROUnVZWkFJYUFsR1VTYwpjRTR6VUxWd3V4RFkrSFZKK1EzWDc1dWNpbC80b2JReHBEdlRFZlVOenlpTFo5WWRUWWlaQ0ZMb2QxV0FrRHFCCmJ6OW1FYnNmN25SM1pHOTVFWW1KQm16MjU0V3QKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="

func main() {
	host := flag.String("h", "0.0.0.0", "interface to serve on")
	port := flag.String("p", "8000", "port to serve on")
	insecure := flag.Bool("k", false, "don't use TLS")
	userpass := flag.String("u", "", "user:pass for basic auth")
	certfile := flag.String("kc", "", "path to cert file")
	keyfile := flag.String("kf", "", "path to key file")
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

		if *certfile == "" || *keyfile == "" {
			fmt.Println("[!] Insecurely using hardcoded SSL cert...")
			pemB, _ := base64.StdEncoding.DecodeString(pem)
			*certfile = string(pemB)

			keyB, _ := base64.StdEncoding.DecodeString(key)
			*keyfile = string(keyB)

			cert, err := tls.X509KeyPair([]byte(*certfile), []byte(*keyfile))
			if err != nil {
				panic(err)
			}
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
			server := http.Server{
				Addr:      *host + ":" + *port,
				TLSConfig: tlsConfig,
			}
			log.Fatal(server.ListenAndServeTLS("", ""))
		}

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

	f, err := os.OpenFile(*dir+"/"+filename, os.O_WRONLY|os.O_CREATE, 0660)
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
