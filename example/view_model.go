package example

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
