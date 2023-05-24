package util

import (
  "context"
  "encoding/json"
  "github.com/chromedp/cdproto/cdp"
  "github.com/chromedp/cdproto/network"
  "github.com/chromedp/chromedp"
  "log"
  "strconv"
  "strings"
  "time"
)

type DaiQuote map[string]string

func parseQuoteRequest(
  requestID network.RequestID, 
  ctxBase context.Context, 
  deadline time.Time,
  callback func(float64, time.Time),
) {
  const DATE_FORMAT = "2006-01-02"

  var rateValue float64
  var rateDate time.Time

  // Get response.
  c := chromedp.FromContext(ctxBase)
  ctx := cdp.WithExecutor(ctxBase, c.Target)
  
  bodyBytes, bodyErr := network.GetResponseBody(requestID).Do(ctx)
  if bodyErr != nil {
    log.Printf("get body err: %s", bodyErr)
    return
  }

  var bodyObj []DaiQuote
  errUnmarshal := json.Unmarshal(bodyBytes, &bodyObj)

  // Iterate quotes.
  if errUnmarshal == nil {
    for i := 0; i < len(bodyObj); i++ {
      date, dateError := time.Parse(DATE_FORMAT, bodyObj[i]["date"])
      buyRate, buyRateErr := strconv.ParseFloat(bodyObj[i]["buy_rate"], 64)
      sellRate, sellRateErr := strconv.ParseFloat(bodyObj[i]["sell_rate"], 64)

      // Check if the entry contains more recent data.
      if dateError == nil && 
         buyRateErr == nil && 
         sellRateErr == nil && 
         date.Month() == deadline.Month() && 
         date.After(rateDate) {
          rateValue = (buyRate + sellRate)/2
          rateDate = date
      }
    }  
  }

  callback(rateValue, rateDate)
}

func FetchDaiArsRate(deadline time.Time, callback func(float64, time.Time)) {
  const TARGET_URL = "https://www.ripio.com/ar/dai/"

  var requestIDofQuote network.RequestID

  // Create Chrome instance.
  ctx, cancel := chromedp.NewContext(
    context.Background(),
    chromedp.WithLogf(log.Printf),
  )
  defer cancel()

  // Create a timeout.
  ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
  defer cancel()

  // List events.
  chromedp.ListenTarget(ctx, func(ev interface{}) {
    switch ev := ev.(type) {

    // Identify ID of request fetching the list of quotes.
    case *network.EventResponseReceived:
      if ev.Type == "Fetch" || ev.Type == "XHR" {
        resp := ev.Response

        if strings.Contains(resp.URL, "quote") {
          requestIDofQuote = ev.RequestID
        }
      }

    // Once request has finished, get list of quotes.
    case *network.EventLoadingFinished:
      if ev.RequestID == requestIDofQuote {
        go parseQuoteRequest(ev.RequestID, ctx, deadline, callback)
      }
    }
  })

  // Navigate to a page, wait for sections to be rendered.
  runErr := chromedp.Run(ctx,
    chromedp.Navigate(TARGET_URL),
    chromedp.WaitVisible("section"),
  )
  if runErr != nil {
    log.Fatal(runErr)
  }
}
