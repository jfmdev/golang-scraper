package main

import (
  "fmt"
  "time"
  "github.com/gocolly/colly"
)

const DATE_ONLY = "2006-01-02"
const INVESTING_DATE_FORMAT = "02.01.2006"

func main() {
  // --- Initialization --- //

  currentTime := time.Now()
  endOfLastMonth := currentTime.AddDate(0, 0, -currentTime.Day())

  var usdArsRate string
  var usdArsRateDate time.Time

  fmt.Println(
    "[Initializing]",
    "Today is " + currentTime.Format(DATE_ONLY) + ", we'll search exchanges rates from " + endOfLastMonth.Format(DATE_ONLY))

  // --- Collection --- //
  
  coll := colly.NewCollector()

  coll.OnHTML("table[data-test='historical-data-table'] tbody tr[data-test='historical-data-table-row']", func(elem *colly.HTMLElement) {
    // Read row data.
  	rowDateStr := elem.ChildText("td:nth-child(1)")
    rowRateStr := elem.ChildText("td:nth-child(2)")
  	rowDate, rowDateError := time.Parse(INVESTING_DATE_FORMAT, rowDateStr)
  	
    // Check if the row contains more recent data.
  	if rowDateError == nil && 
       rowDate.Month() == endOfLastMonth.Month() && 
       rowDate.After(usdArsRateDate) {
      usdArsRate = rowRateStr
      usdArsRateDate = rowDate
    }
  })

  coll.OnScraped(func(response *colly.Response) {
    fmt.Println("[Finishing] Site scraped", usdArsRateDate, usdArsRate)
  })

  coll.OnRequest(func(request *colly.Request) {
    fmt.Println("[Starting] Visiting", request.URL)
  })

  coll.Visit("https://es.investing.com/currencies/usd-ars-historical-data")
}
