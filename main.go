package main

import (
  "fmt"
  "time"
)

import "main/util"

const DATE_ONLY = "2006-01-02"

func main() {
  currentTime := time.Now()
  endOfLastMonth := currentTime.AddDate(0, 0, -currentTime.Day())

  fmt.Println(
    "[Start]",
    "Today is " + currentTime.Format(DATE_ONLY) + ", we'll search exchanges rates from " + endOfLastMonth.Format(DATE_ONLY))

  util.FetchOfficialUsdArsRate(endOfLastMonth, func(rate float64, date time.Time) {
    fmt.Println("[Result] USD to ARS (official)", rate, date.Format(DATE_ONLY))
  })
}
