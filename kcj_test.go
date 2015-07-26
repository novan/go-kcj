package kcj

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestStationPage(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	result, count, _ := ScheduleStationPage("CUK", 0)

	data := make([][10]string, len(result))

	for i, sched := range result {
		data[i] = [...]string{
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
		}
	}

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
				sched.arrivingTime,
				sched.departingTime,
				sched.ls,
				sched.status,
			})
	}

	table.Render()

	fmt.Printf("Total Records is: %v\n", count)

	// fmt.Printf("%+v\n", result)
}
