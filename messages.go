package livereload

type message struct {
	Command string `json:"command"`
}

var (
	updateMessage = struct {
		message
		Url string
	}{
		message{"update"},
		"",
	}

	reloadMessage = struct {
		message
		Path    string `json:"path"`
		LiveCSS bool `json:"liveCSS"`
	}{
		message{"reload"},
		"",
		true,
	}

	alertMessage = struct {
		message
		Message string `json:"message"`
	}{
		message{"alert"},
		"",
	}

	helloMessage = struct {
		message
		Protocols []string `json:"protocols"`
	}{
		message{"hello"},
		protos,
	}

	serverHello = struct {
		message
		Protocols  []string `json:"protocols"`
		ServerName string   `json:"serverName"`
	}{
		helloMessage.message,
		helloMessage.Protocols,
		host,
	}
)
