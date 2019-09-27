package persistence

//google account fields
type Account struct {
	Id  			string     	`json:"id" datastore:"id"`
	Name  			string     	`json:"name" datastore:"name"`
	UpdateTime   	string    	`json:"updateTime,omitempty" datastore:"updateTime,omitempty"`
	CreateTime      string    	`json:"createTime,omitempty" datastore:"createTime,omitempty"`
	Provider     	string		`json:"provider,omitempty" datastore:"provider,omitempty"`
	State 	 		string      `json:"state,omitempty" datastore:"state,omitempty"`
	Approvals    	[]Approval  `json:"approvals,omitempty" datastore:"approvals,omitempty"`
}

type Approval struct {
	Name  			string     	`json:"name" datastore:"name"`
	State  			string     	`json:"state" datastore:"state"`
	Reason  		string     	`json:"reason" datastore:"reason"`
	UpdateTime  	string     	`json:"updateTime" datastore:"updateTime"`
}

//cloudbees signup fields
type Contact struct {
	AccountId 		string     	`json:"accountId" datastore:"accountId"`
	FirstName 		string     	`json:"firstName,omitempty" datastore:"firstName,omitempty"`
	LastName		string     	`json:"lastName,omitempty" datastore:"lastName,omitempty"`
	EmailAddress	string     	`json:"emailAddress" datastore:"emailAddress"`
	Phone			string     	`json:"phone,omitempty" datastore:"phone,omitempty"`
	Company			string     	`json:"company,omitempty" datastore:"company,omitempty"`
	Timezone		string     	`json:"timezone,omitempty" datastore:"timezone,omitempty"`
}

//google entitlement fields
type Entitlement struct {
	Id     				string	`json:"id" datastore:"id"`
	Name     			string	`json:"name" datastore:"name"`
	Account   			string	`json:"account" datastore:"account"`
	Provider    		string	`json:"provider" datastore:"provider"`
	Product  			string	`json:"product" datastore:"product"`
	Plan     	  		string	`json:"plan" datastore:"plan"`
	NewPendingPlan 	  	string	`json:"newPendingPlan" datastore:"newPendingPlan"`
	State    	  		string	`json:"state" datastore:"state"`
	UpdateTime    	  	string	`json:"updateTime" datastore:"updateTime"`
	CreateTime    	  	string	`json:"createTime" datastore:"createTime"`
	UsageReportingId    string	`json:"usageReportingId" datastore:"usageReportingId"`
	MessageToUser    	string	`json:"messageToUser" datastore:"messageToUser"`
}
