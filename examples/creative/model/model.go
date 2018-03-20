package model

// user info
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

// artist info
type Artist struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Username string `json:"username"`
}

// production info
type Production struct {
	Type                  string            `json:"type"`
	Serial                string            `json:"serial"`
	Name                  string            `json:"name"`
	Desc                  string            `json:"desc"`
	CopyrightPriceType    string            `json:"copyright_price_type"` // TODO list
	CopyrightPrice        string            `json:"copyright_price"`      // price（1-100）
	CopyrightNum          string            `json:"copyright_num"`        // total
	Username              string            `json:"username"`
	Supporters            map[string]string `json:"supporters"`
	Buyers                map[string]string `json:"buyers"`
	CopyrightTransferPart string            `json:"copyright_transfer_part"` //1 ~ num
}
