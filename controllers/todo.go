package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	Id     int
	Task   string
	IsDone bool
}

var todoList = make(map[int]TodoController)

func isTodoAleadyInList(doesTodoMatch func(todoFromLoop *TodoController) bool) (bool, int) {
	for key, element := range todoList {
		fmt.Println(element, key)
		if doesTodoMatch(&element) {
			return true, key
		}
	}
	return false, -1
}

func addRequestedTodoIntoList(requestedTodo *TodoController) bool {
	listLength := len(todoList)
	doesTodoMatch := func(todoFromLoop *TodoController) bool {
		return requestedTodo.Task == todoFromLoop.Task || requestedTodo.Id == todoFromLoop.Id
	}
	isRequriedTodoInList, _ := isTodoAleadyInList(doesTodoMatch)
	if requestedTodo.Task != "" && !isRequriedTodoInList {
		todoList[listLength] = *requestedTodo
		return true
	}
	return false
}

func updateTodo(updateRequestedTodo *TodoController) bool {
	doesTodoIdMatch := func(todoFromLoop *TodoController) bool {
		return updateRequestedTodo.Id == todoFromLoop.Id
	}
	isUpdateRequestedTodoIntoList, key := isTodoAleadyInList(doesTodoIdMatch)
	if isUpdateRequestedTodoIntoList {
		todoList[key] = *updateRequestedTodo
		return true
	}
	return false
}

/*
 * Route functions ///////////////////////////////
 */

// Create new entry into the todo list
func (c *TodoController) Create(ctx *gin.Context) {
	var requestedTodo TodoController
	ctx.Bind(&requestedTodo)
	if addRequestedTodoIntoList(&requestedTodo) {
		responseString := fmt.Sprintf("Your task %s was added to list", requestedTodo.Task)
		ctx.String(http.StatusOK, responseString)
	} else {
		responseString := fmt.Sprintf("Your task %s already exists in the list!", requestedTodo.Task)
		ctx.String(http.StatusBadRequest, responseString)
	}
}

// Update the todolist
func (c *TodoController) Update(ctx *gin.Context) {
	//todoId := ctx.Param("id")
	var upadtedTodo TodoController
	ctx.Bind(&upadtedTodo)
	if updateTodo(&upadtedTodo) {
		responseString := fmt.Sprintf("Your task %s was updated!")
		ctx.String(http.StatusOK, responseString)
	} else {
		responseString := fmt.Sprintf("Your task was not found")
		ctx.String(http.StatusBadRequest, responseString)
	}
}

func (c *TodoController) Read(ctx *gin.Context) {
	requestedJsonObj, _ := json.Marshal(todoList)
	ctx.String(http.StatusOK, string(requestedJsonObj))

}

func (c *TodoController) ReadById(ctx *gin.Context) {
	requestedId := ctx.Param("id")
	id, _ := strconv.Atoi(requestedId)
	doesTodoIdMatch := func(todoFromLoop *TodoController) bool {
		return id == todoFromLoop.Id
	}
	doesTheRecordExistsInList, key := isTodoAleadyInList(doesTodoIdMatch)
	requestedJsonObj, _ := json.Marshal(todoList[key])
	if doesTheRecordExistsInList {
		ctx.String(http.StatusOK, string(requestedJsonObj))
	} else {
		ctx.String(http.StatusBadRequest, "The todo not found")
	}
}

func (c *TodoController) Remove(ctx *gin.Context) {
	requestedId := ctx.Param("id")
	id, _ := strconv.Atoi(requestedId)
	doesTodoIdMatch := func(todoFromLoop *TodoController) bool {
		return id == todoFromLoop.Id
	}
	doesTheRecordExistsInList, key := isTodoAleadyInList(doesTodoIdMatch)
	if doesTheRecordExistsInList {
		delete(todoList, key)
		responseString := fmt.Sprintf("Your task has been removed successfully")
		ctx.String(http.StatusOK, responseString)
	} else {
		responseString := fmt.Sprintf("Your task was not found")
		ctx.String(http.StatusOK, responseString)
	}

}
