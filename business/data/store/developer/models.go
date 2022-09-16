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

type NewTeam struct {
	Name    string `validate:"required,max=255"`
	AppId   string `validate:"required"`
	Members []string
}

type TeamMember struct {
	Name   string `db:"member_name"`
	TeamId string `db:"team_id"`
}

type Team struct {
	Id      string `db:"id"`
	Name    string `db:"name"`
	AppId   string `db:"app_id"`
	Members []TeamMember
}
