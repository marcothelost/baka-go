package types

// BakalariError represents an error returned by the Bakalari API.
type BakalariError struct {
	Error string
}

// AccessInfo contains access information including access token, refresh token, token type, and expiration time.
type AccessInfo struct {
	Access_Token string
	Refresh_Token string
	Token_Type string
	Expires_In int
}

// Mark represents a mark for a subject including mark text and weight.
type Mark struct {
	MarkText string
	Weight int
}

// Subject represents a school subject with an ID, abbreviation, and full name.
type Subject struct {
	Id string
	Abbrev string
	Name string
}

// SubjectListing represents a subject along with its marks and average.
type SubjectListing struct {
	Marks []Mark
	Subject Subject
	AverageText string
}

// MarksListing represents a list of subjects with their corresponding marks and averages.
type MarksListing struct {
	Subjects []SubjectListing
}

// FinalMark represents a final mark including mark date, edit date, mark text, subject ID, and ID.
type FinalMark struct {
	MarkDate string
	EditDate string
	MarkText string
	SubjectId string
	Id string
}

// CertificateTerm represents a term in a certificate including final marks, subjects, grade, year in school, school year, repetition status, closure status, and marks average.
type CertificateTerm struct {
	FinalMarks []FinalMark
	Subjects []Subject
	GradeName string
	Grade int
	YearInSchool int
	SchoolYear string
	Repeated bool
	Closed bool
	MarksAverage string
}

// FinalMarksListing represents a list of certificate terms with their final marks and related information.
type FinalMarksListing struct {
	CertificateTerms []CertificateTerm
}
