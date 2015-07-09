package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Host  string `toml:"host"`
	Token string `toml:"token"`
}

func printErr(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{ \"error\": \"%s\" }", message)
}

func invite(config Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			printErr(w, "not allow GET reuqest")
			return
		}

		email := r.FormValue("email")
		if email == "" {
			printErr(w, "empty email")
			return
		}

		values := url.Values{}
		values.Add("email", email)
		values.Add("token", config.Token)
		values.Add("set_active", "true")

		_, err := http.PostForm("https://"+config.Host+"/api/users.admin.invite", values)
		if err != nil {
			printErr(w, fmt.Sprint(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{ \"success\": 1 }")
		fmt.Println("invite: " + email)
	}
}

func main() {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./assets/")))
	http.HandleFunc("/invite", invite(config))
	http.ListenAndServe(":8080", nil)
}
