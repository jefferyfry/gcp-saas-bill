package persistence

type Account struct {
	//google fields
	Name  			string     	`json:"name"`
	UpdateTime   	string    	`json:"updateTime"`
	CreateTime      string    	`json:"createTime"`
	Provider     	string		`json:"provider"`
	State 	 		string      `json:"state"`
	Approvals    	string     	`json:"approvals"`

	//signup fields
	FirstName 		string     	`json:"firstName"`
	LastName		string     	`json:"lastName"`
	EmailAddress	string     	`json:"emailAddress"`
	Phone			string     	`json:"phone"`
	Company			string     	`json:"company"`
	Country			string     	`json:"country"`
	Timezone		string     	`json:"timezone"`
}

type Entitlement struct {
	Name     			string	`json:"name"`
	Account   			string	`json:"account"`
	Provider    		string	`json:"provider"`
	Product  			string	`json:"product"`
	Plan     	  		string	`json:"plan"`
	NewPendingPlan 	  	string	`json:"newPendingPlan"`
	State    	  		int64	`json:"state"`
	UpdateTime    	  	string	`json:"updateTime"`
	CreateTime    	  	string	`json:"createTime"`
	UsageReportingId    string	`json:"usageReportingId"`
	MessageToUser    	string	`json:"messageToUser"`
}
