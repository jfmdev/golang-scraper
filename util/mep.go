package util

import (
  "context"
  "fmt"
  "github.com/chromedp/cdproto/cdp"
  "github.com/chromedp/chromedp"
  "log"
  "strconv"
  "strings"
  "time"
)

func FetchMepUsdArsRate(deadline time.Time, callback func(float64, time.Time)) {
  const TARGET_URL = "https://www.ambito.com/contenidos/dolar-mep-historico.html"
  const DATE_FORMAT = "02/01/2006"

  const TABLE_BODY_SELECTOR = "table.general-historical__table tbody.general-historical__tbody"
  const ROWS_SELECTOR = TABLE_BODY_SELECTOR + " tr"
  const CELL_SELECTOR = ROWS_SELECTOR + ":nth-child(%d) td:nth-child(%d)"

  var rateValue float64
  var rateDate time.Time

  // Create Chrome instance.
  ctx, cancel := chromedp.NewContext(
    context.Background(),
    chromedp.WithLogf(log.Printf),
  )
  defer cancel()

  // Create a timeout.
  ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
  defer cancel()

  // Navigate to a page, wait for table to load and read rows.
  var nodes []*cdp.Node
  err := chromedp.Run(ctx,
    chromedp.Navigate(TARGET_URL),
    chromedp.WaitVisible(TABLE_BODY_SELECTOR),
    chromedp.Nodes(ROWS_SELECTOR, &nodes),
  )
  if err != nil {
    log.Fatal(err)
  }

  // Parse rows of table.
  for i := 0; i < len(nodes); i++ {
    if nodes[i].ChildNodeCount >= 2 {
      var rowDateStr, rowValueStr string
      if err := chromedp.Run(ctx,
        chromedp.Text(fmt.Sprintf(CELL_SELECTOR, i+1, 1), &rowDateStr),
        chromedp.Text(fmt.Sprintf(CELL_SELECTOR, i+1, 2), &rowValueStr),
      ); err != nil {
        log.Fatal(err)
      }

      rowValue, rowValueErr := strconv.ParseFloat(strings.Replace(rowValueStr, ",", ".", -1), 64);
      rowDate, rowDateError := time.Parse(DATE_FORMAT, rowDateStr)

      // Check if the row contains more recent data.
      if rowDateError == nil && 
         rowValueErr == nil &&
         rowDate.Month() == deadline.Month() && 
         rowDate.After(rateDate) {
        rateValue = rowValue
        rateDate = rowDate
      }
    }
  }

  callback(rateValue, rateDate)
}
