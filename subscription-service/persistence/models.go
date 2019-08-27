package persistence

type Account struct {
	Name  			string     	`json:"name"`
	UpdateTime   	string    	`json:"updateTime"`
	CreateTime      string    	`json:"createTime"`
	Provider     	string		`json:"provider"`
	State 	 		string      `json:"state"`
	Approvals    	string     	`json:"approvals"`
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
