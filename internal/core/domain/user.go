package domain

type Geo struct {
	Latitude  float64
	Longitude float64
}

type Address struct {
	Street  string
	Suite   string
	City    string
	Zipcode string
	Geo     Geo
}

type Company struct {
	Name        string
	CatchPhrase string
	BS          string
}

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Website  string
	Address  Address
	Company  Company
}
