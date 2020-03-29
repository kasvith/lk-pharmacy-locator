package main

type PharmacyEntry struct {
	Name           string   `json:"name"`
	District       string   `json:"district"`
	Area           string   `json:"area"`
	Address        string   `json:"address"`
	ContactNo      []string `json:"contact_no"`
	PharmacistName string   `json:"pharmacist_name"`
	Owner          string   `json:"owner"`
	WhatsApp       []string `json:"whatsapp"`
	Viber          []string `json:"viber"`
	Email          []string `json:"email"`
}

type Area struct {
	Name       string
	Pharmacies []*PharmacyEntry
}

type District struct {
	Name  string
	Areas []*Area
}
