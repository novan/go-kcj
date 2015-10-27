package kcj

type ScheduleItem struct {
	TrainNumber     string
	Misc            string
	Class           string
	Relation        string
	StartingStation string
	CurrentStation  string
	EndStation      string
	ArrivingTime    string
	DepartingTime   string
	Ls              string
	Status          string
}

type Schedule struct {
	Items      []ScheduleItem
	TotalItems int
}

func (s *Schedule) Len() int {
	return len(s.Items)
}

func (s *Schedule) Less(i, j int) bool {
	a1 := jktTime(s.Items[i].ArrivingTime)
	a2 := jktTime(s.Items[j].ArrivingTime)

	return a1.Before(a2)
}

func (s *Schedule) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

func (s *Schedule) IsAllSchedule() bool {
	return len(s.Items) == s.TotalItems
}

type ByRelation struct{ *Schedule }

func (b ByRelation) Less(i, j int) bool { return b.Items[i].Relation < b.Items[j].Relation }

type ByTrainNumber struct{ *Schedule }

func (b ByTrainNumber) Less(i, j int) bool { return b.Items[i].TrainNumber < b.Items[j].TrainNumber }

type ScheduleParam struct {
	Station     string
	TrainNumber string
	Relation	string
	HourFrom	string
	HourTo		string
}
