package controllers

import (
	"bufio"
	"fmt"
	"library_management/concurrency"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
	"time"
)

func RunLibraryConsole() {
	library := services.NewLibrary()
	concurrency.StartConcurrentReservationWorker(library)

	library.Members[1] = models.Member{ID: 1, Name: "Hajira"}
	library.Members[2] = models.Member{ID: 2, Name: "Aisha"}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Concurrent Library Management ===")
		fmt.Println("1. Add Book")
		fmt.Println("2. Reserve Book")
		fmt.Println("3. List Available Books")
		fmt.Println("4. Simulate Concurrent Reservations")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(input))

		switch choice {
		case 1:
			fmt.Print("Enter Book ID: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))

			fmt.Print("Enter Title: ")
			title, _ := reader.ReadString('\n')

			fmt.Print("Enter Author: ")
			author, _ := reader.ReadString('\n')

			book := models.Book{
				ID:     id,
				Title:  strings.TrimSpace(title),
				Author: strings.TrimSpace(author),
				Status: "Available",
			}
			library.AddBook(book)
			fmt.Println("Book added successfully!")

		case 2:
			fmt.Print("Enter Book ID to reserve: ")
			bookIDStr, _ := reader.ReadString('\n')
			bookID, _ := strconv.Atoi(strings.TrimSpace(bookIDStr))

			fmt.Print("Enter Member ID: ")
			memberIDStr, _ := reader.ReadString('\n')
			memberID, _ := strconv.Atoi(strings.TrimSpace(memberIDStr))

			err := library.ReserveBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book reserved successfully!")
			}

		case 3:
			fmt.Println("📗 Available Books:")
			for _, book := range library.ListAvailableBooks() {
				fmt.Printf("ID: %d | Title: %s | Status: %s\n", book.ID, book.Title, book.Status)
			}

		case 4:
			fmt.Println("Simulating concurrent reservations...")
			go func() { fmt.Println(library.ReserveBook(1, 1)) }()
			go func() { fmt.Println(library.ReserveBook(1, 2)) }()
			time.Sleep(6 * time.Second)

		case 5:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}