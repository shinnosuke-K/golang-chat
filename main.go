package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	//"os"
	"path/filepath"
	"sync"
	//"trace"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
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
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	// Bootstrapをダウンロードする場合
	//http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/assetsへのパス/"))))

	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()

	// Webサーバーを起動します
	log.Println("Webサーバーを開始します。ポート: ", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
