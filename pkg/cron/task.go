package cron

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	rangeStepMatcher = regexp.MustCompile(`^(\*|\d{1,2}-\d{1,2})/(\d{1,2})$`)
	listMatcher      = regexp.MustCompile(`^(\d{1,2})(?:,(\d{1,2}))*$`)
	rangeMatcher     = regexp.MustCompile(`^(\d{1,2})-(\d{1,2})$`)
)

var (
	Secondly = "@secondly"
	Minutely = "@minutely"
	Hourly   = "@hourly"
)

var (
	macroMap = map[string]string{
		Secondly: "*/1 * *",
		Minutely: "0 */1 *",
		Hourly:   "0 0 */1",
	}
)

type Func func()

type Task struct {
	second []int8
	minute []int8
	hour   []int8
	dom    []int8
	month  []int8
	dow    []int8

	task Func
}

func NewTask(expr string, f Func) (*Task, error) {
	task, err := parseCronExpressionString(expr)

	if err != nil {
		return nil, err
	}

	task.task = f

	return task, nil
}

func (t *Task) Match(now time.Time) bool {
	now = now.Truncate(time.Second)

	// fmt.Printf("now: %v\n", now)
	// fmt.Printf("task: %+v\n", t)

	if !slices.Contains(t.minute, int8(now.Minute())) ||
		!slices.Contains(t.hour, int8(now.Hour())) ||
		!slices.Contains(t.month, int8(now.Month())) ||
		!slices.Contains(t.second, int8(now.Second())) {
		return false
	}

	domMatch := slices.Contains(t.dom, int8(now.Day()))
	dowMatch := slices.Contains(t.dow, int8(now.Weekday()))

	return domMatch && dowMatch
}

func parseCronExpressionString(expr string) (*Task, error) {
	if strings.HasPrefix(expr, "@") {
		e, ok := macroMap[expr]

		if ok != true {
			return nil, fmt.Errorf("unknown macro %s", expr)
		}

		return parseCronExpressionString(e)
	}

	expressions := strings.Split(expr, " ")

	if len(expressions) < 3 {
		return nil, fmt.Errorf("invalid cron format: expected 3-6 fields")
	}

	seconds, err := expandExpression(expressions[0], generateOptions(0, 60, 1))

	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", "seconds", err)
	}
	minutes, err := expandExpression(expressions[1], generateOptions(0, 60, 1))

	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", "minutes", err)
	}

	hours, err := expandExpression(expressions[2], generateOptions(0, 24, 1))

	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", "hours", err)
	}

	task := Task{
		second: seconds,
		minute: minutes,
		hour:   hours,
		dom:    generateOptions(1, 32, 1),
		month:  generateOptions(1, 13, 1),
		dow:    generateOptions(0, 7, 1),
	}

	if len(expressions) > 3 {
		dom, err := expandExpression(expressions[3], generateOptions(1, 32, 1))

		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", "dom", err)
		}

		task.dom = dom
	}

	if len(expressions) > 4 {
		month, err := expandExpression(expressions[4], generateOptions(1, 13, 1))

		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", "month", err)
		}
		task.month = month
	}
	if len(expressions) > 5 {
		dow, err := expandExpression(expressions[5], generateOptions(0, 8, 1))

		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", "dow", err)
		}

		if slices.Contains(dow, 7) && !slices.Contains(dow, 0) {
			dow = append(dow, 0)
		}
	} else {
		task.dow = append(task.dow, 0)
	}

	return &task, nil
}

func expandExpression(expr string, options []int8) ([]int8, error) {
	if expr == "*" {
		return options, nil
	}

	rangeMatches := rangeMatcher.FindStringSubmatch(expr)

	if rangeMatches != nil {
		// NOTE: ignoring errors because using \d matcher in regex

		left, _ := strconv.ParseInt(rangeMatches[1], 10, 8)
		right, _ := strconv.ParseInt(rangeMatches[2], 10, 8)

		return generateOptions(int8(left), int8(right), 1), nil
	}

	listMatches := listMatcher.FindStringSubmatch(expr)

	if listMatches != nil {
		rawValues := strings.Split(listMatches[0], ",")

		values := make([]int8, 0, len(rawValues))

		for _, v := range rawValues {
			d, _ := strconv.ParseInt(v, 10, 8)
			values = append(values, int8(d))
		}

		return values, nil
	}

	rangeStepMatches := rangeStepMatcher.FindStringSubmatch(expr)

	if rangeStepMatches != nil {
		new_options, err := expandExpression(rangeStepMatches[1], options)

		if err != nil {
			return nil, err
		}

		interval, _ := strconv.ParseInt(rangeStepMatches[2], 10, 8)

		return generateOptions(
			new_options[0],
			new_options[len(new_options)-1],
			int8(interval),
		), nil
	}

	return nil, fmt.Errorf("Failed to parse cron subexpression: unknown format")
}

func generateOptions(left, right, step int8) []int8 {
	if right <= left {
		return []int8{}
	}

	if step < 1 {
		return []int8{}
	}

	result := make([]int8, 0, (right-left)/step)
	cur := left
	for cur < right {
		result = append(result, cur)
		cur += step
	}

	return result
}
