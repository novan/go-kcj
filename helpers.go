package kcj

import (
	"github.com/moovweb/gokogiri/xml"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func mapToQuery(m *map[string]string) string {
	if m == nil || len(*m) == 0 {
		return ""
	} else {
		params := url.Values{}
		for k, v := range *m {
			params.Add(k, v)
		}
		return params.Encode()
	}
}

func buildUrl(base string, qs *map[string]string) (url *url.URL, err error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return nil, err
	} else {
		if qs == nil {
			return baseUrl, nil
		} else {
			baseUrl.RawQuery = mapToQuery(qs)
			return baseUrl, nil
		}
	}
}

var rndInitialised bool = false

func randomUserAgents() string {

	if !rndInitialised {
		rand.Seed(time.Now().UTC().UnixNano())
	}
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
	idx := rand.Intn(len(userAgents))
	return userAgents[idx]
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

	// Train service starts at 4, so midnite train are for next day
	theDay := nowjkt.Day()
	if hr < 4 {
		theDay = nowjkt.Day() + 1
	}

	return time.Date(nowjkt.Year(), nowjkt.Month(), theDay, hr, mn, sec, 0, loca)
}

func trNodeToSchedule(scheduleNode xml.Node) (item ScheduleItem, err error) {

	results, err := scheduleNode.Search("./td/text()")

	if err != nil {
		return ScheduleItem{}, err
	}

	item = ScheduleItem{
		trainNumber:     strings.TrimSpace(results[1].String()),
		misc:            strings.TrimSpace(results[2].String()),
		class:           strings.TrimSpace(results[3].String()),
		relation:        strings.TrimSpace(results[4].String()),
		startingStation: strings.TrimSpace(results[5].String()),
		currentStation:  strings.TrimSpace(results[6].String()),
		arrivingTime:    strings.TrimSpace(results[7].String()),
		departingTime:   strings.TrimSpace(results[8].String()),
		ls:              strings.TrimSpace(results[9].String()),
	}

	if len(results) > 10 {
		item.status = strings.TrimSpace(results[10].String())
	}

	stationParts := strings.FieldsFunc(item.relation, func(r rune) bool {
		return r == '-'
	})

	item.endStation = stationParts[1] // [ANGKE BOGOR] BOGOR is end station

	return
}
