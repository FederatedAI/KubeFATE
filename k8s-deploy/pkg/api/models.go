package api

type endpoint struct {
	ip   string
	port string
}

type boostrapParties struct {
	partyId   string   `json:"party_id"`
	endpoint  endpoint `json:"endpoint"`
	partyType string   `json:"party_type"`
}
type installCluster struct {
	Name            string           `json:"name"`
	Namespace       string           `json:"namespace"`
	Version         string           `json:"version"`
	EggNumber       int              `json:"egg_number"`
	BoostrapParties *boostrapParties `json:"boostrap_parties"`
}
