package example

// User is user model
type User struct {
	ID       string
	Name     string
	Contacts map[string]*Contact
}

// Contact is contact model
type Contact struct {
	ID     string
	Email  string
	UserID string
}

// UserMap is user map
type UserMap = map[string]*User

// UserView is user view
type UserView struct {
	ID       string
	Name     string
	Contacts []*ContactView
}

// ContactView is contact view
type ContactView struct {
	ID    string
	Email string
}
