package persistence

type Account struct {
	//google fields
	Name  			string     	`json:"name,omitempty"`
	UpdateTime   	string    	`json:"updateTime,omitempty"`
	CreateTime      string    	`json:"createTime,omitempty"`
	Provider     	string		`json:"provider,omitempty"`
	State 	 		string      `json:"state,omitempty"`
	Approvals    	string     	`json:"approvals,omitempty"`

	//signup fields
	FirstName 		string     	`json:"firstName,omitempty"`
	LastName		string     	`json:"lastName,omitempty"`
	EmailAddress	string     	`json:"emailAddress,omitempty"`
	Phone			string     	`json:"phone,omitempty"`
	Company			string     	`json:"company,omitempty"`
	Timezone		string     	`json:"timezone,omitempty"`
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
