package main

import "fmt"

type Manageable interface {
	Borrow() bool
	Return() bool
	GetInfo() string
}

type Book struct {
	ID          int
	Title       string
	Author      string
	IsAvailable bool
}

type Magazine struct {
	ID          int
	Title       string
	Issue       int
	IsAvailable bool
}

type Library struct {
	Books     []*Book
	Magazines []*Magazine
	Name      string
}

func (b *Book) Borrow() bool {
	if b.IsAvailable {
		b.IsAvailable = true
		return true
	}
	return false
}

func (b *Book) Return() bool {
	if !b.IsAvailable {
		b.IsAvailable = true
		return true
	}
	return false
}

func (b *Book) GetInfo() string {
	return fmt.Sprintf("Book ID: %d, Title: %s, Author: %s, Available: %t", b.ID, b.Title, b.Author, b.IsAvailable)
}

func (m *Magazine) Borrow() bool {
	if m.IsAvailable {
		m.IsAvailable = true
		return true
	}
	return false
}
func (m *Magazine) Return() bool {
	if !m.IsAvailable {
		m.IsAvailable = true
		return true
	}
	return false
}
func (m *Magazine) GetInfo() string {
	return fmt.Sprintf("Magazine ID: %d, Title: %s, Issue: %s, Available: %t", m.ID, m.Title, m.Issue, m.IsAvailable)
}

func (l *Library) AddBook(b *Book) {
	l.Books = append(l.Books, b)
}

func (l *Library) AddMagazine(m *Magazine) {
	l.Magazines = append(l.Magazines, m)
}

func (l *Library) ShowAvailableItems() {
	fmt.Println("Available items in the library:", l.Name)
	fmt.Println("Books:")
	for _, book := range l.Books {
		if book.IsAvailable {
			fmt.Println(book.GetInfo())
		}
	}
	fmt.Println("Magazines:")
	for _, magazine := range l.Magazines {
		if magazine.IsAvailable {
			fmt.Println(magazine.GetInfo())
		}
	}
}

func main() {
	library := Library{Name: "Chinese Library"}

	book1 := &Book{ID: 1, Title: "Language", Author: "wyz", IsAvailable: true}
	book2 := &Book{ID: 1, Title: "Language1", Author: "wz", IsAvailable: true}
	book3 := &Book{ID: 1, Title: "Language2", Author: "wyz11", IsAvailable: true}

	library.AddBook(book1)
	library.AddBook(book2)
	library.AddBook(book3)

	magazine1 := &Magazine{ID: 1, Title: "English", Issue: 11, IsAvailable: true}
	magazine2 := &Magazine{ID: 2, Title: "English1", Issue: 22, IsAvailable: true}

	library.AddMagazine(magazine1)
	library.AddMagazine(magazine2)

	library.ShowAvailableItems()

	fmt.Println("借书")
	book1.Borrow()

	fmt.Println("借杂志")
	magazine1.Borrow()

	fmt.Println("还书")
	book1.Return()

	library.ShowAvailableItems()
}
