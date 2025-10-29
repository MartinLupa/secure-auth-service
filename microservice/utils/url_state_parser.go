package utils

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ExtractDataFromUrlState(c *gin.Context) (string, string, error) {
	state := c.Query("state")
	decoded, _ := base64.StdEncoding.DecodeString(state)

	var info struct {
		RedirectURI                   string `json:"redirectUri"`
		SocialUserNotFoundRedirectURI string `json:"socialUserNotFoundRedirectUri"`
	}

	if err := json.Unmarshal(decoded, &info); err != nil {
		c.String(http.StatusBadRequest, "Invalid state parameter")
		return "", "", err
	}

	redirectUrl, err := url.QueryUnescape(info.RedirectURI)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid redirect URI")
		return "", "", err
	}

	socialUserNotFoundRedirectUrl, err := url.QueryUnescape(info.SocialUserNotFoundRedirectURI)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid social user not found redirect URI")
		return "", "", err
	}

	return redirectUrl, socialUserNotFoundRedirectUrl, nil
}
