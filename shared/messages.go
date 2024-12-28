package shared

type ItemMsg struct {
	Title   string
	Command string
}

type SaveCommandMsg struct {
	Command string
}

type AddNewItemMsg struct {
	Title       string
	Description string
	Command     string
}

type CopyToClipboardMsg struct {
	Command string
}

type ErrorMsg struct {
	Err error
}
