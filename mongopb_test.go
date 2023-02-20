// Copyright 2023 The MongoPB Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

package mongopb

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/structpb"
)

func testClient(ctx context.Context, t *testing.T, connString string) *mongo.Client {
	t.Helper()

	clientOptions := options.Client().ApplyURI(connString)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("failed to connect to the client: %v", err)
	}

	return client
}

func defaultTestClient(ctx context.Context, t *testing.T) *mongo.Client {
	t.Helper()

	return testClient(ctx, t, fmt.Sprintf("mongodb://mongo1:27017/%s", "defaultdb"))
}

func decodeJSONForTest(t *testing.T, data []byte) (*structpb.ListValue, error) {
	t.Helper()

	// If there is no data, return an empty list.
	if len(data) == 0 {
		return &structpb.ListValue{}, nil
	}

	// Check if the first byte of the json is a '{' or '['
	if data[0] == '{' {
		// Unmarshal the json into a structpb.Struct
		record := &structpb.Struct{}
		if err := json.Unmarshal(data, record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json object: %w", err)
		}

		return &structpb.ListValue{
			Values: []*structpb.Value{
				{
					Kind: &structpb.Value_StructValue{
						StructValue: record,
					},
				},
			},
		}, nil
	}

	records := &structpb.ListValue{}
	if err := json.Unmarshal(data, records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json array: %w", err)
	}

	return records, nil
}

//nolint:paralleltest
func TestListWriter(t *testing.T) {
	// Create a client and ping it.
	ctx := context.Background()

	client := defaultTestClient(ctx, t)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("failed to disconnect from the client: %v", err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("failed to ping the client: %v", err)
	}

	//nolint:paralleltest
	for _, tcase := range []struct {
		name       string
		database   string
		collection string
		data       []byte
	}{
		{
			name:       "empty object",
			database:   "defaultdb",
			collection: "test1",
			data:       []byte(`{}`),
		},
		{
			name:       "simple object",
			database:   "defaultdb",
			collection: "test2",
			data:       []byte(`{"name": "test1"}`),
		},
		{
			name:       "empty array",
			database:   "defaultdb",
			collection: "test3",
			data:       []byte(`[]`),
		},
		{
			name:       "simple array",
			database:   "defaultdb",
			collection: "test4",
			data:       []byte(`[{"name": "test1"}, {"name": "test2"}]`),
		},
		{
			name:       "nested object",
			database:   "defaultdb",
			collection: "test5",
			data:       []byte(`{"foo": {"bar": "baz"}}`),
		},
		{
			name:       "nested array",
			database:   "defaultdb",
			collection: "test6",
			data:       []byte(`{"foo": [{"bar": "baz"}]}`),
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			// Create a new ListWriter.
			coll := client.Database(tcase.database).Collection(tcase.collection)
			listWriter := NewListWriter(coll)

			// Defer to drop the collection.
			defer func() {
				if err := coll.Drop(ctx); err != nil {
					t.Fatalf("failed to drop the collection: %v", err)
				}
			}()

			// Decode the json data.
			list, err := decodeJSONForTest(t, tcase.data)
			if err != nil {
				t.Fatalf("failed to decode json: %v", err)
			}

			// Write the ListValue to the collection.
			if err := listWriter.Write(ctx, list); err != nil {
				t.Fatalf("failed to write the list value: %v", err)
			}

			// Check to see if the list was correctly written to the
			// collection.
			for _, value := range list.Values {
				// Turn the value into bytes.
				bytes, err := json.Marshal(value)
				if err != nil {
					t.Fatalf("failed to marshal the value: %v", err)
				}

				// Unmarshal the bytes into a bson.D
				var doc bson.D
				if err := bson.UnmarshalExtJSON(bytes, true, &doc); err != nil {
					t.Fatalf("failed to unmarshal the bytes: %v", err)
				}

				// Find the document in the collection.
				result := coll.FindOne(ctx, doc)
				if result.Err() != nil {
					t.Fatalf("failed to find the document: %v", result.Err())
				}

				// Decode the result into a bson.D
				var found bson.D
				if err := result.Decode(&found); err != nil {
					t.Fatalf("failed to decode the result: %v", err)
				}

				// Check to see if the document was found.
				if len(found) == 0 {
					t.Fatalf("failed to find the document")
				}
			}
		})
	}
}
