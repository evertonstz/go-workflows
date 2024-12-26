package shared

type ItemMsg struct {
	Title   string
	Command string
}

type SaveItem struct {
	Command string
}

type CopyToClipboard struct {
	Command string
}

type ErrorMsg struct {
	Err error
}
