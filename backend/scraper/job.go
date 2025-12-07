package scraper

type Job struct {
	Profession       string `json:"profession"`
	Salary           string `json:"salary"`
	Company          string `json:"company"`
	Location         string `json:"location"`
	StartDate        string `json:"startDate"`
	Telephone        string `json:"telephone"`
	Email            string `json:"email"`
	ShortDescription string `json:"shortDescription"`
	FullDescription  string `json:"fullDescription"`
	RefNr            string `json:"refNr"`
	ExternalLink     string `json:"externalLink"`
	ApplicationLink  string `json:"applicationLink"`
}
