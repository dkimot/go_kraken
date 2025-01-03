package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	ws "github.com/dkimot/go_kraken/websocket"
)

func main() {
  // set up kraken service
  signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	kraken := ws.NewKraken("wss://ws.kraken.com/v2", ws.WithLogLevel(log.InfoLevel))
	if err := kraken.Connect(); err != nil {
		log.Fatalf("Error connecting to web socket: %s", err.Error())
	}

  if err := kraken.Authenticate(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET")); err != nil {
		log.Fatalf("Authenticate error: %s", err.Error())
	}

  // subscribe to BTCUSD`s orders
  // if err := kraken.SubscribeOrders([]string{"BTC/USD"}, ws.Depth10); err != nil {
  //   log.Fatalf("SubscribeOrders error: %s", err.Error())
  // }

  if err := kraken.SubscribeTrades([]string{"BTC/USD"}); err != nil {
    log.Fatalf("SubscribeTrades error: %s", err.Error())
  }

  // 10 - a depth of order book
  // 5 - the price precision from asset info
  // 8 - the volume precision from asset info
  // orderBook := ws.NewOrderBook(10, 5, 8)

	for {
		select {
		case <-signals:
			log.Warn("Stopping...")
			if err := kraken.Close(); err != nil {
				log.Fatal(err)
			}
			return

		case update := <-kraken.Listen():
      switch data := update.Data.(type) {
      case ws.OrdersUpdate: // could be a snapshot or an update
        log.Infof("Orders update: %+v\n", data)
      case []*ws.Trade:
        for _, trade := range data {
          log.Infof("Trade update: %+v\n", trade)
        }
			default:
			}
		}
	}
}
