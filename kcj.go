package kcj

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"errors"
)

const Version = "0.1.0"
const kcjBaseUrl = "http://www.krl.co.id/infonew/rute_jadwal.php" // please open the browser yourself
const maxChan = 10

func kcjHttpRequest(param *map[string]string) (content []byte, err error) {

	finalUrl, err := buildUrl(kcjBaseUrl, param)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", finalUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", randomUserAgents())
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return
}

func SchedulePage(query ScheduleParam, page int) (schedule *Schedule, err error) {
	param := make(map[string]string)
	param["stasiun_singgah"] = query.Station
	param["start"] = fmt.Sprintf("%d", page*10)
	param["no"] = fmt.Sprintf("%d", page*10+1)
	param["p_f"] = "0"
	param["ska_id"] = query.TrainNumber
	param["select"] = query.Relation
	param["select1"] = "10"						// Produk JABOTABEK = COMMUTER LINE
	param["jam1"] = query.HourFrom
	param["jam2"] = query.HourTo

	content, err := kcjHttpRequest(&param)

	doc, err := gokogiri.ParseHtml(content)

	if err != nil {
		return nil, err
	}

	const schXPath = "/html/body/table/tr[2]/td[2]/table/tr/td/table/tr"

	root := doc.Root()

	if root == nil {
		return nil, errors.New("Empty document?")
	}

	html := doc.Root().FirstChild()

	results, err := html.Search(schXPath)

	if err != nil {
		return nil, err
	} else if len(results) <= 0 {
		return nil, errors.New("Parsing failed")
	}

	defer doc.Free()

	schedule = &Schedule{}
	schedule.Items = make([]ScheduleItem, len(results)-1)

	for i, result := range results[1:] {
		schedule.Items[i], _ = trNodeToSchedule(result)
	}

	// getting total data
	const totalSchXPath = "/html/body/table/tr[3]/td[2]/table/tr[2]/td/text()"

	totalResults, err := html.Search(totalSchXPath)

	if err != nil {
		return nil, err
	}

	schedule.TotalItems, err =
		strconv.Atoi(strings.Fields(strings.TrimSpace(totalResults[0].String()))[6])

	if err != nil {
		return nil, err
	}

	return
}

func ScheduleAll(query ScheduleParam) (schedule *Schedule, err error) {
	schedule, err = SchedulePage(query, 0)

	if err != nil {
		return nil, err
	} else if schedule.TotalItems <= 0 {
		return nil, errors.New("No data")
	}

	pageCount := schedule.TotalItems / 10

	if schedule.TotalItems%10 > 0 {
		pageCount++
	}

	var wg sync.WaitGroup
	var chanNum = pageCount

	if chanNum > maxChan {
		chanNum = maxChan
	}

	c := make(chan *Schedule, chanNum)

	rf := func(page int) {
		// delay sometime not to overwhelm the server
		time.Sleep(300 * time.Millisecond)
		result, _ := SchedulePage(query, page)
		c <- result
		defer wg.Done()
	}

	for i := 1; i < pageCount; i++ {
		wg.Add(1)
		go rf(i)
	}

	for i := 1; i < pageCount; i++ {
		s := <-c

		if s == nil {
			continue
		}

		schedule.Items = append(schedule.Items, s.Items...)
	}

	return
}

func AllTrainNumbers() (numbers []string, err error) {

	content, err := kcjHttpRequest(nil)

	if err != nil {
		return nil, err
	}

	doc, err := gokogiri.ParseHtml(content)

	if err != nil {
		return nil, err
	}

	defer doc.Free()

	const trainNumXPath = "/html/body/form/table/tr[2]/td[2]/table/tr[4]/td[2]/select/option/text()"

	html := doc.Root().FirstChild()

	numResults, err := html.Search(trainNumXPath)

	if err != nil {
		return nil, err
	}

	numbers = make([]string, len(numResults)-1)

	for i, num := range numResults[1:] {
		numbers[i] = num.String()
	}

	return
}
