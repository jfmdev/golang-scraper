package util

import (
  // "fmt"
  "strconv"
  "strings"
  "time"
  "github.com/gocolly/colly"
)

const INVESTING_DATE_FORMAT = "02.01.2006"

func FetchOfficialUsdArsRate(deadline time.Time, callback func(float64, time.Time)) {
  var usdArsRateStr string
  var usdArsRateDate time.Time

  coll := colly.NewCollector()

  coll.OnHTML("table[data-test='historical-data-table'] tbody tr[data-test='historical-data-table-row']", func(elem *colly.HTMLElement) {
    // Read row data.
    rowDateStr := elem.ChildText("td:nth-child(1)")
    rowRateStr := elem.ChildText("td:nth-child(2)")
    rowDate, rowDateError := time.Parse(INVESTING_DATE_FORMAT, rowDateStr)
    
    // Check if the row contains more recent data.
    if rowDateError == nil && 
       rowDate.Month() == deadline.Month() && 
       rowDate.After(usdArsRateDate) {
      usdArsRateStr = rowRateStr
      usdArsRateDate = rowDate
    }
  })

  coll.OnScraped(func(response *colly.Response) {
    usdArsRate, usdArsRateErr := strconv.ParseFloat(strings.Replace(usdArsRateStr, ",", ".", -1), 64);
    if usdArsRateErr != nil {
      usdArsRate = -1
    }

    // fmt.Println("[Finishing] Site scraped", usdArsRateDate, usdArsRateStr, usdArsRate)
    callback(usdArsRate, usdArsRateDate)
  })

  // coll.OnRequest(func(request *colly.Request) {
  //   fmt.Println("[Starting] Visiting", request.URL)
  // })

  coll.Visit("https://es.investing.com/currencies/usd-ars-historical-data")
}
