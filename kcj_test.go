package kcj

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
)

func TestAllStation(t *testing.T) {
	schedule, _ := ScheduleStation("CTA")
	sort.Sort(schedule) // sorted by time

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Train Number", "Misc.", "Class", "Relation",
		"Starting", "Current", "End", "Arriving", "Departing", "LS", "Status"})

	for _, sched := range schedule {
		table.Append(
			[]string{
				sched.trainNumber,
				sched.misc,
				sched.class,
				sched.relation,
				sched.startingStation,
				sched.currentStation,
				sched.endStation,
				sched.arrivingTime.Format(time.RFC822Z),
				sched.departingTime.Format(time.RFC822Z),
				sched.ls,
				sched.status,
			})
	}

	table.Render()

}

func TestStationPage(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	result, count, _ := ScheduleStationPage("MRI", 0)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Train Number", "Misc.", "Class", "Relation",
		"Starting", "Current", "Arriving", "Departing", "LS", "Status"})

	for _, sched := range result {
		table.Append(
			[]string{
				sched.trainNumber,
				sched.misc,
				sched.class,
				sched.relation,
				sched.startingStation,
				sched.currentStation,
				sched.arrivingTime.Format(time.RFC822Z),
				sched.departingTime.Format(time.RFC822Z),
				sched.ls,
				sched.status,
			})
	}

	table.Render()

	fmt.Printf("Total Records is: %v\n", count)

	// fmt.Printf("%+v\n", result)
}
