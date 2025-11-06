package services

import(
	"errors"
	"fmt"
	"library_management/models"
	"sync"
	"time"
)

type LibraryManager interface{
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
  ReturnBook(bookID int, memberID int) error
  ListAvailableBooks() []models.Book
  ListBorrowedBooks(memberID int) []models.Book
	ReserveBook(bookID int, memberID int) error

}

type Library struct{
	Books map[int]models.Book
	Members map[int]models.Member
	mu sync.Mutex
	ReserveQueue chan ReservationRequest
}

func NewLibrary() * Library{
	return &Library{
		Books: make(map[int]models.Book),
		Members: make(map[int]models.Member),
		ReserveQueue: make(chan ReservationRequest, 10),
	}
}

/*add book*/
func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Books[book.ID] = book
}

/*remove book*/
func (l *Library) RemoveBook(bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.Books, bookID)
}

/*borrow book*/
func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	return nil
}


/*return book*/
func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("book was not borrowed by this member")
	}

	book.Status = "Available"
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

/*available books*/
func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

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
	l.mu.Lock()
	defer l.mu.Unlock()

	member, exists := l.Members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

// ---------- New: Concurrent Reservation ----------
type ReservationRequest struct {
	BookID   int
	MemberID int
	Response chan error
}

func (l *Library) ReserveBook(bookID int, memberID int) error {
	req := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Response: make(chan error),
	}

	l.ReserveQueue <- req
	return <-req.Response
}

// Process reservation requests concurrently
func (l *Library) StartReservationWorker() {
	for req := range l.ReserveQueue {
		go l.processReservation(req)
	}
}

func (l *Library) processReservation(req ReservationRequest) {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[req.BookID]
	if !exists {
		req.Response <- errors.New("book not found")
		return
	}

	if book.Status != "Available" {
		req.Response <- errors.New("book is not available for reservation")
		return
	}

	book.Status = "Reserved"
	l.Books[req.BookID] = book
	fmt.Printf("📘 Book '%s' reserved by Member %d\n", book.Title, req.MemberID)

	req.Response <- nil

	// Start auto-cancel timer
	go func(bookID int) {
		time.Sleep(5 * time.Second)

		l.mu.Lock()
		defer l.mu.Unlock()

		book, exists := l.Books[bookID]
		if exists && book.Status == "Reserved" {
			book.Status = "Available"
			l.Books[bookID] = book
			fmt.Printf("⏳ Reservation for book '%s' expired automatically.\n", book.Title)
		}
	}(req.BookID)
}