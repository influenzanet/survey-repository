package models

type PlatformInfo struct {
	ID string `json:"id"`
	Label string `json:"label"`
	Country string `json:"country"`
}

var WellKnownPlatforms []PlatformInfo = []PlatformInfo{
	{ID: "it", Label:"Influweb.it", Country: "it"},
	{ID: "fr", Label:"Grippenet.fr", Country: "it"},
	{ID: "uk", Label:"Flusurvey", Country: "gb"},
	{ID: "nl", Label:"Infectieradar.nl", Country: "nl"},
	{ID: "be", Label:"Infectieradar.be", Country: "be"},
	{ID: "ch", Label:"Grippenet.ch", Country: "be"},
	{ID: "dki", Label:"Influmeter", Country: "dk"},
	{ID: "es", Label:"Gripenet.es", Country: "es"},
	{ID: "pt", Label:"Gripenet.pt", Country: "pt"},
}