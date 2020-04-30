package api

type endpoint struct {
	ip   string
	port string
}

type boostrapParties struct {
	PartyId   string   `json:"party_id"`
	Endpoint  endpoint `json:"endpoint"`
	PartyType string   `json:"party_type"`
}
type installCluster struct {
	Name            string           `json:"name"`
	Namespace       string           `json:"namespace"`
	Version         string           `json:"version"`
	EggNumber       int              `json:"egg_number"`
	BoostrapParties *boostrapParties `json:"boostrap_parties"`
}
