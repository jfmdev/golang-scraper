package util

import (
  "strconv"
  "strings"
  "time"
  "github.com/gocolly/colly"
)

func FetchOfficialUsdArsRate(deadline time.Time, callback func(float64, time.Time)) {
  const TARGET_URL = "https://es.investing.com/currencies/usd-ars-historical-data"
  const DATE_FORMAT = "02.01.2006"

  var rateValue float64
  var rateDate time.Time

  coll := colly.NewCollector()

  coll.OnHTML("table[data-test='historical-data-table'] tbody tr[data-test='historical-data-table-row']", func(elem *colly.HTMLElement) {
    // Read row data.
    rowDateStr := elem.ChildText("td:nth-child(1)")
    rowValueStr := elem.ChildText("td:nth-child(2)")
    rowDate, rowDateError := time.Parse(DATE_FORMAT, rowDateStr)
    rowValue, rowValueErr := strconv.ParseFloat(strings.Replace(rowValueStr, ",", ".", -1), 64);
    
    // Check if the row contains more recent data.
    if rowDateError == nil && 
       rowValueErr == nil && 
       rowDate.Month() == deadline.Month() && 
       rowDate.After(rateDate) {
      rateValue = rowValue
      rateDate = rowDate
    }
  })

  coll.OnScraped(func(response *colly.Response) {
    callback(rateValue, rateDate)
  })

  coll.Visit(TARGET_URL)
}
