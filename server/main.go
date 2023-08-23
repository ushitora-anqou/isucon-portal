package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/uta8a/isucon-portal/server/gen/rpc/auth/v1/authv1connect"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AppServer struct{}

func (s *AppServer) GetUser(ctx context.Context, req *connect.GetUserRequest) (*authv1connect.GetUserResponse, error) {
	return &authv1connect.GetUserResponse{
		User: &authv1connect.User{
			Id:   "1",
			Name: "test",
		},
	}, nil
}

func main() {
	baseUrl := os.Getenv("BASE_URL")
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/api/auth/callback", baseUrl),
		Scopes:       []string{}, // 最小限。public dataのみを取得
		Endpoint:     github.Endpoint,
	}
	url := conf.AuthCodeURL("github_auth")
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	var authCode string
	fmt.Scan(&authCode)
	fmt.Printf("%v\n", authCode)
	// Handle the exchange code to initiate a transport.
	tok, err := conf.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(context.TODO(), tok)
	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", body)
	appServer := &AppServer{}
	path, handler := authv1connect.NewUserServiceHandler(appServer)
}
