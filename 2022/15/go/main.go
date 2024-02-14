package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Sensor struct {
	pos           Point
	closestBeacon Point
}

func (s Sensor) manhattan() int {
	return s.manhattanTo(s.closestBeacon)
}

func (s Sensor) manhattanTo(p Point) int {
	return abs(s.pos.x-p.x) + abs(s.pos.y-p.y)
}

type Range struct {
	xLeft, xRight int
}

func (s Sensor) getRangeX(y int) (Range, bool) {
	r := s.manhattan()
	dr := r - abs(s.pos.y-y)
	if dr < 0 {
		return Range{}, false
	}
	xl := s.pos.x - dr
	xr := s.pos.x + dr

	return Range{xl, xr}, true
}

func parseSensor(input string) Sensor {
	inputs := strings.Split(input, " ")

	parseNumber := func(input string, prefix string, suffix string) int {
		c, _ := strings.CutPrefix(input, prefix)
		c, _ = strings.CutSuffix(c, suffix)
		res, _ := strconv.Atoi(c)
		return res
	}

	return Sensor{
		pos: Point{
			x: parseNumber(inputs[2], "x=", ","),
			y: parseNumber(inputs[3], "y=", ":"),
		},
		closestBeacon: Point{
			x: parseNumber(inputs[8], "x=", ","),
			y: parseNumber(inputs[9], "y=", ""),
		},
	}
}

func checkInArea(ss []Sensor, p Point) bool {
	for _, s := range ss {
		r := s.manhattan()
		d := s.manhattanTo(p)

		if d <= r {
			return true
		}
	}
	return false
}

func checkBeaconIntersect(ss []Sensor, p Point) bool {
	for _, s := range ss {
		if p == s.closestBeacon {
			return true
		}
	}
	return false
}

func compressSensorsX(sensors []Sensor, y int) []Range {
	ranges := []Range{}
	for _, s := range sensors {
		r, exists := s.getRangeX(y)
		if exists {
			ranges = append(ranges, r)
		}
	}

	compressedRanges := []Range{}

	if len(ranges) == 0 {
		return compressedRanges
	}

	sort.Sort(Ranges(ranges))
	// fmt.Println(ranges)

	compressedRanges = append(compressedRanges, ranges[0])
	ranges = ranges[1:]

	for len(ranges) > 0 {
		lastRange := compressedRanges[len(compressedRanges)-1]
		nextRange := ranges[0]

		compressedRanges = compressedRanges[:len(compressedRanges)-1]
		ranges = ranges[1:]

		if lastRange.xRight >= nextRange.xLeft {
			compressedRanges = append(compressedRanges, Range{lastRange.xLeft, max(lastRange.xRight, nextRange.xRight)})
		} else {
			compressedRanges = append(compressedRanges, lastRange)
			compressedRanges = append(compressedRanges, nextRange)
		}

		// fmt.Println("====================================")
		// fmt.Println(ranges)
		// fmt.Println(compressedRanges)
	}

	return compressedRanges
}

func countBeaconsX(sensors []Sensor, y int) int {
	res := 0
	xs := map[int]bool{}
	for _, sensor := range sensors {
		if sensor.closestBeacon.y == y {
			_, was := xs[sensor.closestBeacon.x]
			if !was {
				res += 1
				xs[sensor.closestBeacon.x] = true
			}
		}
	}
	return res
}

func main() {
	content, _ := os.ReadFile("input.txt")
	contentStr := strings.Trim(string(content), "\n\r\t ")

	sensors := []Sensor{}
	left, right := math.MaxInt, math.MinInt
	for _, line := range strings.Split(contentStr, "\n") {
		sensor := parseSensor(line)
		sensors = append(sensors, sensor)

		dist := abs(sensor.pos.x - sensor.closestBeacon.x)
		lBound := sensor.pos.x - dist
		if lBound < left {
			left = lBound
		}

		rBound := sensor.pos.x + dist
		if rBound > right {
			right = rBound
		}
	}

	res := 0

    const Y = 4000000

	// ranges := compressSensorsX(sensors, Y)
	// fmt.Println(ranges)
	// for _, rng := range ranges {
	// 	res += rng.xRight - rng.xLeft + 1
	// }
	//
	// res -= countBeaconsX(sensors, Y)

    for y := 0; y <= Y; y += 1 {
        ranges := compressSensorsX(sensors, y)
        // fmt.Println(ranges)

        if len(ranges) > 1 {
            fmt.Println(ranges)
            res = (ranges[0].xRight + 1) * Y + y
            break
        }
    }

	fmt.Println(res)
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func max(ns ...int) int {
	res := math.MinInt
	for _, n := range ns {
		if n > res {
			res = n
		}
	}
	return res
}

type Ranges []Range

func (rs Ranges) Len() int {
	return len(rs)
}

func (rs Ranges) Less(i, j int) bool {
	return rs[i].xLeft < rs[j].xLeft
}

func (rs Ranges) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}
