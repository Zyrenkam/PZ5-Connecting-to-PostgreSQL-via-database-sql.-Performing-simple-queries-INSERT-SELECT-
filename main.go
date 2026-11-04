package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен; если файла нет — ошибка игнорируется
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback — прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:root@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// Проверочное 1
	doneTasks, err := repo.ListDone(ctx, false)
	if err != nil {
		log.Fatal("ListDone error: %v", err)
	}
	fmt.Println("Not done tasks")
	for _, t := range doneTasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// Проверочное 2
	task, err := repo.FindByID(ctx, 2)
	if err != nil {
		log.Fatal("FindByID error: %v", err)
	}
	fmt.Printf("Task #%d \n", task.ID)
	fmt.Printf("Title: %s\nDone: %v\nCreated: %s\n",
		task.Title, task.Done, task.CreatedAt.Format(time.RFC3339))

	// Проверочное 3
	more := []string{"CreateMany #1", "CreateMany #2", "CreateMany #3"}
	if err := repo.CreateMany(ctx, more); err != nil {
		log.Fatal("CreateMany error: %v", err)
	}
	log.Println("CreateMany: inserted 3 tasks in one transaction")
}
