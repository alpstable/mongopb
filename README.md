# GMongo

Gidari Mongo is the MongoDB storage implementation for Gidari.

## Usage

```go
package main

import (
	"context"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/config"
	"github.com/alpstable/gmongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.TODO()

	// Create a MongoDB client using the official MongoDB Go Driver.
	clientOptions := options.Client().ApplyURI("mongodb://mongo1:27017/defaultcoll")
	client, _ := mongo.Connect(ctx, clientOptions)

	// Plug the client into a Gidari MongoDB Storage adapater.
	mdbStorage, _ := gmongo.New(ctx, client)

	// Include the adapter in the storage slice of the transport configuration.
	err := gidari.Transport(ctx, &config.Config{
		Storage: []Storage{mdbStorage},
	})
}
```
