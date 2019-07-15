package timetracker

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Given an Unix epoch timestamp a binning function returns a string
// representing the bin.
type BinningFunc func(timestamp float64) string

type Context struct {
	Filestats map[string]map[string]float64
	Projects  map[string]map[string]float64

	lastEvent *DataPoint

	binning_func BinningFunc
	Debug        bool `json:"-"`
}

func (self *Context) AddDataPoint(point *DataPoint) {
	if point.Timestamp == 0 {
		return
	}

	if self.lastEvent == nil {
		self.lastEvent = point
	}

	// Another timestamp of the same file - just add time to the
	// file and project.
	added_time := point.Timestamp - self.lastEvent.Timestamp
	if added_time < 600 {
		if self.Debug {
			fmt.Printf(
				"Adding %v seconds to %s @ %v\n",
				added_time, self.lastEvent.Filename, time.Unix(int64(point.Timestamp), 0))
		}
		self._addTimeToStats(self.lastEvent, added_time)
	} else {
		if self.Debug {
			fmt.Printf(
				"Skipping update of %v seconds to %s\n",
				added_time, self.lastEvent.Filename)
		}
	}

	self.lastEvent = point
}

func (self *Context) _addTimeToStats(point *DataPoint, added_time float64) {
	timestamp := self.binning_func(point.Timestamp)
	file_stats_map, pres := self.Filestats[timestamp]
	if !pres {
		file_stats_map = make(map[string]float64)
		self.Filestats[timestamp] = file_stats_map
	}

	file_stats, _ := file_stats_map[point.Filename]
	file_stats_map[point.Filename] = file_stats + added_time

	project_stats_map, pres := self.Projects[timestamp]
	if !pres {
		project_stats_map = make(map[string]float64)
		self.Projects[timestamp] = project_stats_map
	}

	project_stats, _ := project_stats_map[point.Project]
	project_stats_map[point.Project] = project_stats + added_time

}

func (self *Context) Dump() string {
	dates := []string{}
	for k, _ := range self.Projects {
		dates = append(dates, k)
	}

	sort.Strings(dates)

	serialized, _ := json.MarshalIndent(self, " ", " ")
	return string(serialized)
}

func NewContext(binning_func BinningFunc) *Context {

	if binning_func == nil {
		binning_func = func(timestamp float64) string {
			return time.Unix(int64(timestamp), 0).Format("2006-01-02")
		}
	}

	return &Context{
		Filestats:    make(map[string]map[string]float64),
		Projects:     make(map[string]map[string]float64),
		binning_func: binning_func,
	}
}
