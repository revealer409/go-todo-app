package mongodb

import (
	"errors"
	"github.com/mecitsemerci/clean-go-todo-api/app/core/domain"
	"github.com/mecitsemerci/clean-go-todo-api/app/core/domain/todo"
	"github.com/mecitsemerci/clean-go-todo-api/app/infra/datetime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type TodoAdapter struct {
	DbCtx DbContext
}

func NewTodoAdapter(dbContext DbContext) *TodoAdapter {
	return &TodoAdapter{DbCtx: dbContext}
}

func (adapter *TodoAdapter) GetAll() ([]*todo.Todo, error) {
	var todos []*todo.Todo

	//Configure Find Query
	findOptions := options.Find()

	//Connect
	adapter.DbCtx.Connect()

	cur, err := adapter.DbCtx.TodoCollection.Find(adapter.DbCtx.Context, bson.D{}, findOptions)

	if err != nil {
		return nil, err
	}

	for cur.Next(adapter.DbCtx.Context) {
		var entity Todo
		err := cur.Decode(&entity)
		if err != nil {
			log.Println(err.Error())
		}
		todos = append(todos, entity.ToModel())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	_ = cur.Close(adapter.DbCtx.Context)

	//Disconnect
	defer adapter.DbCtx.Disconnect()

	return todos, nil
}

func (adapter *TodoAdapter) GetById(id domain.ID) (*todo.Todo, error) {
	var t Todo
	oid, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return nil, err
	}
	//Filter
	filter := bson.M{"_id": bson.M{"$eq": oid}}

	//Connect
	adapter.DbCtx.Connect()
	//Find Item by Id
	err = adapter.DbCtx.TodoCollection.FindOne(adapter.DbCtx.Context, filter).Decode(&t)

	//Disconnect
	defer adapter.DbCtx.Disconnect()

	if err != nil {
		return nil, err
	}

	return t.ToModel(), nil
}

func (adapter *TodoAdapter) Insert(todo todo.Todo) (domain.ID, error) {
	// Data object
	var t Todo
	err := t.FromModel(&todo)
	if err != nil {
		return domain.NilID, err
	}
	//Connect
	adapter.DbCtx.Connect()

	//Insert Item
	result, err := adapter.DbCtx.TodoCollection.InsertOne(adapter.DbCtx.Context, t)

	//Disconnect
	defer adapter.DbCtx.Disconnect()

	if err != nil {
		return domain.NilID, err
	}

	// Return inserted item id
	var oid ObjectId
	oid.Set(result.InsertedID.(primitive.ObjectID).Hex())
	return &oid, nil
}

func (adapter *TodoAdapter) Update(todo todo.Todo) error {
	oid, err := primitive.ObjectIDFromHex(todo.Id.String())
	if err != nil {
		return err
	}
	//Filter
	filter := bson.M{"_id": bson.M{"$eq": oid}}

	//Update fields
	document := bson.M{"$set": bson.M{
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
		"updated_at":  datetime.Now(),
	}}

	//Connect
	adapter.DbCtx.Connect()

	//Update Item
	result, err := adapter.DbCtx.TodoCollection.UpdateOne(adapter.DbCtx.Context, filter, document)

	//Disconnect
	defer adapter.DbCtx.Disconnect()

	if err != nil {
		return err
	}

	if result.MatchedCount > 0 {
		return nil
	}
	return errors.New("no items have been updated")
}

func (adapter *TodoAdapter) Delete(id domain.ID) error {
	oid, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return err
	}

	//Filter
	filter := bson.M{"_id": bson.M{"$eq": oid}}

	//Connect
	adapter.DbCtx.Connect()

	//Delete Item
	result, err := adapter.DbCtx.TodoCollection.DeleteOne(adapter.DbCtx.Context, filter)

	//Disconnect
	defer adapter.DbCtx.Disconnect()

	if err != nil {
		return err
	}

	if result.DeletedCount > 0 {
		return nil
	}
	return errors.New("no item has been deleted")
}
