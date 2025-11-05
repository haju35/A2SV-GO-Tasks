package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

func RunLibraryConsole() {
	library := service.NewLibrary()

	library.Members[1] = models.Member{ID: 1, Name: "Hajira"}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Library Management System ===")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Exit")
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
			fmt.Print("Enter Book ID to remove: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			library.RemoveBook(id)
			fmt.Println("Book removed successfully!")

		case 3:
			fmt.Print("Enter Book ID to borrow: ")
			bookIDStr, _ := reader.ReadString('\n')
			bookID, _ := strconv.Atoi(strings.TrimSpace(bookIDStr))

			fmt.Print("Enter Member ID: ")
			memberIDStr, _ := reader.ReadString('\n')
			memberID, _ := strconv.Atoi(strings.TrimSpace(memberIDStr))

			err := library.BorrowBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book borrowed successfully!")
			}

		case 4:
			fmt.Print("Enter Book ID to return: ")
			bookIDStr, _ := reader.ReadString('\n')
			bookID, _ := strconv.Atoi(strings.TrimSpace(bookIDStr))

			fmt.Print("Enter Member ID: ")
			memberIDStr, _ := reader.ReadString('\n')
			memberID, _ := strconv.Atoi(strings.TrimSpace(memberIDStr))

			err := library.ReturnBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book returned successfully!")
			}

		case 5:
			fmt.Println("Available Books:")
			for _, book := range library.ListAvailableBooks() {
				fmt.Printf("ID: %d | Title: %s | Author: %s\n", book.ID, book.Title, book.Author)
			}

		case 6:
			fmt.Print("Enter Member ID: ")
			memberIDStr, _ := reader.ReadString('\n')
			memberID, _ := strconv.Atoi(strings.TrimSpace(memberIDStr))

			fmt.Println("Borrowed Books:")
			for _, book := range library.ListBorrowedBooks(memberID) {
				fmt.Printf("ID: %d | Title: %s | Author: %s\n", book.ID, book.Title, book.Author)
			}

		case 7:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
