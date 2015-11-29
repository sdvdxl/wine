package time

import (
	"errors"
	"github.com/sdvdxl/go-tools/collections"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	FORMAT_ERROR = "format error"
)

type Crontab struct {
	ticker         *time.Ticker
	tickerDuration time.Duration //周期值, 秒1-59, 分 1-59， 日 1-31, 月 1-12，星期 0-6，0代表周日， 其中星期不可和日或者月同时使用
	isTicker       bool
	seconds        *collections.Set //0-59
	minutes        *collections.Set // 0-59
	hours          *collections.Set //0-23
	days           *collections.Set // 1-31
	monthes        *collections.Set // 1-12
	weekdays       *collections.Set //0-6
	handlers       []Handler
	stop           bool
}

type Handler interface{}

//create crontab jobs from args
func NewCrontab(args string, handlers ...Handler) (Crontab, error) {
	return buildCronab(args, handlers...)
}

// start the crontab jobs
func (c *Crontab) Start() {
	if c.isTicker {
		c.ticker = time.NewTicker(c.tickerDuration)
		for {
			<-c.ticker.C
			execHandlers(c.handlers...)
		}
	} else {
		for {
			now := time.Now()
			// s, m, d, M|W
			if c.minutes.Size() == 0 { //只有秒
				if c.seconds.Contains(now.Second()) {
					execHandlers(c.handlers...)
				}
			} else if c.hours.Size() == 0 { //有秒，分
				if c.seconds.Contains(now.Second()) && c.minutes.Contains(now.Minute()) {
					execHandlers(c.handlers...)
				}
			} else if c.days.Size() == 0 { //有秒，分, 时
				if c.seconds.Contains(now.Second()) && c.minutes.Contains(now.Minute()) && c.hours.Contains(now.Hour()) {
					execHandlers(c.handlers...)
				}
			} else if c.monthes.Size() == 0 && c.weekdays.Size() == 0 { //有秒，分, 时,天
				if c.seconds.Contains(now.Second()) && c.minutes.Contains(now.Minute()) && c.hours.Contains(now.Hour()) && c.days.Contains(now.Day()) {
					execHandlers(c.handlers...)
				}
			} else {
				if c.seconds.Contains(now.Second()) && c.minutes.Contains(now.Minute()) && c.hours.Contains(now.Hour()) && c.days.Contains(now.Day()) && (c.monthes.Contains(now.Month()) || c.weekdays.Contains(now.Weekday())) {
					execHandlers(c.handlers...)
				}
			}

			time.Sleep(time.Second)
		}
	}
}

// stop the crontab jobs
func (c *Crontab) Stop() {
	if c.isTicker {
		if c.ticker != nil {
			c.ticker.Stop()
		}
	} else {
		c.stop = true
	}

}

func buildCronab(args string, handlers ...Handler) (Crontab, error) {
	crontab := Crontab{}
	for _, handler := range handlers {
		err := validateHandler(handler)
		if err != nil {
			return crontab, err
		}

		crontab.handlers = append(crontab.handlers, handler)

	}

	if err := crontab.parse(args); err != nil {
		return Crontab{}, err
	}

	return crontab, nil
}

func (c *Crontab) parse(argline string) error {
	args := strings.Split(argline, " ")
	if len(args) != 6 {
		err := errors.New(FORMAT_ERROR)
		return err
	}

	if err := c.parseSecond(args[0]); err != nil {
		return err
	}

	if err := c.parseMinute(args[1]); err != nil {
		return err
	}

	if err := c.parseHour(args[2]); err != nil {
		return err
	}

	if err := c.parseDay(args[3]); err != nil {
		return err
	}

	if err := c.parseMonth(args[4]); err != nil {
		return err
	}

	if err := c.parseWeek(args[5]); err != nil {
		return err
	}

	if c.weekdays.Size() == 7 && c.monthes.Size() == 12 && c.days.Size() == 31 && c.hours.Size() == 24 && c.minutes.Size() == 60 && c.seconds.Size() == 60 {
		c.isTicker = true
		c.tickerDuration = time.Second
	}

	return nil

}

//check handler if is a function type
func validateHandler(handler Handler) error {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		return errors.New("crontab handler must be a callable func")
	}

	return nil
}

//解析星期， 只允许
// 0-6 星期天到星期六
// 0，1 星期日和星期一
// * 忽略
//允许用英文全称和简写
func (c *Crontab) parseWeek(arg string) error {
	set, isticker, err := parseArgument(arg, 0, 6, "week")
	if err != nil {
		return err
	}

	if isticker {
		c.isTicker = true
		v, _ := set.Values()[0].(int)
		c.tickerDuration = time.Hour * 24 * 7 * time.Duration(v)
	}

	c.weekdays = set
	return nil
}

//解析月份， 只允许
// 1-12 1月到12月，相当于 */1 每个月
// 1,2 1月和2月
// * 忽略
//允许用英文全称和简写
func (c *Crontab) parseMonth(arg string) error {
	set, isticker, err := parseArgument(arg, 1, 12, "month")
	if err != nil {
		return err
	}

	if isticker {
		c.isTicker = true
		v, _ := set.Values()[0].(int)
		c.tickerDuration = time.Hour * 24 * 30 * time.Duration(v)
	}

	c.monthes = set
	return nil
}

//解析日， 只允许
// 1-12 1月到12月，相当于 */1 每个月
// 1,2 1月和2月
// * 忽略
//允许用英文全称和简写
func (c *Crontab) parseDay(arg string) error {
	set, isticker, err := parseArgument(arg, 1, 31, "")
	if err != nil {
		return err
	}

	if isticker {
		c.isTicker = true
		v, _ := set.Values()[0].(int)
		c.tickerDuration = time.Hour * 24 * time.Duration(v)
	}

	c.days = set
	return nil
}

func (c *Crontab) parseHour(arg string) error {
	set, isticker, err := parseArgument(arg, 0, 23, "")
	if err != nil {
		return err
	}

	if isticker {
		c.isTicker = true
		v, _ := set.Values()[0].(int)
		c.tickerDuration = time.Hour * time.Duration(v)
	}

	c.hours = set
	return nil
}

func (c *Crontab) parseMinute(arg string) error {
	set, isticker, err := parseArgument(arg, 0, 59, "")
	if err != nil {
		return err
	}

	if isticker {
		c.isTicker = true
		v, _ := set.Values()[0].(int)
		c.tickerDuration = time.Minute * time.Duration(v)
	}

	c.minutes = set
	return nil
}

func (c *Crontab) parseSecond(arg string) error {
	set, isticker, err := parseArgument(arg, 0, 59, "")
	if err != nil {
		return err
	}

	if isticker {
		c.isTicker = true
		v, _ := set.Values()[0].(int)
		c.tickerDuration = time.Second * time.Duration(v)
	}

	c.seconds = set
	return nil
}

//解析 分，时，天，月，星期
func parseArgument(arg string, unitMinValue, unitMaxValue int, unitTye string) (result *collections.Set, isTicker bool, err error) {
	result = collections.NewSet(unitMaxValue - unitMinValue + 1)

	monthesAbbr := map[string]int{"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6, "Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12}

	monthes := map[string]int{"January": 1, "February": 2, "March": 3, "April": 4,
		"May": 5, "June": 6, "July": 7, "Auguet": 8, "September": 9, "October": 10, "November": 11, "December": 12}

	weekdays := map[string]int{"Sunday": 0, "Monday": 1, "Tuesday": 2, "Wednesday": 3, "Thursday": 4, "Friday": 5, "Saturday": 6}
	weekdaysAbbr := map[string]int{"Sun": 0, "Mon": 1, "Tue": 2, "Wed": 3, "Thu": 4, "Fri": 5, "Sat": 6}

	//先根据comma分割， 因为有 2-4,6格式
	values := strings.Split(arg, ",")
	for _, param := range values {
		var startVal, endVal int
		var convErr error

		if pos := strings.Index(param, "-"); pos != -1 { //类似 1-4 的格式
			startParam := param[:pos]
			endParam := param[pos+1:]

			if unitTye == "month" {
				if montheIndex, ok := monthes[strings.Title(startParam)]; ok { //是英文书写形式
					startVal = montheIndex
					endVal = monthes[endParam]
				} else if montheIndex, ok = monthesAbbr[strings.Title(startParam)]; ok { //是英文简写形式
					startVal = montheIndex
					endVal = monthesAbbr[strings.Title(endParam)]
				} else { //是数字
					startVal, convErr = strconv.Atoi(startParam)
					if convErr != nil {
						err = errors.New(FORMAT_ERROR + " at:" + startParam)
						return
					}

					endVal, convErr = strconv.Atoi(endParam)
					if convErr != nil {
						err = errors.New(FORMAT_ERROR + " at:" + endParam)
						return
					}
				}

			} else if unitTye == "week" {
				if weekIndex, ok := weekdays[strings.Title(startParam)]; ok { //是英文书写形式
					startVal = weekIndex
					endVal = weekdays[endParam]
				} else if weekIndex, ok = weekdaysAbbr[strings.Title(startParam)]; ok { //是英文简写形式
					startVal = weekIndex
					endVal = weekdaysAbbr[strings.Title(endParam)]
				} else { //是数字
					startVal, convErr = strconv.Atoi(startParam)
					if convErr != nil {
						err = errors.New(FORMAT_ERROR + " at:" + startParam)
						return
					}

					endVal, convErr = strconv.Atoi(endParam)
					if convErr != nil {
						err = errors.New(FORMAT_ERROR + " at:" + endParam)
						return
					}
				}
			} else { //不是week和month
				startVal, convErr = strconv.Atoi(startParam)
				if convErr != nil {
					err = errors.New(FORMAT_ERROR + " at:" + startParam)
					return
				}

				endVal, convErr = strconv.Atoi(endParam)
				if convErr != nil {
					err = errors.New(FORMAT_ERROR + " at:" + endParam)
					return
				}
			}

			if endVal < startVal {
				for i := startVal; i <= unitMaxValue; i++ {
					result.Add(i)
				}

				for i := unitMinValue; i <= endVal; i++ {
					result.Add(i)
				}
			} else {
				for i := startVal; i <= endVal; i++ {
					result.Add(i)
				}
			}
		} else if param == "*" {
			for i := unitMinValue; i <= unitMaxValue; i++ {
				result.Add(i)
			}
			return
		} else if idx := strings.Index(param, "*/"); idx != -1 {
			value, convErr := strconv.Atoi(param[idx+2:])
			if convErr != nil {
				err = convErr
				return
			}

			isTicker = true
			result.Add(value)
			return
		} else { //是 单个数字或者单个单词
			if unitTye == "month" {
				if montheIndex, ok := monthes[strings.Title(param)]; ok { //是英文书写形式
					result.Add(montheIndex)
				} else if montheIndex, ok := monthesAbbr[strings.Title(param)]; ok { //是英文书写形式
					result.Add(montheIndex)
				}
			} else if unitTye == "week" {
				if weekIndex, ok := weekdays[strings.Title(param)]; ok { //是英文书写形式
					result.Add(weekIndex)
				} else if weekIndex, ok := weekdaysAbbr[strings.Title(param)]; ok { //是英文书写形式
					result.Add(weekIndex)
				}
			} else {
				convValue, convErr := strconv.Atoi(param)
				if convErr != nil {
					err = errors.New(FORMAT_ERROR + " at:" + param)
					return
				}
				result.Add(convValue)
			}
		}
	}

	for _, v := range result.Values() {
		val, _ := v.(int)
		if val < unitMinValue || val > unitMaxValue {
			err = errors.New(FORMAT_ERROR + " invalid range")
			return
		}
	}

	return
}

func execHandler(handler Handler) {
	fn := reflect.TypeOf(handler)
	params := make([]reflect.Value, fn.NumIn())
	for i := 0; i < fn.NumIn(); i++ {
		params = append(params, reflect.ValueOf(fn.In(i)))
	}

	value := reflect.ValueOf(handler)
	value.Call(params)
}

func execHandlers(handlers ...Handler) {
	for _, handler := range handlers {
		go execHandler(handler)
	}
}
