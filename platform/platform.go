package platform

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"turbojet/cli"
)

const (
	gameEndpoint      = "xxx/game.html"
	lobbyEndpoint     = "lobby.html"
	miniLobbyEndpoint = "mlobby.html"
)

func GetGameCDNPath(c *cli.Context, domain string, pID string, gID string) (string, error) {
	url := domain + "/" + gameURL(pID, gID)
	fmt.Printf("GetGameCDNPath: %s\n", url)
	rURL, err := getRedirect(url)
	if err != nil {
		return "", err
	}
	return rURL, err
}

func GetLobbyCDNPath(c *cli.Context, domain string, pID string) (string, error) {
	url := domain + "/" + lobbyURL(pID)
	fmt.Printf("URL: %s\n", url)
	rURL, err := getRedirect(url)
	if err != nil {
		return "", err
	}
	return rURL, err
}

func GetMiniLobbyCDNPath(c *cli.Context, domain string, pID string) (string, error) {
	url := domain + "/" + miniLobbyURL(pID)
	rURL, err := getRedirect(url)
	if err != nil {
		return "", err
	}
	return rURL, err
}

func getRedirect(url string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	var redirectURL string
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectURL = req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return redirectURL, nil
}

func gameURL(pID string, gID string) string {
	url := fmt.Sprintf("%s?property_id=%s&game_id=%s", gameEndpoint, pID, gID)
	return url
}

func lobbyURL(pID string) string {
	url := fmt.Sprintf("%s?property_id=%s&login_name=aaa&sessiontoken=abcd", lobbyEndpoint, pID)
	return url
}

func miniLobbyURL(pID string) string {
	url := fmt.Sprintf("%s?property_id=%s&login_name=aaa&sessiontoken=abcd", miniLobbyEndpoint, pID)
	return url
}
