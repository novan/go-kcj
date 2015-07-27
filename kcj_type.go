package kcj

type ScheduleItem struct {
	trainNumber     string
	misc            string
	class           string
	relation        string
	startingStation string
	currentStation  string
	endStation      string
	arrivingTime    string
	departingTime   string
	ls              string
	status          string
}

type Schedule struct {
	items      []ScheduleItem
	totalItems int
}

func (s *Schedule) Len() int {
	return len(s.items)
}

func (s *Schedule) Less(i, j int) bool {
	a1 := jktTime(s.items[i].arrivingTime)
	a2 := jktTime(s.items[j].arrivingTime)

	return a1.Before(a2)
}

func (s *Schedule) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *Schedule) IsAllSchedule() bool {
	return len(s.items) == s.totalItems
}

type ByRelation struct{ *Schedule }

func (b ByRelation) Less(i, j int) bool { return b.items[i].relation < b.items[j].relation }

type ByTrainNumber struct{ *Schedule }

func (b ByTrainNumber) Less(i, j int) bool { return b.items[i].trainNumber < b.items[j].trainNumber }

type ScheduleParam struct {
	station     string
	trainNumber string
}
