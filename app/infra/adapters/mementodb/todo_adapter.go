package mementodb

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mecitsemerci/clean-go-todo-api/app/core/domain/todo"
	"github.com/mecitsemerci/clean-go-todo-api/app/infra/utils"
)

type TodoAdapter struct{}

func (adapter *TodoAdapter) GetById(id uuid.UUID) (todo.Todo, error) {
	if existingEntity, ok := todoCollection[id.String()]; ok {
		return existingEntity, nil
	}
	return todo.Todo{}, errors.New("Item not found")
}

func (adapter *TodoAdapter) GetAll() ([]todo.Todo, error) {
	var result []todo.Todo
	for _, v := range todoCollection {
		result = append(result, v)
	}
	return result, nil
}

func (adapter *TodoAdapter) Insert(entity todo.Todo) (uuid.UUID, error) {
	entity.Id = uuid.New()
	entity.CreatedAt = utils.UtcNow()
	entity.UpdatedAt = utils.UtcNow()
	todoCollection[entity.Id.String()] = entity
	return entity.Id, nil
}

func (adapter *TodoAdapter) Update(entity todo.Todo) (bool, error) {

	if existingEntity, ok := todoCollection[entity.Id.String()]; ok {
		updatedEntity := existingEntity
		updatedEntity.Title = entity.Title
		updatedEntity.UpdatedAt = utils.UtcNow()
		todoCollection[entity.Id.String()] = updatedEntity
		return true, nil
	}
	return false, errors.New("Item not found")
}

func (adapter *TodoAdapter) Delete(id uuid.UUID) (bool, error) {

	if _, ok := todoCollection[id.String()]; ok {
		delete(todoCollection, id.String())
		return true, nil
	}
	return false, errors.New("Item not found")
}
