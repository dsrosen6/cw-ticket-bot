package webex

type Webhook struct {
	Name      string `json:"name"`
	TargetUrl string `json:"targetUrl"`
	Resource  string `json:"resource"`
	Event     string `json:"event"`
	Filter    string `json:"filter"`
}
