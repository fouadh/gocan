package developer

type EntityDeveloper struct {
	Entity     string
	Author     string
	Added      int
	TotalAdded int
	Ownership  float64
}

type Developer struct {
	Name            string
	NumberOfCommits int
}

type EntityEffort struct {
	Entity          string
	Author          string
	AuthorRevisions int
	TotalRevisions  int
}

type KnowledgeMapHierarchy struct {
	Name          string                   `json:"name"`
	Children      []*KnowledgeMapHierarchy `json:"children,omitempty"`
	Weight        float64                  `json:"weight,omitempty"`
	Size          int                      `json:"size,omitempty"`
	MainDeveloper string                   `json:"mainDeveloper,omitempty"`
}
