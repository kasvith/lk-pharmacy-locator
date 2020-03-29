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

type AreaData map[string][]*PharmacyEntry
type DistrictData map[string]AreaData

type DistrictsResponse struct {
	Districts []string `json:"districts"`
}

type AreasResponse struct {
	Areas []string `json:"areas"`
}

type PharmacyResponse struct {
	Pharmacies []*PharmacyEntry `json:"pharmacies"`
}
