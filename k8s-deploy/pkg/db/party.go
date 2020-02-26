package db

type Party struct {
	PartyId   string `json:"party_id"`
	Endpoint  string `json:"endpoint"`
	PartyType string `json:"party_type"`
}

func NewParty(partyId string, endpoint string, partyType string) *Party {
	party := &Party{
		PartyId:   partyId,
		Endpoint:  endpoint,
		PartyType: partyType,
	}

	return party
}
