package repository

import (
	"fmt"
	"github.com/Warh40k/go-todo-rest"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(userId, listId int, list todo.TodoItem) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	var id int
	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createListItemsQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery = fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
								FROM %s ti 
								INNER JOIN %s li ON ti.id = li.item_id 
								INNER JOIN %s ul ON li.list_id = ul.list_id 
							 	WHERE ul.user_id = $1 and li.list_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Select(&items, query, userId, listId)

	return items, err
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
								FROM %s ti 
								INNER JOIN %s li ON ti.id = li.item_id 
								INNER JOIN %s ul ON li.list_id = ul.list_id 
							 	WHERE ul.user_id = $1 and li.item_id = $2 `, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, query, userId, itemId)

	return item, err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s ul, %s li WHERE ti.id = li.item_id and li.list_id = ul.list_id and ul.user_id = $1 and li.item_id = $2", todoItemsTable, usersListsTable, listsItemsTable)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}

//func (r *TodoItemPostgres) Update(userId, listId, itemId int, input todo.UpdateItemInput) error {
//	panic("function not implemented")
//}
