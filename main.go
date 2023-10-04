package main

// https://github.com/arg0naut91/authenticateAndGo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"tgs-devops/api"
	"tgs-devops/auth"
	"tgs-devops/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	_ "github.com/joho/godotenv/autoload"
)

type Greeting struct {
	Message string `json:"message"`
}

func greeting(w http.ResponseWriter, r *http.Request) {
	greeting := Greeting{
		"欢迎访问学院君个人网站?",
	}
	message, _ := json.Marshal(greeting)
	w.Write(message)
}

func main() {

	utils.InitDB()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	registerAPI(router)

	http.Handle("/", router)
	http.ListenAndServe(utils.GetPort(), router)
}

func registerAPI(router *chi.Mux) {
	router.Route("/", func(r chi.Router) {
		// use the Bearer Authentication middleware

		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		router.Get("/", api.LoginPageHandler)
		router.Post("/login", api.LoginHandler)
		router.Post("/register", api.RegisterHandler)

		router.Get("/index", api.IndexPageHandler)
		router.Get("/loginGoogle", auth.HandleGoogleLogin)

		router.Get("/callback", auth.HandleGoogleCallback)
		router.Post("/logout", api.LogoutHandler)

		router.Group(func(r chi.Router) {
			secretKey := "mySecretKey-10101"

			s := oauth.NewBearerServer(
				secretKey,
				time.Second*120,
				&TestUserVerifier{},
				nil)

			r.Use(oauth.Authorize(secretKey, nil))

			r.Post("/token", s.UserCredentials)
			r.Post("/auth", s.ClientCredentials)

			r.Get("/greeting", greeting)
			r.Get("/customers", GetCustomers)
			r.Get("/customers/{id}/orders", GetOrders)
		})

		//https://levelup.gitconnected.com/oauth-2-0-in-go-846b257d32b4
		fileServer := http.StripPrefix("/web/", http.FileServer(http.Dir("./web/")))
		router.Mount("/web/", fileServer)
		//https://stackoverflow.com/questions/63465062/how-to-set-context-path-for-go-chi
	})
}

func GetCustomers(w http.ResponseWriter, _ *http.Request) {
	renderJSON(w, `{
		"Status":        "verified",
		"Customer":      "test001",
		"Customer_name":  "Max",
		"Customer_email": "test@test.com",
	}`, http.StatusOK)
}

func GetOrders(w http.ResponseWriter, _ *http.Request) {
	renderJSON(w, `{
		"status":            "sent",
		"customer":          "test001",
		"order_id":          "100234",
		"total_order_items": "199",
	}`, http.StatusOK)
}

func renderJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}

// TestUserVerifier provides user credentials verifier for testing.
type TestUserVerifier struct {
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*TestUserVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	if username == "user01" && password == "12345" {
		return nil
	}

	return errors.New("wrong user")
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*TestUserVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientID == "abcdef" && clientSecret == "12345" {
		return nil
	}

	return errors.New("wrong client")
}

// ValidateCode validates token ID
func (*TestUserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (*TestUserVerifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	claims["customer_id"] = "1001"
	claims["customer_data"] = `{"order_date":"2016-12-14","order_id":"9999"}`
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*TestUserVerifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	props["customer_name"] = "Gopher"
	return props, nil
}

// ValidateTokenID validates token ID
func (*TestUserVerifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*TestUserVerifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
