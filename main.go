package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"html/template"
	"log"
	"net/http"
	//"os"
	"path/filepath"
	"sync"
	//"trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var host = flag.String("host", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈

	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("1")
	gomniauth.WithProviders(
		facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
		github.New("", "", "http://localhost:8080/auth/callback/github"),
		google.New("", "", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// Bootstrapをダウンロードする場合
	//http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/assetsへのパス/"))))

	// チャットルームを開始します
	go r.run()

	// Webサーバーを起動します
	log.Println("Webサーバーを開始します。ポート: ", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
