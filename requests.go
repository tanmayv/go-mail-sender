package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

type TokenResponse struct {
  TokenType string `json:"token_type"`
  ExpiresIn int `json:"expires_in"`
  ExtExpiresIn int `json:"ext_expires_in"`
  AccessToken string `json:"access_token"`
}

func open(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}

func OpenLoginPage(clientId string) {
    open("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=" + clientId + "&response_type=code&redirect_uri=http://localhost:4200/redirect&response_mode=query&scope=offline_access%20user.read%20mail.read%20mail.send")
}

func GetTokenFromCode(clientId string, clientSecret string, code string) string {
  uri := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
  
  requestBody := url.Values{}
  requestBody.Set("client_id", clientId)
  requestBody.Set("client_secret", clientSecret)
  requestBody.Set("grant_type", "authorization_code")
  requestBody.Set("code", code)
  requestBody.Set("scope", "User.Read Mail.Read Mail.Send")
  requestBody.Set("redirect_uri", "http://localhost:4200/redirect")
  // requestBody, err := json.Marshal(TokenRequest{
  //   // ClientId: clientId,
  //   // ClientSecret: clientSecret,
  //   GrantType: "client_credentials",
  //   Scope: "https%3A%2F%2Fgraph.microsoft.com%2F.default",
  // })

  // if err != nil {
  //   fmt.Println("Not able to parst request")
  //   return
  // }

  req, err := http.NewRequest("POST", uri, strings.NewReader(requestBody.Encode()))
  // Set the Content-Type header to indicate JSON data
  // req.Header.Set("Content-Type", "application/json")
  // Set the Content-Type header for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


  if err != nil {
    fmt.Println("Request err:", err)
    return ""
  }

  // Create an HTTP client and send the request
  client := &http.Client{}
  // reqBody, _ := ioutil.ReadAll(req.Body)
  // fmt.Println("Request body from request", string(reqBody))
  resp, err := client.Do(req)

  if err != nil {
    fmt.Println("Error sending request:", err)
    return ""
  }
  defer resp.Body.Close()
  // Read and print the response body
  respBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error reading response:", err)
    return ""
  }


  // Check the response status code
  fmt.Println("Response Status Code:", resp.Status)
  fmt.Println("Response body", string(respBody))
  var tokenResponse TokenResponse

  // Unmarshal the JSON data into the struct
  err = json.Unmarshal(respBody, &tokenResponse)
  if err != nil {
    fmt.Println("Error unmarshalling JSON:", err)
    return ""
  }
  return tokenResponse.AccessToken
}

func SendMail(token string, to string, subject string, body string, delayedTimeUtc string) {
  json_message := fmt.Sprintf(`
  {
    "message": {
      "subject": "%s",
      "body": {
	"contentType": "Text",
	"content": "%s"
      },
      "toRecipients": [
      {
	"emailAddress": {
	  "address": "%s"
	}
      }
      ],
      "singleValueExtendedProperties": [
      {
	"id": "SystemTime 0x3FEF",
	"value": "%s"
      }
      ]
    },
    "saveToSentItems": "true"
  }
  `, subject, body, to, delayedTimeUtc)
  fmt.Printf("Sending to %s: %s at %s \n", to, subject, delayedTimeUtc)
  SendMailRaw(token, json_message)
}

func SendMailRaw(token string, json_message string) {
  uri := "https://graph.microsoft.com/v1.0/me/sendMail"

  req, err := http.NewRequest("POST", uri, strings.NewReader(json_message))
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", "Bearer " + token)


  if err != nil {
    fmt.Println("Request err:", err)
    return
  }

  client := &http.Client{}
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
