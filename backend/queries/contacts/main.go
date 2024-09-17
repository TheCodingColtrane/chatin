package contactQueries

type FindContacts struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Code      string `json:"code"`
}
