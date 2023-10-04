package main

import (
  "fmt"
  "encoding/csv"
  "os"
)

type UserEmailMessage struct {
  Email string
  Subject string
  Body string
  Timestamp string
}

func validateStruct(s UserEmailMessage) bool {
  return true 
}

func parseUserEmailMessage(data [][] string) []UserEmailMessage {
  var result []UserEmailMessage

  for row_num, row := range data {
    if row_num > 0 {
      var record UserEmailMessage
      if len(row) > 2 {
	record.Email = row[0]	
	record.Subject = row[1]
	record.Body = row[2]
	record.Timestamp = ""
	if len(row) > 3 {
	  record.Timestamp = row[3]
	}
	if validateStruct(record) {
	  result = append(result, record)
	}
      } else {
	fmt.Printf("Skipping %d: %s", row_num, row)
      }
    }
  }
  return result
}

func ReadCSVFile(filepath string) []UserEmailMessage {
  file, err := os.Open(filepath)
  if err != nil {
    fmt.Println("Error open file", err)
  }
  defer file.Close()

  csvReader := csv.NewReader(file)
  data, err := csvReader.ReadAll()

  if err != nil {
    fmt.Println("Error reading csv", err)
  }

  return parseUserEmailMessage(data)

}
