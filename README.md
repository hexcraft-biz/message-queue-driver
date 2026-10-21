# message-queue-driver

Pub/Sub publishing uses `cloud.google.com/go/pubsub/v2` and requires Go 1.25.7 or later.
This Go version requirement applies to all consumers of the module, including
those importing only the `message` package.
The wrapper retains `NewClient`, `Topic`, `Exists`, `Publish`, and `Close`.

`PubsubClient.Entity` is now a v2 `*pubsub.Client`, and `PubsubTopic.Entity` is
now a `*pubsub.Publisher`. Callers that use v1 client or topic types or access
v1-specific methods through `Entity` must update those uses.
`Publish` returns a v2 `*pubsub.PublishResult`. In the currently verified versions
(v1.36.1 and v2.7.0), both packages alias the same underlying result type, so an
existing v1 `PublishResult` declaration alone does not require a source change.
This compatibility should be rechecked when upgrading dependencies.

Create topic wrappers with `client.Topic(...)` so `Exists()` has access to the
associated admin client. Constructing `PubsubTopic` with only `Entity` leaves
that client uninitialized and causes `Exists()` to panic.

## Quick start

Set `GOOGLE_APPLICATION_CREDENTIALS` or configure Application Default Credentials.
The topic must already exist; `Topic` only constructs a publisher reference.

```go
package main

import (
    "context"
    "encoding/json"
    "log"

    "github.com/hexcraft-biz/message-queue-driver/message"
    "github.com/hexcraft-biz/message-queue-driver/pubsub"
)

type Msg struct {
    message.MessagePrototype
}

func (m *Msg) Bytes() []byte {
    data, _ := json.Marshal(m)
    return data
}

func main() {
    client, err := pubsub.NewClient("my-project")
    if err != nil {
        log.Print(err)
        return
    }
    defer client.Close()

    topic := client.Topic("my-topic")
    defer topic.Close()

    msg := &Msg{
        MessagePrototype: message.MessagePrototype{
            Entity: map[string]string{"text": "hello"},
        },
    }
    msg.Init()

    result := topic.Publish(msg, map[string]string{"source": "example"})
    id, err := result.Get(context.Background())
    if err != nil {
        log.Print(err)
        return
    }
    log.Printf("Published message %s", id)
}
```

`Publish` queues a message asynchronously. Use `PublishResult.Get` to wait for
confirmation and retrieve publishing errors. Message JSON retains the `uuid`,
`time`, and `entity` fields; attributes are sent separately.

`Exists` returns `(false, nil)` for a missing topic and returns an error for
other failures, including insufficient permissions.

Reuse each topic publisher. Close every topic before closing the shared client:
`topic.Close()` flushes and stops only that publisher, while `client.Close()`
releases the shared connections. Deferred calls in the example run in this order.
