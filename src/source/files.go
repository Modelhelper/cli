package source

type Files struct {
	Connection Connection
	Path       string
}

type fileEntity struct {
	name        string
	schema      string
	description string
	columns     []fileColumn
}

type fileColumn struct {
	name        string
	datatype    string
	nullable    bool
	references  string
	identity    bool
	description string
}
