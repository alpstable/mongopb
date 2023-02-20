// Copyright 2023 The MongoPB Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

// package mongopb streams structpb-typed data to a MongoDB collection using
// the official MongoDB Go Driver.
package mongopb

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/structpb"
)

// ListWriter is used to write a structpb.ListValue to a MongoDB collection.
type ListWriter struct {
	coll *mongo.Collection
}

// NewListWriter creates a new ListWriter for writing a structpb.ListValue to
// a MongoDB collection.
func NewListWriter(coll *mongo.Collection) *ListWriter {
	listWriter := &ListWriter{coll: coll}

	return listWriter
}

type model struct {
	mongo.WriteModel
}

// mapAsWriteModel converts a map[string]interface{} to a mongo.WriteModel.
func mapAsWriteModel(m map[string]interface{}) model {
	return model{mongo.NewUpdateOneModel().
		SetFilter(m).
		SetUpdate(bson.D{primitive.E{Key: "$set", Value: m}}).
		SetUpsert(true)}
}

// structAsWriteModel converts a structpb.Struct to a mongo.WriteModel.
func structAsWriteModel(spb *structpb.Struct) model {
	return mapAsWriteModel(spb.AsMap())
}

// asWriteModel converts a structpb.Value to a slice of mongo.WriteModel
// objects.
func asWriteModels(value *structpb.Value) <-chan model {
	models := make(chan model)

	go func() {
		defer close(models)

		valType := value.Kind

		obj, ok := valType.(*structpb.Value_StructValue)
		if ok {
			models <- structAsWriteModel(obj.StructValue)
		}
	}()

	return models
}

// Write writes the ListValue to a MongoDB collection.
func (w *ListWriter) Write(ctx context.Context, list *structpb.ListValue) error {
	listLen := len(list.Values)
	if listLen == 0 {
		return nil
	}

	models := []mongo.WriteModel{}

	wg := &sync.WaitGroup{}
	wg.Add(listLen)

	for _, value := range list.GetValues() {
		go func(value *structpb.Value) {
			defer wg.Done()

			for model := range asWriteModels(value) {
				models = append(models, model.WriteModel)
			}
		}(value)
	}

	wg.Wait()

	_, err := w.coll.BulkWrite(ctx, models)
	if err != nil {
		return fmt.Errorf("error writing to MongoDB: %w", err)
	}

	return nil
}
