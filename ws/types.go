package ws

// Command is the JSON format between web server and docker server
type Command struct {
	Command     string   `json:"command"`
	PWD         string   `json:"pwd"`
	ENV         []string `json:"env"`
	UserName    string   `json:"user"`
	ProjectName string   `json:"project"`
	Entrypoint  []string `json:"entrypoint"`
	Type        string   `json:"type"`
}

// ClientDebugMessage stores the data received from user
type ClientDebugMessage struct {
	BreakPoints string
	Command     string
}
