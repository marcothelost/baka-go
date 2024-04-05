package types

type BakalariError struct {
	Error string
}

type AccessInfo struct {
	Access_Token string
	Refresh_Token string
	Token_Type string
	Expires_In int
}

type Mark struct {
	MarkText string
	Weight int
}

type Subject struct {
	Id string
	Abbrev string
	Name string
}

type SubjectListing struct {
	Marks []Mark
	Subject Subject
	AverageText string
}

type MarksListing struct {
	Subjects []SubjectListing
}
