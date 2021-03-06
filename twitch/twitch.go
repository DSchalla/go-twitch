package twitch

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const rootURL = "https://api.twitch.tv/kraken/"

type Client struct {
	client      *http.Client
	BaseURL     *url.URL
	AppInfo     *AppInfo
	AccessToken string

	// Twitch api methods
	Channels *ChannelsMethod
	Chat     *ChatMethod
	Games    *GamesMethod
	Ingests  *IngestsMethod
	Search   *SearchMethod
	Streams  *StreamsMethod
	Teams    *TeamsMethod
	Users    *UsersMethod
	Videos   *VideosMethod

	// OAuth
	Auth          *AuthMethods
	OAuthResponse *OAuthResponse
}

type AppInfo struct {
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	State        string `url:"state"`
	RedirectURI  string `url:"redirect_uri"`
	Scope        string `url:"scope"`
}

// Returns a new twitch client used to communicate with the API.
func NewClient(httpClient *http.Client) *Client {
	baseURL, _ := url.Parse(rootURL)

	c := &Client{client: httpClient, BaseURL: baseURL}
	c.Channels = &ChannelsMethod{client: c}
	c.Chat = &ChatMethod{client: c}
	c.Games = &GamesMethod{client: c}
	c.Ingests = &IngestsMethod{client: c}
	c.Search = &SearchMethod{client: c}
	c.Streams = &StreamsMethod{client: c}
	c.Teams = &TeamsMethod{client: c}
	c.Users = &UsersMethod{client: c}
	c.Videos = &VideosMethod{client: c}
	c.Auth = &AuthMethods{client: c}

	clientId := os.Getenv("GO-TWITCH_CLIENTID")
	clientSecret := os.Getenv("GO-TWITCH_CLIENTSECRET")
	state := os.Getenv("GO-TWITCH_CLIENTID")
	redirectURL := os.Getenv("GO-TWITCH_REDIRECTURL")
	scope := os.Getenv("GO-TWITCH_SCOPE")

	c.AppInfo = &AppInfo{}
	c.AppInfo.ClientID = clientId
	c.AppInfo.ClientSecret = clientSecret
	c.AppInfo.State = state
	c.AppInfo.RedirectURI = redirectURL
	c.AppInfo.Scope = scope

	return c
}

// Issues an API get request and returns the API response. The response body is
// decoded and stored in the value pointed by r.
func (c *Client) Get(path string, r interface{}) (*http.Response, error) {
	rel, err := url.Parse(path)

	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err

	}
	req.Header.Add("Accept", "application/vnd.twitchtv.v2+json")

	if len(c.AppInfo.ClientID) != 0 {
		req.Header.Add("Client-ID", c.AppInfo.ClientID)
	}

	if len(c.AccessToken) != 0 {
		req.Header.Add("Authorization", "OAuth "+c.AccessToken)
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, errors.New("api error, response code: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	if r != nil {
		err = json.NewDecoder(resp.Body).Decode(r)
	}

	return resp, err
}
