package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mecitsemerci/go-todo-app/internal/api/dto"
	"github.com/mecitsemerci/go-todo-app/internal/api/httperror"
	"github.com/mecitsemerci/go-todo-app/internal/core/domain"
	"github.com/mecitsemerci/go-todo-app/internal/core/interfaces"
	"github.com/mecitsemerci/go-todo-app/internal/pkg/check"
	"github.com/mecitsemerci/go-todo-app/internal/pkg/validator"
	"net/http"
)

type TodoController struct {
	TodoService interfaces.TodoService
}

func NewTodoController(todoService interfaces.TodoService) TodoController {
	return TodoController{TodoService: todoService}
}

func (controller *TodoController) Register(apiRouteGroup *gin.RouterGroup) {
	router := apiRouteGroup.Group("/todo")
	{
		router.GET("/", controller.GetAll)
		router.GET("/:id", controller.Find)
		router.POST("/", controller.Create)
		router.PUT("/:id", controller.Update)
		router.DELETE("/:id", controller.Delete)
	}
}

// GetAll godoc
// @Summary Get all todo
// @Description Get all todo array
// @Tags Todo
// @Accept json
// @Produce json
// @Success 200 {array} dto.TodoOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 422 {object} dto.ErrorOutput
// @Router /api/v1/todo [get]
func (controller *TodoController) GetAll(ctx *gin.Context) {
	todoList, err := controller.TodoService.GetAll()
	if err != nil {
		httperror.NewError(ctx, http.StatusUnprocessableEntity, "Something went wrong!", err)
		return
	}
	var result = make([]dto.TodoOutput, 0)
	todoOutputDto := dto.TodoOutput{}
	for _, model := range todoList {
		result = append(result, todoOutputDto.FromModel(model))
	}
	ctx.JSON(http.StatusOK, result)
}

// Find Todo godoc
// @Summary Find a todo
// @Description Get todo by id
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.TodoOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 404 {object} dto.ErrorOutput
// @Router /api/v1/todo/{id} [get]
func (controller *TodoController) Find(ctx *gin.Context) {
	todoId := ctx.Param("id")

	if check.IsEmptyOrWhiteSpace(todoId) {
		httperror.NewError(ctx, http.StatusBadRequest, "ID is empty.", nil)
		return
	}

	todoModel, err := controller.TodoService.Find(domain.ID(todoId))

	if err != nil {
		httperror.NewError(ctx, http.StatusUnprocessableEntity, "Item is not exist.", err)
		return
	}

	todoOutput := dto.TodoOutput{}

	ctx.JSON(http.StatusOK, todoOutput.FromModel(*todoModel))
}

// Create Todo godoc
// @Summary Create a todo
// @Description add by json todo
// @Tags Todo
// @Accept  json
// @Produce  json
// @Param todo body dto.CreateTodoInput true "Create todo"
// @Success 200 {object} dto.CreateTodoOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 422 {object} dto.ErrorOutput
// @Router /api/v1/todo [post]
func (controller *TodoController) Create(ctx *gin.Context) {
	var createTodoInput dto.CreateTodoInput

	if err := ctx.ShouldBindJSON(&createTodoInput); err != nil {
		httperror.NewError(ctx, http.StatusBadRequest, "Request model is invalid.", err)
		return
	}

	if err := validator.Validate(createTodoInput); err != nil {
		httperror.NewError(ctx, http.StatusBadRequest, "Validation error", err)
		return
	}

	todoId, err := controller.TodoService.Create(createTodoInput.ToModel())

	if err != nil {
		httperror.NewError(ctx, http.StatusUnprocessableEntity, "The item failed to create.", err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.CreateTodoOutput{
		TodoId: todoId.String(),
	})

}

// Update Todo godoc
// @Summary Update a todo
// @Description update by json todo
// @Tags Todo
// @Accept  json
// @Produce  json
// @Param id path string true "Todo ID"
// @Param todo body dto.UpdateTodoInput true "Update todo"
// @Success 204 {object} dto.TodoOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 404 {object} dto.ErrorOutput
// @Router /api/v1/todo/{id} [put]
func (controller *TodoController) Update(ctx *gin.Context) {
	todoId := ctx.Param("id")

	if check.IsEmptyOrWhiteSpace(todoId) {
		httperror.NewError(ctx, http.StatusBadRequest, "ID is invalid.", nil)
		return
	}

	var updateTodoInput dto.UpdateTodoInput

	if err := ctx.ShouldBindJSON(&updateTodoInput); err != nil {
		httperror.NewError(ctx, http.StatusBadRequest, "Request model is invalid.", err)
		return
	}

	if err := validator.Validate(updateTodoInput); err != nil {
		httperror.NewError(ctx, http.StatusBadRequest, "Validation error", err)
		return
	}

	model := updateTodoInput.ToModel(todoId)

	err := controller.TodoService.Update(*model)

	if err != nil {
		httperror.NewError(ctx, http.StatusNotFound, "The item failed to update.", err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

// Delete Todo godoc
// @Summary Delete a todo
// @Description Delete by todo id
// @Tags Todo
// @Accept  json
// @Produce  json
// @Param id path string true "Todo ID"
// @Success 204 {object} dto.TodoOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 404 {object} dto.ErrorOutput
// @Router /api/v1/todo/{id} [delete]
func (controller *TodoController) Delete(ctx *gin.Context) {
	todoId := ctx.Param("id")

	if check.IsEmptyOrWhiteSpace(todoId) {
		httperror.NewError(ctx, http.StatusBadRequest, "ID is empty.", nil)
		return
	}

	err := controller.TodoService.Delete(domain.ID(todoId))

	if err != nil {
		httperror.NewError(ctx, http.StatusNotFound, "The item failed to delete.", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}