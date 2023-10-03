package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func printUserInformation(token string) {
  uri := "https://graph.microsoft.com/v1.0/me/messages"
  req, err := http.NewRequest("GET", uri, strings.NewReader(""))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Set("Authorization", "Bearer " + token)


  if err != nil {
    fmt.Println("Request err:", err)
    return
  }

  // Create an HTTP client and send the request
  client := &http.Client{}
  // reqBody, _ := ioutil.ReadAll(req.Body)
  // fmt.Println("Request body from request", string(reqBody))
  resp, err := client.Do(req)

  if err != nil {
    fmt.Println("Error sending request:", err)
    return 
  }
  defer resp.Body.Close()
  // Read and print the response body
  respBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error reading response:", err)
    return 
  }


  // Check the response status code
  fmt.Println("Response Status Code:", resp.Status)
  fmt.Println("Response body", string(respBody))
}

func main() {
  var clientId string
  var clientSecret string
  var tenant string
  var token string

  flag.StringVar(&clientId, "client_id", "", "Client id")
  flag.StringVar(&clientSecret, "client_secret", "", "Client secret")
  flag.StringVar(&token, "token", "", "Token (Optional)")

  flag.Parse()

  if clientSecret == "" || clientId == "" {
    fmt.Printf("Error: Required flags are missing %s %s %s !\n", clientId, clientSecret, tenant)
    flag.PrintDefaults()
    os.Exit(1)
  }

  if token == "" {
    var code string
    OpenLoginPage(clientId)
    fmt.Print("code> ")
    fmt.Scanln(&code)
    token = GetTokenFromCode(clientId, clientSecret, code)
  }
  
  fmt.Println("================")
  fmt.Println("Token: ", token)
  fmt.Println("================")
  // printUserInformation(token)
  SendMail(token, "12tanmayvijay@gmail.com", "Test message from script", "Hello world", time.Now().Format(time.RFC3339))
}



