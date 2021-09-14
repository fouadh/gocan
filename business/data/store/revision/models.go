package revision

type Revision struct {
	Entity                      string
	NumberOfRevisions           int
	NumberOfAuthors             int
	NormalizedNumberOfRevisions float64
	Code                        int
}

type HotspotHierarchy struct {
	Name          string              `json:"name"`
	Children      []*HotspotHierarchy `json:"children,omitempty"`
	Weight        float64             `json:"weight,omitempty"`
	Size          int                 `json:"size,omitempty"`
}
