package service

import(
	"errors"
	"library_management/models"
)

type LibraryManager interface{
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
  ReturnBook(bookID int, memberID int) error
  ListAvailableBooks() []models.Book
  ListBorrowedBooks(memberID int) []models.Book

}

type Library struct{
	Books map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() * Library{
	return &Library{
		Books: make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

/*add book*/
func (l*Library) AddBook(book models.Book){
	l.Books[book.ID] = book
}

/*remove book*/
func (l*Library) RemoveBook(bookID int){
	delete(l.Books, bookID)
}

/*borrow book*/
func (l*Library) BorrowBook(bookID int, memberID int) error{
	book, exists := l.Books[bookID]
	if !exists{
		return errors.New("book not found")
	}
	member, exists := l.Members[memberID]
	if!exists{
		return errors.New("member not found")
	}

	if book.Status == "Borrowed"{
		return errors.New("book is already borrowed")
	}
	book.Status = "Borrowed"
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil

}

/*return book*/
func (l*Library) ReturnBook(bookID int, memberID int) error{
	book, exists := l.Books[bookID]
	if !exists{
		return errors.New("book not found")
	}
	member, exists := l.Members[memberID]
	if !exists{
		return errors.New("member not found")
	}

	found:= false
	for i, b := range member.BorrowedBooks{
		if b.ID == bookID{
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
	    found = true
	    break
		}
	}
	if !found{
		return  errors.New("book was not borrowed by this member")
	}
	book.Status = "Available"
	l.Books[bookID] = book
	l.Members[memberID] = member

	return nil
}

/*available books*/
func (l *Library) ListAvailableBooks() []models.Book {
	var available []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}

/*all borrowed books*/
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}