package ogame

// Chapter is a "directives" chapter
type Chapter struct {
	ID       int64
	ClaimAll *int64
	Tasks    []ChapterTask
}

// ChapterTask ...
type ChapterTask struct {
	ID        int64
	Completed bool
	Collected bool
}
