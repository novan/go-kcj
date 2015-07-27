package kcj

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"testing"
)

type ScheduleQueryFunc func(ScheduleParam) (*Schedule, error)

func queryAndPrint(scheduleFunction ScheduleQueryFunc, param ScheduleParam, t *testing.T) {
	schedule, err := scheduleFunction(param)

	if err != nil {
		t.Error(err)
		return
	}

	sort.Sort(schedule) // sorted by time

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Train Number", "Misc.", "Class", "Relation",
		"Starting", "Current", "End", "Arriving", "Departing", "LS", "Status"})

	for _, sched := range schedule.items {
		table.Append(
			[]string{
				sched.trainNumber,
				sched.misc,
				sched.class,
				sched.relation,
				sched.startingStation,
				sched.currentStation,
				sched.endStation,
				sched.arrivingTime,
				sched.departingTime,
				sched.ls,
				sched.status,
			})
	}

	table.Render()
}

func TestAllTrain(t *testing.T) {
	queryAndPrint(ScheduleAll, ScheduleParam{trainNumber: "1272"}, t)
}

func TestAllStation(t *testing.T) {
	queryAndPrint(ScheduleAll, ScheduleParam{station: "JNG"}, t)
}

func TestAllTrainNumbers(t *testing.T) {
	trainNumbers, err := AllTrainNumbers()
	sort.Strings(trainNumbers)

	if err != nil {
		t.Error(err)
		return
	}

	if trainNumbers == nil {
		t.Error("Cannot Get Train Numbers")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Train Number"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, num := range trainNumbers {
		table.Append([]string{num})
	}

	table.Render()
}

func TestStationPage(t *testing.T) {
	const currentPage = 0
	schedule, err := SchedulePage(ScheduleParam{station: "MRI"}, currentPage)

	if err != nil {
		t.Error(err)
		return
	}

	if schedule == nil {
		t.Error("Cannot Get Schedule")
		return
	}

	sort.Sort(ByTrainNumber{schedule})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Train Number", "Misc.", "Class", "Relation",
		"Starting", "Current", "Arriving", "Departing", "LS", "Status"})

	for _, sched := range schedule.items {
		table.Append(
			[]string{
				sched.trainNumber,
				sched.misc,
				sched.class,
				sched.relation,
				sched.startingStation,
				sched.currentStation,
				sched.arrivingTime,
				sched.departingTime,
				sched.ls,
				sched.status,
			})
	}

	table.Render()

	fmt.Printf("Total Records is: %v, currently Show Page %v, %v items\n",
		schedule.totalItems, currentPage, len(schedule.items))
}
