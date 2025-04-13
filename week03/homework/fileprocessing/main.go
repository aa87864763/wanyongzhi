package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	//如果显示2025/04/11 20:15:50 打开数据库：./data.db失败
	//那么使用下面的sqlite驱动而不是go-sqlite3
	//同时在initDatabase将db, err := sql.Open("sqlite3", dbPath)中的sqlite3更改为sqlite
	//	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
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

// 定义预编译语句
var (
	insertWordStmt         *sql.Stmt
	insertTranslationsStmt *sql.Stmt
	insertPhrasesStmt      *sql.Stmt
)

// 初始化数据库
func initDatabase(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("打开数据库：%s失败", dbPath)
	}

	// 设置WAL模式
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatalf("无法设置WAL模式：%v", err)
	}

	createTable(db, "words", `
	CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY,
		word TEXT UNIQUE NOT NULL
	)`)

	createTable(db, "translations", `
	CREATE TABLE IF NOT EXISTS translations (
		id INTEGER PRIMARY KEY,
		word_id INTEGER NOT NULL,
		translation TEXT NOT NULL,
		type TEXT NOT NULL,
		FOREIGN KEY (word_id) REFERENCES words (id),
		UNIQUE(word_id, translation, type) ON CONFLICT IGNORE
	)`)

	createTable(db, "phrases", `
	CREATE TABLE IF NOT EXISTS phrases (
		id INTEGER PRIMARY KEY,
		word_id INTEGER NOT NULL,
		phrase TEXT NOT NULL,
		translation TEXT NOT NULL,
		FOREIGN KEY (word_id) REFERENCES words (id),
		UNIQUE(word_id, phrase, translation) ON CONFLICT IGNORE
	)`)
	return db
}

// 判断表的存在与否进行创建
func createTable(db *sql.DB, tableName, createSQL string) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?", tableName).Scan(&count)

	if count == 0 {
		_, err := db.Exec(createSQL)
		if err != nil {
			log.Fatalf("无法创建%s表", tableName)
		}
	}
}

// 插入word并获取word_id
func getWord(tx *sql.Tx, word string) int {
	var wordID int
	err := tx.QueryRow("SELECT id FROM words WHERE word = ?", word).Scan(&wordID)
	if err == sql.ErrNoRows {
		result, err := tx.Stmt(insertWordStmt).Exec(word)
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

// 插入Translations数据
func insertTranslations(tx *sql.Tx, wordID int, translations []Translation) error {
	for _, t := range translations {
		_, err := tx.Stmt(insertTranslationsStmt).Exec(wordID, t.Translation, t.Type)
		if err != nil {
			tx.Rollback()
			log.Fatalf("无法插入翻译数据：%v", err)
		}
	}
	return nil
}

// 插入Phrases数据
func insertPhrases(tx *sql.Tx, wordID int, phrases []Phrase) error {
	for _, p := range phrases {
		_, err := tx.Stmt(insertPhrasesStmt).Exec(wordID, p.Phrase, p.Translation)
		if err != nil {
			tx.Rollback()
			log.Fatalf("无法插入短语数据：%v", err)
		}
	}
	return nil
}

// 读取Json文件并反解析为结构体
func readFile(filePath string, wg *sync.WaitGroup, wordLists chan<- []WordData) {
	defer wg.Done()

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("无法读取文件%s：%v", filePath, err)
		return
	}

	var wordList []WordData
	if err := json.Unmarshal(data, &wordList); err != nil {
		log.Printf("无法将Json文件%s反序列化：%v", filePath, err)
		return
	}

	wordLists <- wordList
}

// 使用事务批量写入数据库
func writeDatabase(db *sql.DB, wordLists <-chan []WordData, done chan<- struct{}) {
	defer close(done)

	batchSize := 1000
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("无法开始事务: %v", err)
	}

	for wordList := range wordLists {
		for i, wordData := range wordList {
			// 插入单词
			wordID := getWord(tx, wordData.Word)
			if wordID == 0 {
				log.Fatalf("插入单词'%s'失败：%v", wordData.Word, err)
			}

			if err := insertTranslations(tx, wordID, wordData.Translations); err != nil {
				tx.Rollback()
				log.Fatalf("无法为单词'%s'插入Translations：%v", wordData.Word, err)
			}

			if err := insertPhrases(tx, wordID, wordData.Phrases); err != nil {
				tx.Rollback()
				log.Fatalf("无法为单词'%s'插入Phrases：%v", wordData.Word, err)
			}

			// 每1000条数据提交一次事务
			if (i+1)%batchSize == 0 {
				if err := tx.Commit(); err != nil {
					log.Fatalf("无法提交事务：%v", err)
				}
				tx, err = db.Begin()
				if err != nil {
					log.Fatalf("无法开始事务：%v", err)
				}
			}
		}
	}

	// 提交剩余的数据
	if err := tx.Commit(); err != nil {
		log.Fatalf("无法提交事务：%v", err)
	}
}

func main() {
	db := initDatabase("./data.db")
	defer db.Close()

	// 进行预编译
	var err error
	insertWordStmt, err = db.Prepare("INSERT OR IGNORE INTO words (word) VALUES (?)")
	if err != nil {
		log.Fatalf("无法预编译插入Word语句：%v", err)
	}
	defer insertWordStmt.Close()

	insertTranslationsStmt, err = db.Prepare("INSERT OR IGNORE INTO translations (word_id, translation, type) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("无法预编译插入Translations语句：%v", err)
	}
	defer insertTranslationsStmt.Close()

	insertPhrasesStmt, err = db.Prepare("INSERT OR IGNORE INTO phrases (word_id, phrase, translation) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("无法预编译插入Phrases语句：%v", err)
	}
	defer insertPhrasesStmt.Close()

	jsonFiles := []string{
		//	"./json/1-初中-顺序.json",
		//	"./json/2-高中-顺序.json",
		"./json/3-CET4-顺序.json",
		"./json/4-CET6-顺序.json",
		//	"./json/5-考研-顺序.json",
		//	"./json/6-托福-顺序.json",
		//"./json/7-SAT-顺序.json",
	}

	var wg sync.WaitGroup
	wordLists := make(chan []WordData, len(jsonFiles))
	done := make(chan struct{})

	// 计时
	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeStart)
		fmt.Printf("处理所有文件共花费时间：%v秒\n", duration.Seconds())
	}()

	// 并行读取Json文件
	for _, filePath := range jsonFiles {
		wg.Add(1)
		go readFile(filePath, &wg, wordLists)
	}

	// 并行写入数据库
	go writeDatabase(db, wordLists, done)

	// 等待所有Json文件读取完成
	wg.Wait()
	close(wordLists)

	// 等待数据库写入完成
	<-done
}
