package timetracker

import (
	"bufio"
	"encoding/json"
	"io"
)

func ProcessFile(fd io.Reader, ctx *Context) error {
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		point := &DataPoint{}
		err := json.Unmarshal([]byte(scanner.Text()), &point)
		if err != nil {
			return err
		}

		ctx.AddDataPoint(point)
	}

	return nil
}
