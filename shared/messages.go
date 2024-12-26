package shared

type ItemMsg struct {
	Title string
	Desc  string
}

type SaveItem struct {
	Desc string
}

type CopyToClipboard struct {
	Command string
}

type ErrorMsg struct {
	Err error
}
