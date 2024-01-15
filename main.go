package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3" //_ solely for importing packages for its side effects
	"github.com/urfave/cli/v2"
)

type Tasks struct {
	taskID   int
	taskName string
}

// create a new database
func createDB(db *sql.DB) error {
	const createTable string = `CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER NOT NULL,
		task TEXT
	);`

	_, err := db.Exec(createTable)
	if err != nil {
		return err
	}
	return nil
}

//show database
func showDB(db *sql.DB) error {
	const showDB string = `SELECT * FROM tasks;`
	rows, err := db.Query(showDB)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var task Tasks
		err := rows.Scan(&task.taskID, &task.taskName)
		if err != nil {
			return err
		}
		fmt.Println(task.taskID, task.taskName)
	}
	return nil
}

//insert to database
func insertToDB(db *sql.DB, taskID int, task string) error {
	const insertToDB string = `INSERT INTO tasks VALUES(?,?);`
	_, err := db.Exec(insertToDB, taskID, task)
	if err != nil {
		return err
	}
	return nil
}

//delete from database
func deleteFromDB(db *sql.DB, taskID int) error {
	const deleteFromDB string = `DELETE FROM tasks WHERE id=?;`
	_, err := db.Exec(deleteFromDB, taskID)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	const file string = "tasks.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("Connection failed")
		log.Fatal(err)
	}

	err = createDB(db)
	if err != nil {
		fmt.Println("Failed to create table tasks")
		log.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "show",
				Aliases: []string{"s"},
				Usage:   "show the database",
				Action: func(ctx *cli.Context) error {
					err := showDB(db)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "insert",
				Aliases: []string{"i"},
				Usage:   "insert a task to database (args: taskID, taskName)",
				Action: func(ctx *cli.Context) error {
					taskID, _ := strconv.Atoi(ctx.Args().Get(0))
					taskName := ctx.Args().Get(1)
					err := insertToDB(db, taskID, taskName)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "delete a task from database (args: taskID)",
				Action: func(ctx *cli.Context) error {
					taskID, _ := strconv.Atoi(ctx.Args().Get(0))
					err := deleteFromDB(db, taskID)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
