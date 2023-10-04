package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
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
  var token string
  var delayInSeconds int
  var csvFile string

  flag.StringVar(&clientId, "client_id", "", "Client id")
  flag.StringVar(&clientSecret, "client_secret", "", "Client secret")
  flag.StringVar(&csvFile, "csv", "", "CSV file with email address and messages")
  flag.IntVar(&delayInSeconds, "delay_in_seconds", 60, "Delay in between mails (Optional)")
  flag.StringVar(&token, "token", "", "Token (Optional)")

  flag.Parse()

  if clientSecret == "" || clientId == "" || csvFile == "" {
    fmt.Printf("Error: Required flags are missing %s %s %s !\n", clientId, clientSecret, csvFile)
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
  printUserInformation(token)

  var lastMailTime = time.Now()
  for _, record := range ReadCSVFile(csvFile) {
    var sendTime string
    if record.Timestamp != "" {
      sendTime = record.Timestamp
    } else {
      lastMailTime = lastMailTime.Add(time.Second * time.Duration(delayInSeconds))
      sendTime = lastMailTime.Format(time.RFC3339)
    }
    SendMail(token, record.Email, record.Subject, record.Body, sendTime)
  }
}



