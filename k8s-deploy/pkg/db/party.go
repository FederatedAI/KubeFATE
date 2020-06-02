/*
* Copyright 2019-2020 VMware, Inc.
* 
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
* 
*/
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
