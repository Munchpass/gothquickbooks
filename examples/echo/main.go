package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/munchpass/gothquickbooks"
)

func OAuthStart(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, "quickbooks")
	newReq := c.Request().WithContext(ctx)
	user, err := gothic.CompleteUserAuth(c.Response(), newReq)
	if err != nil {
		fmt.Println("err: ", err)
		gothic.BeginAuthHandler(c.Response(), newReq)
		return nil
	}

	fmt.Printf("user (oauth start): %+v\n", user)
	fmt.Printf("raw data:\n%+v\n", user.RawData)
	return nil
}

func OAuthCallback(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, "quickbooks")
	queryParams := c.Request().URL.Query()
	realmId := queryParams.Get("realmId")
	fmt.Println("realmId: ", realmId)
	newReq := c.Request().WithContext(ctx)
	user, err := gothic.CompleteUserAuth(c.Response(), newReq)
	if err != nil {
		return err
	}

	fmt.Printf("user (from callback): %+v\n", user)
	fmt.Printf("raw data:\n%+v\n", user.RawData)

	// Fetch profile
	// Change this if you are using production key
	// From: https://developer.api.intuit.com/.well-known/openid_configuration
	// userInfoEndpoint := "https://accounts.platform.intuit.com/v1/openid_connect/userinfo"

	// From: https://developer.api.intuit.com/.well-known/openid_sandbox_configuration
	userInfoEndpoint := "https://sandbox-accounts.platform.intuit.com/v1/openid_connect/userinfo"
	request, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//set header
	request.Header.Set("accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+user.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("failed to get user info: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	type Address struct {
		StreetAddress string `json:"streetAddress"`
		Locality      string `json:"locality"`
		Region        string `json:"region"`
		PostalCode    string `json:"postalCode"`
		Country       string `json:"country"`
	}

	type UserInfoResponse struct {
		Sub                 string  `json:"sub"`
		Email               string  `json:"email"`
		EmailVerified       bool    `json:"emailVerified"`
		GivenName           string  `json:"givenName"`
		FamilyName          string  `json:"familyName"`
		PhoneNumber         string  `json:"phoneNumber"`
		PhoneNumberVerified bool    `json:"phoneNumberVerified"`
		Address             Address `json:"address"`
	}

	var userInfoResp = new(UserInfoResponse)
	err = json.Unmarshal(body, &userInfoResp)
	if err != nil {
		log.Fatalln("error parsing userInfoResponse:", err)
	}
	fmt.Printf("User Info Resp: %+v\n", userInfoResp)
	log.Println("Ending GetUserInfo")
	return nil
}

func main() {
	qbClientId := os.Getenv("QB_CLIENT_ID")
	qbClientSecret := os.Getenv("QB_CLIENT_SECRET")
	callbackUrl := "http://localhost:3000/quickbooks/callback"
	qbProvider := gothquickbooks.New(qbClientId, qbClientSecret, callbackUrl,
		gothquickbooks.ScopeAccounting, gothquickbooks.ScopeOpenID,
		gothquickbooks.ScopeEmail, gothquickbooks.ScopePhone, gothquickbooks.ScopeProfile)
	goth.UseProviders(qbProvider)

	e := echo.New()
	e.GET("/quickbooks/start", OAuthStart)
	e.GET("/quickbooks/callback", OAuthCallback)

	e.Logger.Fatal(e.Start(":3000"))
}
