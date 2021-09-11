package cloc

type Clocs struct {
	Files []FileInfo
}

type FileInfo struct {
	Location string
	Code     int
}

type Info struct {
	File  string
	Lines int
}
