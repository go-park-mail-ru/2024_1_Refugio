package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type Response struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

var AUTH_URL = "https://oauth.vk.com/authorize?client_id=51916655&redirect_uri=https://mailhub.su/auth-vk&response_type=code&scope=email"

func main() {
	http.HandleFunc("/auth-vk", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := r.URL.Query().Get("code") //r.FormValue("code")
		conf := oauth2.Config{
			ClientID:     os.Getenv("APP_ID"),
			ClientSecret: os.Getenv("APP_KEY"),
			RedirectURL:  os.Getenv("API_URL"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.vk.com/authorize",
				TokenURL: "https://oauth.vk.com/access_token",
			},
		}

		if code == "" {
			w.Write([]byte(`
				<div>
					<a href="` + AUTH_URL + `">auth</a>
				</div>
			`))
			return
		}

		token, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Println("cannot exchange", err)
			w.Write([]byte("=("))
			return
		}

		fmt.Println("TOKEN OK")

		client := conf.Client(ctx, token)
		resp, err := client.Get(fmt.Sprintf(os.Getenv("API_URL"), token.AccessToken))
		if err != nil {
			log.Println("cannot request data", err)
			w.Write([]byte("=("))
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("cannot read buffer", err)
			w.Write([]byte("=("))
			return
		}

		data := &Response{}
		json.Unmarshal(body, data)
		fmt.Println("photo: ", data.Response[0].Photo, "first_name: ", data.Response[0].FirstName)

		w.Write([]byte(`
		<div>
			<img src="` + data.Response[0].Photo + `"/>
			` + data.Response[0].FirstName + `
		</div>
		`))
	})

	http.ListenAndServe(":8007", nil)
}
