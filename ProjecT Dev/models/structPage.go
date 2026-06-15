package models

import "html/template"

type Home struct {
	Profil       User
	ListProperty []Property
}

type LegalPageData struct {
	Profil  User
	Title   string
	Content template.HTML
}

type ContactPageData struct {
	Profil  User
	Success bool
	Error   string
}
