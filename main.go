package main

import (
  "fmt"
  "sync"
  "time"
)

import "main/util"

const DATE_ONLY = "2006-01-02"

func main() {
  currentTime := time.Now()
  endOfLastMonth := currentTime.AddDate(0, 0, -currentTime.Day())

  var wg sync.WaitGroup

  fmt.Println(
    "[Start]",
    "Today is " + currentTime.Format(DATE_ONLY) + ", we'll search exchanges rates from " + endOfLastMonth.Format(DATE_ONLY))

  wg.Add(1)
  go util.FetchOfficialUsdArsRate(endOfLastMonth, func(rate float64, date time.Time) {
    defer wg.Done()
    fmt.Println("[Result] USD to ARS (official):", date.Format(DATE_ONLY), "=", rate)
  })

  wg.Add(1)
  go util.FetchMepUsdArsRate(endOfLastMonth, func(rate float64, date time.Time) {
    defer wg.Done()
    fmt.Println("[Result] USD to ARS (MEP):", date.Format(DATE_ONLY), "=", rate)
  })

  wg.Add(1)
  go util.FetchUvaArsRate(endOfLastMonth, func(rate float64, date time.Time) {
    defer wg.Done()
    fmt.Println("[Result] UVA to ARS:", date.Format(DATE_ONLY), "=", rate)
  })

  wg.Add(1)
  go util.FetchDaiArsRate(endOfLastMonth, func(rate float64, date time.Time) {
    defer wg.Done()
    fmt.Println("[Result] DAI to ARS:", date.Format(DATE_ONLY), "=", rate,)
  })

  wg.Wait()
}
