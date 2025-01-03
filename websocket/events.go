package websocket

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (k *Kraken) handleEvent(msg []byte) error {
	var event EventType
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

  switch event.Channel {
  case ChanOrders:
    return k.handleOrdersEvent(msg, event)
  default:
		log.Warnf("unknown event: %s", msg)
  }

	return nil
}

func (k *Kraken) handleOrdersEvent(msg []byte, event EventType) error {
  switch event.Type {
  case "update":
    return k.handleOrdersUpdateEvent(msg)
  case "snapshot":
    return k.handleOrdersSnapshotEvent(msg)

  default:
    log.Warnf("unknown orders event type: %s", msg)
  }

  return nil
}

func (k *Kraken) handleOrdersSnapshotEvent(msg []byte) error {
  var snapshotMsg OrdersUpdateFullMessage
  if err := json.Unmarshal(msg, &snapshotMsg); err != nil {
    return err
  }

  snapshotMsg.Data[0].IsSnapshot = true

  k.msg <- Update{
    ChannelName: fmt.Sprintf("%s:snapshot", ChanOrders),
    Data:        *snapshotMsg.Data[0],
  }

  return nil
}

func (k *Kraken) handleOrdersUpdateEvent(msg []byte) error {
  var updateMsg OrdersUpdateFullMessage
  if err := json.Unmarshal(msg, &updateMsg); err != nil {
    return err
  }

  k.msg <- Update{
    ChannelName: fmt.Sprintf("%s:update", ChanOrders),
    Data:        *updateMsg.Data[0],
  }

  return nil
}
