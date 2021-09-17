package revision

type Revision struct {
	Entity                      string `json:"entity"`
	NumberOfRevisions           int `json:"numberOfRevisions"`
	NumberOfAuthors             int `json:"numberOfAuthors"`
	NormalizedNumberOfRevisions float64 `json:"normalizedNumberOfRevisions"`
	Code                        int `json:"code"`
}

type HotspotHierarchy struct {
	Name          string              `json:"name"`
	Children      []*HotspotHierarchy `json:"children,omitempty"`
	Weight        float64             `json:"weight,omitempty"`
	Size          int                 `json:"size,omitempty"`
}
