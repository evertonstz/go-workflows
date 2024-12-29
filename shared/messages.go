package shared

type (
	ItemMsg struct {
		Title   string
		Command string
	}

	SaveCommandMsg struct {
		Command string
	}

	AddNewItemMsg struct {
		Title       string
		Description string
		Command     string
	}

	CopyToClipboardMsg struct {
		Command string
	}

	ErrorMsg struct {
		Err error
	}
)
