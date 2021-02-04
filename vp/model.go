package vp

type Data struct {
	ChartProps struct {
		XDelta int `json:"XDelta"`
		Nodes  []struct {
			Name string  `json:"name"`
			X    float64 `json:"x"`
			Y    int     `json:"y"`
		} `json:"nodes"`
		TotalDays int     `json:"totalDays"`
		H         int     `json:"h"`
		W         int     `json:"w"`
		TagAng    float64 `json:"TagAng"`
	} `json:"chartProps"`
	Marking []struct {
		Name   string `json:"Name"`
		Childs []struct {
			Name    string `json:"Name"`
			Percent string `json:"Percent"`
		} `json:"Childs"`
	} `json:"marking"`
	BasicInfo []struct {
		Title  string `json:"Title"`
		Childs []struct {
			K string `json:"k"`
			V string `json:"v"`
		} `json:"Childs"`
	} `json:"basicInfo"`
	BasicInfo2 []struct {
		Title  string `json:"Title"`
		Childs []struct {
			Title  string `json:"Title"`
			Childs []struct {
				K string `json:"k"`
				V string `json:"v"`
			} `json:"Childs"`
		} `json:"Childs"`
	} `json:"basicInfo2"`
	Progress []struct {
		Name   string `json:"Name"`
		Childs []struct {
			Name    string `json:"Name"`
			Percent string `json:"Percent"`
			Date    string `json:"Date"`
		} `json:"Childs"`
	} `json:"progress"`
}