package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Data struct {
	Time       time.Time
	User_id    string
	Action     string
	ActionInfo string
}

type userStats struct {
	UserID      string
	ActionCount int
	FirstAction time.Time
	LastAction  time.Time
}

type actionStats struct {
	Action string
	count  int
}

type minStats struct {
	Time      string
	UserNum   int
	ActionNum int
}

func getData(pathfile string) []Data {
	file, err := os.Open(pathfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var data []Data
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var temp Data
		line := scanner.Text()
		fmt.Println(line)
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			fmt.Println("Invalid log format:", line)
			continue
		}
		temp.Time, _ = time.Parse("2006-01-02 15:04:05", strings.TrimSpace(parts[0]))
		temp.User_id = strings.TrimSpace(parts[1])
		temp.Action = strings.TrimSpace(parts[2])
		temp.ActionInfo = strings.TrimSpace(parts[3])
		data = append(data, temp)
	}
	return data
}

func generateUserStats(data []Data) {
	userMap := make(map[string]userStats)

	for _, d := range data {
		if user, existd := userMap[d.User_id]; !existd {
			userMap[d.User_id] = userStats{
				UserID:      d.User_id,
				ActionCount: 1,
				FirstAction: d.Time,
				LastAction:  d.Time,
			}
		} else { //遇到重复id的时候，更新数据
			user.ActionCount++
			if d.Time.Before(user.FirstAction) {
				user.FirstAction = d.Time
			}
			if d.Time.After(user.LastAction) {
				user.LastAction = d.Time
			}
			userMap[d.User_id] = user
		}
	}

	file, err := os.OpenFile("./user_statistics.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.WriteString("用户ID,操作次数,首次操作时间,最后操作时间\n")
	if err != nil {
		panic(err)
	}
	for _, user := range userMap {
		_, err = writer.WriteString(fmt.Sprintf("%s,%d,%s,%s\n", user.UserID, user.ActionCount, user.FirstAction.Format("2006-01-02 15:04:05"), user.LastAction.Format("2006-01-02 15:04:05")))
	}
}

func generateActionStats(data []Data) {
	var actionMap = make(map[string]int)
	for _, d := range data {
		if _, existd := actionMap[d.Action]; !existd {
			actionMap[d.Action] = 1
		} else {
			actionMap[d.Action]++
		}
	}

	file, err := os.OpenFile("./action_statistics.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.WriteString("操作类型,操作次数\n")
	for action, count := range actionMap {
		_, err = writer.WriteString(fmt.Sprintf("%s,%d\n", action, count))
	}
}

func generateTimeWindowStats(data []Data) {
	var timeMap = make(map[string]minStats)        //时间点对应status
	var userMap = make(map[string]map[string]bool) //时间点对应用户id
	var temp minStats

	for _, d := range data { //对于每一行数据，进行处理
		timeNow := d.Time.Format("2006-01-02 15:04")

		if _, existd := timeMap[timeNow]; !existd {
			//如果没有这个时间点的记录，创建一个新的记录
			//并且将用户id存入userMap中
			timeMap[timeNow] = minStats{
				Time:      timeNow,
				UserNum:   1,
				ActionNum: 1,
			}
			userMap[timeNow] = make(map[string]bool)
			userMap[timeNow][d.User_id] = true
		} else {
			//如果有这个时间点的记录，更新记录
			//并且判断用户id是否已经存在于userMap中，如果不存在，则添加
			//如果存在，则不处理
			temp = timeMap[timeNow]
			temp.ActionNum++
			if _, existd := userMap[timeNow][d.User_id]; !existd {
				temp.UserNum++
				userMap[timeNow][d.User_id] = true
			}
			timeMap[timeNow] = temp
		}
	}

	file, err := os.OpenFile("./minute_statistics.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.WriteString("时间段,活跃用户数,操作次数\n")
	if err != nil {
		panic(err)
	}
	for _, v := range timeMap {
		_, err = writer.WriteString(fmt.Sprintf("%s,%d,%d\n",
			v.Time,
			v.UserNum,
			v.ActionNum))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	data := getData("./user_actions.log")
	generateUserStats(data)
	generateActionStats(data)
	generateTimeWindowStats(data)
}
