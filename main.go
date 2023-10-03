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

type TokenResponse struct {
  TokenType string `json:"token_type"`
  ExpiresIn int `json:"expires_in"`
  ExtExpiresIn int `json:"ext_expires_in"`
  AccessToken string `json:"access_token"`
}
// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func getToken(clientId string, clientSecret string, tenant string, code string) string {
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

func printUserInformation(token string) {
  uri := "https://graph.microsoft.com/v1.0/me/messages"
  req, err := http.NewRequest("GET", uri, strings.NewReader(""))
  // Set the Content-Type header to indicate JSON data
  // req.Header.Set("Content-Type", "application/json")
  // Set the Content-Type header for form data
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

// func sendTestMail(token string) {
//   uri := "https://graph.microsoft.com/v1.0/me/sendMail"
//   
//   // requestBody, err := json.Marshal(TokenRequest{
//   //   // ClientId: clientId,
//   //   // ClientSecret: clientSecret,
//   //   GrantType: "client_credentials",
//   //   Scope: "https%3A%2F%2Fgraph.microsoft.com%2F.default",
//   // })
//
//   // if err != nil {
//   //   fmt.Println("Not able to parst request")
//   //   return
//   // }
//
//   body := fmt.Sprintf(`
// {
//     "message": {
//       "subject": "Meet for lunch?",
//       "body": {
// 	"contentType": "Text",
// 	"content": "The new cafeteria is open."
//       },
//       "toRecipients": [
//       {
// 	"emailAddress": {
// 	  "address": "frannis@contoso.onmicrosoft.com"
// 	}
//       }
//       ],
//       "singleValueExtendedProperties": [
//       {
// 	"id": "SystemTime 0x3FEF",
// 	"value": "%s"
//       }
//       ],
//       "ccRecipients": [
//       {
// 	"emailAddress": {
// 	  "address": "danas@contoso.onmicrosoft.com"
// 	}
//       }
//       ]
//     },
//     "saveToSentItems": "true"
//   }
//   `, time.Now().Add(time.Minute).UTC().Format(time.RFC3339))
//   fmt.Println(body)
//   req, err := http.NewRequest("POST", uri, strings.NewReader(body))
//   // Set the Content-Type header to indicate JSON data
//   // req.Header.Set("Content-Type", "application/json")
//   // Set the Content-Type header for form data
//   req.Header.Set("Content-Type", "application/json")
//   req.Header.Set("Authorization", "Bearer " + token)
//
//
//   if err != nil {
//     fmt.Println("Request err:", err)
//     return
//   }
//
//   // Create an HTTP client and send the request
//   client := &http.Client{}
//   // reqBody, _ := ioutil.ReadAll(req.Body)
//   // fmt.Println("Request body from request", string(reqBody))
//   resp, err := client.Do(req)
//
//   if err != nil {
//     fmt.Println("Error sending request:", err)
//     return
//   }
//   defer resp.Body.Close()
//   // Read and print the response body
//   respBody, err := ioutil.ReadAll(resp.Body)
//   if err != nil {
//     fmt.Println("Error reading response:", err)
//     return
//   }
//
//
//   // Check the response status code
//   fmt.Println("Response Status Code:", resp.Status)
//   fmt.Println("Response body", string(respBody))
// }

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
  
  // token := "EwB4A8l6BAAUAOyDv0l6PcCVu89kmzvqZmkWABkAAYgW929T/sjSRB4PRsBu2XjdL4rc2V+wQdeGdGN6uO8vNkgTOaz5aU9DtpHQo4HZAaoIdcTDdhnbtaZ261PS28PLvZSk7v1YkcXsJvxqYhAOffSmVOUOhbkG32mUL1n63IJ3bvvrxPoPzuVYYTPEbtoo4fvzqQQYK4pdomCeOgg4fAkDDOrJp2IXUZajjK0VOmRXdQaYsMjZp0/rETe98ehIAdvEtLfg0h41ID93u7WiSKagCGcQO6+an60/4o3NSrgofJQwOl68vgkJ0+VWW3p8TtXqlMZQ55ffL7VzWwJ3OBL6DYdwym7VbmgCEuHgNlBl/0Kg2/sSkdyKPqQvjgEDZgAACPIUwuBjauQZSAKIEvt2Dd9K6gPkz8t/8Cc2CjhDht/Ep2RRKHR3sOE/SFAzj0G/fkoIbclc3JYK/siWDTbFLybB3QChyxA2KB1lZ3NO1IS+T4gRlerNCb6ZYQqlkMpPCpaQ+KfLhzYxOhjXray1ShGYmK57vyZB3WoOIlTrY9uptqMdm0xH+loWF/Zs+JpQQq0KsTxB9uQ/HhTGNaGPHQARIu09U17DLNuBiDXYytzXaS+Ym+QELeXi07C+haIU0TDOTnsz1SUnmV0hbv/icBoagsZ6T7bKPvrTekrBQ0MZLLATVSy9PWUudaSor9ih7WNr4ml/IY/fN3VSUJitKb4D4kPOwhijqrnA49w/Op8Ev8MLUdPyOgn1LV9XDl0/Ljh0j1TqPFrnDkZ5lqcpubsO8qMjj41Bl+a8f9ScP1Ptj70y9S8reDpWRq/VaKM2UcrN10H0n2kK8hpbWICp9/rBhBiMWHcgIcAVr6OzjgUyZHHQujUwMmgiGKD1k7pMIR/8U2vbCD3R+ULpZONFu+iWBwKmTMgjd//NLsVZeC88A11TZ+kGkFC9ddhozYRVHXnV/fOdHqejdhR0zzWHY49PVtfRd0aGvER6zfrH0c1MHk4/gFnWs5DsiZi6bPUwD0mW+33zg74keXCm1wBxH3xYspskG3twwHa167rZs4ucO7U9kX/neN+9HeTWp4mBx0Yq0+v5AuQbqj+DUiEopuwZRrxqTl7gvMCRJC6wdi6Wa/9yCfSv8tOkDUQgeeKT0JsKViK9qUhVyhkWharObzdCrIoC"
  //
  fmt.Println("================")
  fmt.Println("Token: ", token)
  fmt.Println("================")
  // printUserInformation(token)
  SendMail(token, "12tanmayvijay@gmail.com", "Test message from script", "Hello world", time.Now().Format(time.RFC3339))
}



