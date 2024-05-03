package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"golang.org/x/oauth2"
	"strconv"
)

const (
	APP_ID  = "51916655"
	APP_KEY = "oz3r7Pyakfeg25JpJsQV"
	API_URL = "https://api.vk.com/method/users.get?fields=id,photo_max,email,sex,bdate&access_token=%s&v=5.131"
)

type Response struct {
	Response []struct {
		Id        int    `json:"id"`
		Name      string `json:"first_name"`
		Email     string `json:"email"`
		Photo     string `json:"photo_max"`
		Sex       int    `json:"sex"`
		BirthDate string `json:"bdate"`
		InvitedBy string `json:"invited_by"`
		LastName string `json:"last_name"`
	}
}

var AUTH_URL = "https://oauth.vk.com/authorize?client_id=51916655&redirect_uri=https://mailhub.su/auth/auth-vk&response_type=code&scope=email"

func main() {
	http.HandleFunc("/auth-vk", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := r.FormValue("code")
		conf := oauth2.Config{
			ClientID:     APP_ID,
			ClientSecret: APP_KEY,
			RedirectURL:  "https://mailhub.su/auth/auth-vk",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.vk.com/authorize",
				TokenURL: "https://oauth.vk.com/access_token",
			},
			Scopes: []string{"email"},
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
		fmt.Println("Token: ", token)
		if err != nil {
			log.Println("cannot exchange", err)
			w.Write([]byte("=("))
			return
		}

		fmt.Println("TOKEN OK")

		client := conf.Client(ctx, token)
		resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
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
		fmt.Println("photo: ", data.Response[0].Photo, "first_name: ", data.Response[0].Name)
		w.Write([]byte(`
                <div>
                        <img src="` + data.Response[0].Photo + `"/>
                        <div>` + data.Response[0].Name + `</div>
                        <div>` + data.Response[0].LastName + `</div>
                        <div>` + strconv.Itoa(data.Response[0].Id) + `</div>
                        <div>` + data.Response[0].Email + `</div>
                        <div>` + strconv.Itoa(data.Response[0].Sex) + `</div>
                        <div>` + data.Response[0].BirthDate + `</div>
                        <div>` + data.Response[0].LastName + `</div>
			<div>` + data.Response[0].InvitedBy + `</div>
                </div>
                `))
	})

	http.ListenAndServe(":8007", nil)
}
