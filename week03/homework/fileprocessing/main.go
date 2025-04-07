package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 定义单词结构体
type WordData struct {
	Word         string        `json:"word"`
	Translations []Translation `json:"translations"`
	Phrases      []Phrase      `json:"phrases"`
}

// 定义Translation结构体
type Translation struct {
	Translation string `json:"translation"`
	Type        string `json:"type"`
}

// 定义Phrase结构体
type Phrase struct {
	Phrase      string `json:"phrase"`
	Translation string `json:"translation"`
}

// 数据库连接
func openDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库：%s失败: %w", dbPath, err)
	}
	return db, nil
}

// 获取单词 ID（如果不存在则插入）
func getWord(tx *sql.Tx, word string) int {
	var wordID int
	err := tx.QueryRow("SELECT id FROM words WHERE word = ?", word).Scan(&wordID)
	if err == sql.ErrNoRows {
		result, err := tx.Exec("INSERT INTO words (word) VALUES (?)", word)
		if err != nil {
			return 0
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return 0
		}
		return int(lastInsertID)
	} else if err != nil {
		return 0
	}
	return wordID
}

// 插入翻译数据
func insertTranslations(tx *sql.Tx, wordID int, translations []Translation) error {
	//使用prepare对插入Translations的SQL语句进行预编译
	stmt, _ := tx.Prepare("INSERT OR IGNORE INTO translations (word_id, translation, type) VALUES (?, ?, ?)")
	defer stmt.Close()

	for _, t := range translations {
		if _, err := stmt.Exec(wordID, t.Translation, t.Type); err != nil {
			return fmt.Errorf("无法插入translation：%w", err)
		}
	}
	return nil
}

func insertPhrases(tx *sql.Tx, wordID int, phrases []Phrase) error {
	//使用prepare对插入Phrases的SQL语句进行预编译
	stmt, _ := tx.Prepare("INSERT OR IGNORE INTO phrases (word_id, phrase, translation) VALUES (?, ?, ?)")
	defer stmt.Close()

	for _, p := range phrases {
		if _, err := stmt.Exec(wordID, p.Phrase, p.Translation); err != nil {
			return fmt.Errorf("无法插入phrase：%w", err)
		}
	}
	return nil
}

// 处理单个JSON 文件
func processJson(db *sql.DB, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("无法读取文件%s：%w", filePath, err)
	}

	// 将json文件转换为结构体并存入wordList中
	var wordList []WordData
	if err := json.Unmarshal(data, &wordList); err != nil {
		return fmt.Errorf("无法将Json文件%s反序列化：%w", filePath, err)
	}

	// 使用事务处理减少运行时间
	batchSize := 1000
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("无法开始事务: %w", err)
	}

	for i, wordData := range wordList {
		// 获取单词 ID
		wordID := getWord(tx, wordData.Word)
		if wordID == 0 {
			return fmt.Errorf("获取单词'%s'ID失败", wordData.Word)
		}

		//插入Translations
		if err := insertTranslations(tx, wordID, wordData.Translations); err != nil {
			return fmt.Errorf("无法为单词'%s'插入Translations：%w", wordData.Word, err)
		}

		//插入Phrases
		if err := insertPhrases(tx, wordID, wordData.Phrases); err != nil {
			return fmt.Errorf("无法为单词'%s'插入Phrases：%w", wordData.Word, err)
		}

		// 每1000条数据提交一次事务
		if (i+1)%batchSize == 0 || i == len(wordList)-1 {
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("无法提交事务：%w", err)
			}

			// 开始新的事务
			tx, err = db.Begin()
			if err != nil {
				return fmt.Errorf("无法开始事务：%w", err)
			}
		}
	}

	return nil
}

func main() {
	db, err := openDatabase("./data.db")
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	jsonFiles := []string{
		//	"./json/1-初中-顺序.json",
		//	"./json/2-高中-顺序.json",
		"./json/3-CET4-顺序.json",
		"./json/4-CET6-顺序.json",
		//	"./json/5-考研-顺序.json",
		//	"./json/6-托福-顺序.json",
		//"./json/7-SAT-顺序.json",
	}
	//开始计时
	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeStart)
		fmt.Printf("处理所有文件共耗时：%v", duration.Seconds())
	}()

	// 分批处理每个JSON 文件
	for x, filePath := range jsonFiles {
		timeStart := time.Now()
		if err := processJson(db, filePath); err != nil {
			log.Printf("处理第%v个文件出现错误: %v", x+1, err)
		}
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeStart)
		fmt.Printf("第%v个Json文件插入花费时间: %v\n", x+1, duration.Seconds())
	}

}
