package livereload

type (
  message struct {
	Command string `json:"command"`
}

updateMessage struct {
		message
		Url string
	  }
	reloadMessage struct {
		message
		Path    string `json:"path"`
		LiveCSS bool `json:"liveCSS"`
	  }

	alertMessage struct {
		message
		Message string `json:"message"`
	}

	helloMessage struct {
		message
		Protocols []string `json:"protocols"`
	}

	serverHello struct {
		message
		Protocols  []string `json:"protocols"`
		ServerName string   `json:"serverName"`
	}
)
