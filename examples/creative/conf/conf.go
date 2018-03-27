package conf

// type symbol
const (
	UserPrefix       string = "USER_"
	ArtistPrefix     string = "ARTIST_"
	ProductionPrefix string = "PRODUCTION_"
	StateSplitSymbol string = ":"
	StateStartSymbol string = "0*"
	StateEndSymbol   string = "z*"
)

// invoke func name
const (
	AddUser    string = "AddUser"
	DeleteUser string = "DeleteUser"
	ModifyUser string = "ModifyUser"
	QueryUser  string = "QueryUser"
	ListOfUser string = "ListOfUser"

	AddArtist    string = "AddArtist"
	DeleteArtist string = "DeleteArtist"
	ModifyArtist string = "ModifyArtist"
	QueryArtist  string = "QueryArtist"
	ListOfArtist string = "ListOfArtist"

	AddProduction     string = "AddProduction"
	DeleteProduction  string = "DeleteProduction"
	ModifyProduction  string = "ModifyProduction"
	QueryProduction   string = "QueryProduction"
	ListOfProduction  string = "ListOfProduction"
	ListOfProduction2 string = "ListOfProduction2"

	ListOfSupporter string = "ListOfSupporter"
	AddSupporter    string = "AddSupporter"

	AddBuyer    string = "AddBuyer"
	ListOfBuyer string = "ListOfBuyer"

	// TODO
	DeleteBuyer string = "DeleteBuyer"
	ModifyBuyer string = "ModifyBuyer"
)
