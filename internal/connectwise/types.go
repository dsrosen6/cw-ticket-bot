package connectwise

import "time"

type WebhookBody struct {
	MessageId         string      `json:"MessageId"`
	FromUrl           string      `json:"FromUrl"`
	CompanyId         string      `json:"CompanyId"`
	MemberId          string      `json:"MemberId"`
	Action            string      `json:"Action"`
	Type              string      `json:"Type"`
	ID                int         `json:"ID"`
	ProductInstanceId interface{} `json:"ProductInstanceId"`
	PartnerId         interface{} `json:"PartnerId"`
	Entity            string      `json:"Entity"`
	Metadata          struct {
		KeyUrl string `json:"key_url"`
	} `json:"Metadata"`
	CallbackObjectRecId int `json:"CallbackObjectRecId"`
}

type Member struct {
	Id           int    `json:"id"`
	PrimaryEmail string `json:"primaryEmail"`
	Identifier   string `json:"identifier"`
}

type Note struct {
	Id                    int    `json:"id"`
	TicketId              int    `json:"ticketId"`
	Text                  string `json:"text"`
	DetailDescriptionFlag bool   `json:"detailDescriptionFlag"`
	InternalAnalysisFlag  bool   `json:"internalAnalysisFlag"`
	ResolutionFlag        bool   `json:"resolutionFlag"`
	IssueFlag             bool   `json:"issueFlag"`
	Member                struct {
		Id         int    `json:"id"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
		Info       struct {
			MemberHref string `json:"member_href"`
			ImageHref  string `json:"image_href"`
		} `json:"_info"`
	} `json:"member"`
	DateCreated  time.Time `json:"dateCreated"`
	CreatedBy    string    `json:"createdBy"`
	InternalFlag bool      `json:"internalFlag"`
	ExternalFlag bool      `json:"externalFlag"`
	Info         struct {
		LastUpdated time.Time `json:"lastUpdated"`
		UpdatedBy   string    `json:"updatedBy"`
	} `json:"_info"`
}

type Ticket struct {
	Id                      int    `json:"id"`
	Summary                 string `json:"summary"`
	InitialDescription      string `json:"initialDescription,omitempty"`
	InitialInternalAnalysis string `json:"initialInternalAnalysis,omitempty"`
	Board                   struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"board"`
	Company struct {
		Name string `json:"name"`
	} `json:"company"`
	ContactName string `json:"contactName"`
	Owner       struct {
		Id         int    `json:"id,omitempty"`
		Identifier string `json:"identifier,omitempty"`
	} `json:"owner"`
	ClosedFlag bool   `json:"closedFlag"`
	Resources  string `json:"resources,omitempty"`
	Info       struct {
		LastUpdated         time.Time `json:"lastUpdated"`
		UpdatedBy           string    `json:"updatedBy"`
		DateEntered         time.Time `json:"dateEntered"`
		EnteredBy           string    `json:"enteredBy"`
		ActivitiesHref      string    `json:"activities_href"`
		ScheduleentriesHref string    `json:"scheduleentries_href"`
		DocumentsHref       string    `json:"documents_href"`
		ConfigurationsHref  string    `json:"configurations_href"`
		TasksHref           string    `json:"tasks_href"`
		NotesHref           string    `json:"notes_href"`
		ProductsHref        string    `json:"products_href"`
		TimeentriesHref     string    `json:"timeentries_href"`
		ExpenseEntriesHref  string    `json:"expenseEntries_href"`
	} `json:"_info"`
}
