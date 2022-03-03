package developer

type EntityDeveloper struct {
	Entity     string
	Author     string
	Added      int
	TotalAdded int
	Ownership  float64
}

type Developer struct {
	Name            string `json:"name"`
	NumberOfCommits int    `json:"number_of_commits"`
}

type EntityEffortPerAuthor struct {
	Entity          string
	Author          string
	AuthorRevisions int
	TotalRevisions  int
}

type EntityEffort struct {
	Entity string  `db:"entity"`
	Effort float64 `db:"effort"`
}

type KnowledgeMapHierarchy struct {
	Name          string                   `json:"name"`
	Children      []*KnowledgeMapHierarchy `json:"children,omitempty"`
	Weight        float64                  `json:"weight,omitempty"`
	Size          int                      `json:"size,omitempty"`
	MainDeveloper string                   `json:"mainDeveloper,omitempty"`
	Effort        float64                  `json:"effort,omitempty"`
	DevDiffusion  float64                  `json:"devDiffusion,omitempty"`
}
