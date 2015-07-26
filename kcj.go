package kcj

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const Version = "0.1.0"

// This is redacted, because I don't want this to be able to be searched
// Please check the URL yourself

const kcjBaseUrl = "REDACTED_SO_IT_WONT_BE_ABLE_TO_BE_SEARCHED" // please open the browser yourself

type ScheduleItem struct {
	trainNumber     string
	misc            string
	class           string
	relation        string
	startingStation string
	currentStation  string
	arrivingTime    time.Time
	departingTime   time.Time
	ls              string //?
	status          string //?
}

func mapToQuery(m map[string]string) string {
	if len(m) == 0 {
		return ""
	} else {
		params := url.Values{}
		for k, v := range m {
			params.Add(k, v)
		}
		return params.Encode()
	}
}

func buildUrl(base string, qs map[string]string) (url *url.URL, err error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return nil, err
	} else {
		baseUrl.RawQuery = mapToQuery(qs)
		return baseUrl, nil
	}
}

func ScheduleStationPage(station string, page int) (schedule []ScheduleItem, totalCount int, err error) {
	// Randomise User Agents just for fun, we'll use Console's UA, and OLD OS
	var userAgents = [...]string{
		"Mozilla/5.0 (PlayStation 4 2.57) AppleWebKit/536.26 (KHTML, like Gecko)",
		"Opera/9.50 (Nintendo DSi; Opera/507; U; en-US)",
		"Mozilla/5.0 (Nintendo 3DS; U; ; en) Version/1.7498.US",
		"AmigaVoyager/3.2 (AmigaOS/MC680x0)",
		"NCSA Mosaic/3.0 (Windows 95)",
		"Mozilla/3.0 (Planetweb/2.100 JS SSL US; Dreamcast US)",
		"Mozilla/5.0 (PlayStation Vita 1.80) AppleWebKit/531.22.8 (KHTML, like Gecko) Silk/3.2",
		"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
	}

	param := make(map[string]string)
	param["stasiun_singgah"] = station
	param["start"] = fmt.Sprintf("%d", page*10)
	param["no"] = fmt.Sprintf("%d", page*10+1)
	param["p_f"] = "0"

	finalUrl, err := buildUrl(kcjBaseUrl, param)

	if err != nil {
		return nil, 0, err
	}

	idx := rand.Intn(len(userAgents))
	userAgent := userAgents[idx]

	client := &http.Client{}
	req, err := http.NewRequest("GET", finalUrl.String(), nil)

	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)

	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, 0, err
	}

	doc, err := gokogiri.ParseHtml(content)

	if err != nil {
		return nil, 0, err
	}

	schXPath := "/html/body/table/tr[2]/td[2]/table/tr/td/table/tr"
	html := doc.Root().FirstChild()

	results, err := html.Search(schXPath)

	if err != nil {
		return nil, 0, err
	}

	defer doc.Free()

	sched := make([]ScheduleItem, len(results)-1)

	for i := 1; i < len(results); i++ {
		sched[i-1], _ = trNodeToSchedule(results[i])
	}

	// getting total data
	totalSchXPath := "/html/body/table/tr[3]/td[2]/table/tr[2]/td/text()"

	totalResults, err := html.Search(totalSchXPath)

	tot, _ := strconv.Atoi(strings.Fields(strings.TrimSpace(totalResults[0].String()))[6])

	return sched, tot, nil
	// return schedule, nil
	// return string(results[0].String()), nil
}

// Convert timestamp to Jakarta Time
func jktTime(timestr string) time.Time {
	loca, _ := time.LoadLocation("Asia/Jakarta")
	nowjkt := time.Now().In(loca)
	parts := strings.FieldsFunc(timestr, func(r rune) bool {
		return r == ':'
	})
	hr, _ := strconv.Atoi(parts[0])
	mn, _ := strconv.Atoi(parts[1])
	sec, _ := strconv.Atoi(parts[2])

	return time.Date(nowjkt.Year(), nowjkt.Month(), nowjkt.Day(), hr, mn, sec, 0, loca)
}

func trNodeToSchedule(scheduleNode xml.Node) (item ScheduleItem, err error) {

	results, err := scheduleNode.Search("./td/text()")

	if err != nil {
		return ScheduleItem{}, err
	}

	if len(results) > 10 {
		item = ScheduleItem{
			trainNumber:     strings.TrimSpace(results[1].String()),
			misc:            strings.TrimSpace(results[2].String()),
			class:           strings.TrimSpace(results[3].String()),
			relation:        strings.TrimSpace(results[4].String()),
			startingStation: strings.TrimSpace(results[5].String()),
			currentStation:  strings.TrimSpace(results[6].String()),
			arrivingTime:    jktTime(strings.TrimSpace(results[7].String())),
			departingTime:   jktTime(strings.TrimSpace(results[8].String())),
			ls:              strings.TrimSpace(results[9].String()),
			status:          strings.TrimSpace(results[10].String()),
		}
	} else {
		item = ScheduleItem{
			trainNumber:     strings.TrimSpace(results[1].String()),
			misc:            strings.TrimSpace(results[2].String()),
			class:           strings.TrimSpace(results[3].String()),
			relation:        strings.TrimSpace(results[4].String()),
			startingStation: strings.TrimSpace(results[5].String()),
			currentStation:  strings.TrimSpace(results[6].String()),
			arrivingTime:    jktTime(strings.TrimSpace(results[7].String())),
			departingTime:   jktTime(strings.TrimSpace(results[8].String())),
			ls:              strings.TrimSpace(results[9].String()),
		}
	}
	return
}

func ScheduleStation(station string) (schedule []ScheduleItem, err error) {

	// get first page
	result, count, _ := ScheduleStationPage(station, 0)

	// Allocate all schedule
	schedule = make([]ScheduleItem, len(result))

	// Copy the first page
	copy(schedule, result)

	var pageCount int = count / 10

	if count%10 > 0 {
		pageCount++
	}

	// allocate WaitGroup
	var wg sync.WaitGroup

	c := make(chan []ScheduleItem)

	for i := 1; i < pageCount; i++ {
		wg.Add(1)
		go func(page int) {
			result, _, _ := ScheduleStationPage(station, page)
			c <- result
			defer wg.Done()
		}(i)
	}

	for i := 1; i < pageCount; i++ {
		items := <-c
		for _, item := range items {
			schedule = append(schedule, item)
		}
	}
	wg.Wait()
	err = nil
	return
}
