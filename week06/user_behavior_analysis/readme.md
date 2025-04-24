将需求构建为结构体，然后将log文件中的内容转换为结构体切片进行保存。然后在遍历结构体切片进行数据的输出。

其中`UserAction`结构体为log内的数据格式，对log文件的内容进行接收

`UserStats`结构体为user_statistics.csv要求的内容

```go
type UserStats struct {
	UserID      string
	ActionCount int
	FirstAction time.Time
	LastAction  time.Time
}
```

`ActionTypeStats`结构体为action_statistics.csv要求的内容

```go
type ActionTypeStats struct {
    ActionType string
    Count int
}
```

`TimeWindowStats`结构体为minute_statistics.csv要求的内容

```go
type ActionTypeStats struct {
    TimeWindow   time.Time
	ActiveUsers  int
	TotalActions int
}
```

其中分为三个func对三个csv输出要求进行统计，分别使用循环对其遍历后进行保存