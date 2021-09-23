package revision

type Revision struct {
	Entity                      string  `json:"entity"`
	NumberOfRevisions           int     `json:"numberOfRevisions"`
	NumberOfAuthors             int     `json:"numberOfAuthors"`
	NormalizedNumberOfRevisions float64 `json:"normalizedNumberOfRevisions"`
	Code                        int     `json:"code"`
}

type HotspotHierarchy struct {
	Name     string              `json:"name"`
	Children []*HotspotHierarchy `json:"children,omitempty"`
	Weight   float64             `json:"weight,omitempty"`
	Size     int                 `json:"size,omitempty"`
}

type RevisionTrend struct {
	Date      string     `json:"date"`
	Revisions []Revision `json:"revisions"`
}

func (rt RevisionTrend) FindEntityRevision(entityName string) Revision {
	for _, r := range rt.Revisions {
		if r.Entity == entityName {
			return r
		}
	}

	return Revision{}
}