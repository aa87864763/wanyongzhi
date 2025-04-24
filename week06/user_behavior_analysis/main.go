package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// UserAction 表示一条用户行为日志
type UserAction struct {
	Timestamp    time.Time
	UserID       string
	ActionType   string
	ActionDetail string
}

// UserStats 用户统计信息
type UserStats struct {
	UserID      string
	ActionCount int
	FirstAction time.Time
	LastAction  time.Time
}

// ActionTypeStats 行为类型统计
type ActionTypeStats struct {
	ActionType string
	Count      int
}

// TimeWindowStats 时间窗口统计
type TimeWindowStats struct {
	TimeWindow   time.Time
	ActiveUsers  int
	TotalActions int
}

func ParseLogLine(line string) (*UserAction, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid log format: %s", line)
	}

	// 解析时间戳
	timestamp, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp: %v", err)
	}

	return &UserAction{
		Timestamp:    timestamp,
		UserID:       strings.TrimSpace(parts[1]),
		ActionType:   strings.TrimSpace(parts[2]),
		ActionDetail: strings.TrimSpace(parts[3]),
	}, nil
}

// ReadLogFile 读取并解析日志文件
func ReadLogFile(filepath string) ([]UserAction, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var actions []UserAction
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		action, err := ParseLogLine(scanner.Text())
		if err != nil {
			fmt.Printf("Warning: Skipping line due to error: %v\n", err)
			continue
		}
		actions = append(actions, *action)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return actions, nil
}

// 生成用户统计信息
func generateUserStats(actions []UserAction) []UserStats {
	userStatsMap := make(map[string]*UserStats)

	for _, action := range actions {
		if stats, exists := userStatsMap[action.UserID]; exists {
			stats.ActionCount++
			if action.Timestamp.Before(stats.FirstAction) {
				stats.FirstAction = action.Timestamp
			}
			if action.Timestamp.After(stats.LastAction) {
				stats.LastAction = action.Timestamp
			}
		} else {
			userStatsMap[action.UserID] = &UserStats{
				UserID:      action.UserID,
				ActionCount: 1,
				FirstAction: action.Timestamp,
				LastAction:  action.Timestamp,
			}
		}
	}

	// 转换为切片并排序
	userStats := make([]UserStats, 0, len(userStatsMap))
	for _, stats := range userStatsMap {
		userStats = append(userStats, *stats)
	}
	sort.Slice(userStats, func(i, j int) bool {
		return userStats[i].ActionCount > userStats[j].ActionCount
	})

	return userStats
}

// 生成行为类型统计
func generateActionTypeStats(actions []UserAction) []ActionTypeStats {
	actionTypeMap := make(map[string]int)

	for _, action := range actions {
		actionTypeMap[action.ActionType]++
	}

	// 转换为切片并排序
	actionStats := make([]ActionTypeStats, 0, len(actionTypeMap))
	for actionType, count := range actionTypeMap {
		actionStats = append(actionStats, ActionTypeStats{
			ActionType: actionType,
			Count:      count,
		})
	}
	sort.Slice(actionStats, func(i, j int) bool {
		return actionStats[i].Count > actionStats[j].Count
	})

	return actionStats
}

// 生成时间窗口统计
func generateTimeWindowStats(actions []UserAction) []TimeWindowStats {
	timeWindowMap := make(map[time.Time]map[string]bool)
	actionCountMap := make(map[time.Time]int)

	for _, action := range actions {
		// 将时间戳规整到分钟
		windowTime := time.Date(
			action.Timestamp.Year(),
			action.Timestamp.Month(),
			action.Timestamp.Day(),
			action.Timestamp.Hour(),
			action.Timestamp.Minute(),
			0, 0, action.Timestamp.Location(),
		)

		if _, exists := timeWindowMap[windowTime]; !exists {
			timeWindowMap[windowTime] = make(map[string]bool)
		}
		timeWindowMap[windowTime][action.UserID] = true
		actionCountMap[windowTime]++
	}

	// 转换为切片并排序
	timeStats := make([]TimeWindowStats, 0, len(timeWindowMap))
	for windowTime, users := range timeWindowMap {
		timeStats = append(timeStats, TimeWindowStats{
			TimeWindow:   windowTime,
			ActiveUsers:  len(users),
			TotalActions: actionCountMap[windowTime],
		})
	}
	sort.Slice(timeStats, func(i, j int) bool {
		return timeStats[i].TimeWindow.Before(timeStats[j].TimeWindow)
	})

	return timeStats
}

// 保存用户统计到CSV
func saveUserStatsToCSV(stats []UserStats, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"用户ID", "操作次数", "首次操作时间", "最后操作时间"})

	// 写入数据
	for _, stat := range stats {
		writer.Write([]string{
			stat.UserID,
			fmt.Sprintf("%d", stat.ActionCount),
			stat.FirstAction.Format("2006-01-02 15:04:05"),
			stat.LastAction.Format("2006-01-02 15:04:05"),
		})
	}
	return nil
}

// 保存行为类型统计到CSV
func saveActionTypeStatsToCSV(stats []ActionTypeStats, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"行为类型", "总次数"})

	// 写入数据
	for _, stat := range stats {
		writer.Write([]string{
			stat.ActionType,
			fmt.Sprintf("%d", stat.Count),
		})
	}
	return nil
}

// 保存时间窗口统计到CSV
func saveTimeWindowStatsToCSV(stats []TimeWindowStats, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"时间段", "活跃用户数", "操作总数"})

	// 写入数据
	for _, stat := range stats {
		writer.Write([]string{
			stat.TimeWindow.Format("2006-01-02 15:04"),
			fmt.Sprintf("%d", stat.ActiveUsers),
			fmt.Sprintf("%d", stat.TotalActions),
		})
	}
	return nil
}

func main() {
	// 读取日志文件
	actions, err := ReadLogFile("user_actions.log")
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		return
	}

	// 生成统计信息
	userStats := generateUserStats(actions)
	actionTypeStats := generateActionTypeStats(actions)
	timeWindowStats := generateTimeWindowStats(actions)

	// 保存统计结果到CSV文件
	if err := saveUserStatsToCSV(userStats, "user_statistics.csv"); err != nil {
		fmt.Printf("Error saving user statistics: %v\n", err)
	}

	if err := saveActionTypeStatsToCSV(actionTypeStats, "action_statistics.csv"); err != nil {
		fmt.Printf("Error saving action type statistics: %v\n", err)
	}

	if err := saveTimeWindowStatsToCSV(timeWindowStats, "minute_statistics.csv"); err != nil {
		fmt.Printf("Error saving time window statistics: %v\n", err)
	}
}
