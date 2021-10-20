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

type NewRevisionTrends struct {
	Id         string `db:"id"`
	Name       string `db:"name"`
	BoundaryId string `db:"boundary_id"`
	Entries    []NewRevisionTrend
}

type RevisionTrends struct {
	Id         string `db:"id"`
	Name       string `db:"name"`
	Entries    []RevisionTrend
}

type TrendRevision struct {
	EntryId           string `db:"entry_id"`
	Entity            string `db:"entity"`
	NumberOfRevisions int    `db:"number_of_revisions"`
}

type NewRevisionTrend struct {
	Id              string `db:"id"`
	Date            string `db:"date"`
	Revisions       []TrendRevision
	RevisionTrendId string `db:"revision_trend_id"`
}

type RevisionTrend struct {
	Date      string          `json:"date"`
	Revisions []TrendRevision `json:"revisions"`
}

func (rt RevisionTrend) FindEntityRevision(entityName string) TrendRevision {
	for _, r := range rt.Revisions {
		if r.Entity == entityName {
			return r
		}
	}

	return TrendRevision{}
}
