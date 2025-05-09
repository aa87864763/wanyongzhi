package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"question-generator/models"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type StorageService struct {
	DataDir string
	DB      *sql.DB
}

// 创建新的存储服务
func NewStorageService() *StorageService {
	dataDir := "./data"
	os.MkdirAll(dataDir, 0755)

	// 打开或创建SQLite数据库
	dbPath := filepath.Join(dataDir, "questions.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}

	// 创建题目表，区分选择题和编程题字段
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		question_type INTEGER NOT NULL, -- 1=单选题, 2=多选题, 3=编程题
		difficulty INTEGER DEFAULT 2, -- 1=简单, 2=中等, 3=困难，默认为中等
		answer TEXT, -- 对于选择题存储选项
		right_answer TEXT -- 对于选择题存储正确答案
	)`)

	if err != nil {
		log.Fatalf("无法创建数据库表: %v", err)
	}

	return &StorageService{
		DataDir: dataDir,
		DB:      db,
	}
}

// 保存问题数据到SQLite数据库
func (s *StorageService) SaveQuestion(data *models.QuestionData) error {
	questionType := data.AIReq.GetQuestionType()

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}

	var stmt *sql.Stmt
	var execParams []interface{}

	if questionType == models.Programming {
		stmt, err = tx.Prepare(`INSERT INTO questions (
			title, question_type, difficulty
		) VALUES (?, ?, ?)`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("准备SQL语句失败: %w", err)
		}

		execParams = []interface{}{
			data.AIRes.Title,
			int(questionType),
			int(data.Difficulty),
		}
	} else {
		answerJSON, err := json.Marshal(data.AIRes.Answer)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("序列化选项失败: %w", err)
		}

		rightJSON, err := json.Marshal(data.AIRes.Right)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("序列化正确答案失败: %w", err)
		}

		stmt, err = tx.Prepare(`INSERT INTO questions (
			title, question_type, difficulty, answer, right_answer
		) VALUES (?, ?, ?, ?, ?)`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("准备SQL语句失败: %w", err)
		}

		execParams = []interface{}{
			data.AIRes.Title,
			int(questionType),
			int(data.Difficulty),
			string(answerJSON),
			string(rightJSON),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(execParams...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("插入数据失败: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

func (s *StorageService) SaveQuestions(questionList []models.QuestionData) error {
	if len(questionList) == 0 {
		return nil
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}

	stmtProgramming, err := tx.Prepare(`INSERT INTO questions (
		title, question_type, difficulty
	) VALUES (?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("准备编程题SQL语句失败: %w", err)
	}
	defer stmtProgramming.Close()

	stmtChoice, err := tx.Prepare(`INSERT INTO questions (
		title, question_type, difficulty, answer, right_answer
	) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("准备选择题SQL语句失败: %w", err)
	}
	defer stmtChoice.Close()

	for _, question := range questionList {
		if question.AIReq.GetQuestionType() == models.Programming {
			_, err = stmtProgramming.Exec(
				question.AIRes.Title,
				int(question.AIReq.GetQuestionType()),
				int(question.Difficulty),
			)
		} else {
			answerJSON, err := json.Marshal(question.AIRes.Answer)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("序列化选项失败: %w", err)
			}

			rightJSON, err := json.Marshal(question.AIRes.Right)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("序列化正确答案失败: %w", err)
			}

			_, err = stmtChoice.Exec(
				question.AIRes.Title,
				int(question.AIReq.GetQuestionType()),
				int(question.Difficulty),
				string(answerJSON),
				string(rightJSON),
			)
		}

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("插入数据失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// 从数据库中获取所有题目
func (s *StorageService) GetAllQuestions() ([]models.QuestionData, error) {
	rows, err := s.DB.Query(`SELECT 
		id, title, question_type, difficulty, answer, right_answer
	FROM questions`)

	if err != nil {
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}

	defer rows.Close()

	var questions []models.QuestionData

	for rows.Next() {
		var q models.QuestionData
		var questionType int
		var difficulty int
		var answerJSON, rightJSON sql.NullString

		err := rows.Scan(
			&q.ID,
			&q.AIRes.Title,
			&questionType,
			&difficulty,
			&answerJSON,
			&rightJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("扫描数据库行失败: %w", err)
		}

		q.AIReq.Type = models.QuestionType(questionType)
		q.Difficulty = models.QuestionDifficulty(difficulty)
		q.AIReq.Difficulty = models.QuestionDifficulty(difficulty)

		if q.AIReq.Type == models.Programming {
		} else {
			// 选择题解析选项和正确答案
			if answerJSON.Valid && answerJSON.String != "" {
				json.Unmarshal([]byte(answerJSON.String), &q.AIRes.Answer)
			}

			if rightJSON.Valid && rightJSON.String != "" {
				json.Unmarshal([]byte(rightJSON.String), &q.AIRes.Right)
			}
		}

		questions = append(questions, q)
	}

	return questions, nil
}

// 查询题目列表，支持分页和条件查询
func (s *StorageService) ListQuestions(page, pageSize int, questionType int, difficulty int, title string) ([]models.QuestionData, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	var conditions []string
	var args []interface{}

	if questionType > 0 {
		conditions = append(conditions, "question_type = ?")
		args = append(args, questionType)
	}

	if difficulty > 0 {
		conditions = append(conditions, "difficulty = ?")
		args = append(args, difficulty)
	}

	if title != "" {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+title+"%")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM questions %s", whereClause)
	var total int
	err := s.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %w", err)
	}

	query := fmt.Sprintf(`SELECT 
		id, title, question_type, difficulty, answer, right_answer
	FROM questions
	%s
	ORDER BY id DESC
	LIMIT ? OFFSET ?`, whereClause)

	queryArgs := append(args, pageSize, offset)

	rows, err := s.DB.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询数据失败: %w", err)
	}
	defer rows.Close()

	var questions []models.QuestionData

	for rows.Next() {
		var q models.QuestionData
		var questionType int
		var answerJSON, rightJSON sql.NullString

		err := rows.Scan(
			&q.ID,
			&q.AIRes.Title,
			&questionType,
			&q.Difficulty,
			&answerJSON,
			&rightJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("扫描数据库行失败: %w", err)
		}

		// 设置题目类型
		q.AIReq.Type = models.QuestionType(questionType)

		if q.AIReq.Type == models.Programming {
		} else {
			// 选择题解析选项和正确答案
			if answerJSON.Valid && answerJSON.String != "" {
				json.Unmarshal([]byte(answerJSON.String), &q.AIRes.Answer)
			}

			if rightJSON.Valid && rightJSON.String != "" {
				json.Unmarshal([]byte(rightJSON.String), &q.AIRes.Right)
			}
		}

		questions = append(questions, q)
	}

	return questions, total, nil
}

// 获取单个题目
func (s *StorageService) GetQuestionByID(id int64) (*models.QuestionData, error) {
	query := `SELECT 
		id, title, question_type, difficulty, answer, right_answer
	FROM questions
	WHERE id = ?`

	var q models.QuestionData
	var questionType int
	var difficulty int
	var answerJSON, rightJSON sql.NullString

	err := s.DB.QueryRow(query, id).Scan(
		&q.ID,
		&q.AIRes.Title,
		&questionType,
		&difficulty,
		&answerJSON,
		&rightJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("题目不存在: ID=%d", id)
		}
		return nil, fmt.Errorf("查询题目失败: %w", err)
	}

	q.AIReq.Type = models.QuestionType(questionType)
	q.Difficulty = models.QuestionDifficulty(difficulty)
	q.AIReq.Difficulty = models.QuestionDifficulty(difficulty)

	// 根据题目类型解析不同字段
	if q.AIReq.Type == models.Programming {
	} else {
		// 选择题解析选项和正确答案
		if answerJSON.Valid && answerJSON.String != "" {
			json.Unmarshal([]byte(answerJSON.String), &q.AIRes.Answer)
		}

		if rightJSON.Valid && rightJSON.String != "" {
			json.Unmarshal([]byte(rightJSON.String), &q.AIRes.Right)
		}
	}

	return &q, nil
}

// 手动添加题目
func (s *StorageService) AddQuestion(data *models.QuestionData) (int64, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("启动事务失败: %w", err)
	}

	var result sql.Result
	var stmt *sql.Stmt

	if data.AIReq.GetQuestionType() == models.Programming {
		stmt, err = tx.Prepare(`INSERT INTO questions (
			title, question_type, difficulty
		) VALUES (?, ?, ?)`)

		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("准备SQL语句失败: %w", err)
		}
		defer stmt.Close()

		result, err = stmt.Exec(
			data.AIRes.Title,
			int(data.AIReq.GetQuestionType()),
			int(data.Difficulty),
		)
	} else {
		answerJSON, err := json.Marshal(data.AIRes.Answer)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("序列化选项失败: %w", err)
		}

		sortedRight := make([]int, len(data.AIRes.Right))
		copy(sortedRight, data.AIRes.Right)
		sort.Ints(sortedRight)

		rightJSON, err := json.Marshal(sortedRight)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("序列化正确答案失败: %w", err)
		}

		stmt, err = tx.Prepare(`INSERT INTO questions (
			title, question_type, difficulty, answer, right_answer
		) VALUES (?, ?, ?, ?, ?)`)

		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("准备SQL语句失败: %w", err)
		}
		defer stmt.Close()

		result, err = stmt.Exec(
			data.AIRes.Title,
			int(data.AIReq.GetQuestionType()),
			int(data.Difficulty),
			string(answerJSON),
			string(rightJSON),
		)
	}

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("插入数据失败: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("提交事务失败: %w", err)
	}

	// 获取新插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	return id, nil
}

// 编辑题目
func (s *StorageService) EditQuestion(id int64, data *models.QuestionData) error {
	var originalType int
	err := s.DB.QueryRow("SELECT question_type FROM questions WHERE id = ?", id).Scan(&originalType)
	if err != nil {
		return fmt.Errorf("获取原始题目信息失败: %w", err)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}

	var stmt *sql.Stmt
	newType := int(data.AIReq.GetQuestionType())

	if newType == int(models.Programming) {
		// 如果是编程题：清空answer和right_answer字段，只更新标题、题目类型和难度
		stmt, err = tx.Prepare(`UPDATE questions SET
			title = ?,
			question_type = ?,
			difficulty = ?,
			answer = NULL,
			right_answer = NULL
		WHERE id = ?`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("准备SQL语句失败: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			data.AIRes.Title,
			newType,
			int(data.Difficulty),
			id,
		)
	} else {
		answerJSON, err := json.Marshal(data.AIRes.Answer)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("序列化选项失败: %w", err)
		}

		sortedRight := make([]int, len(data.AIRes.Right))
		copy(sortedRight, data.AIRes.Right)
		sort.Ints(sortedRight)

		rightJSON, err := json.Marshal(sortedRight)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("序列化正确答案失败: %w", err)
		}

		stmt, err = tx.Prepare(`UPDATE questions SET
			title = ?,
			question_type = ?,
			difficulty = ?,
			answer = ?,
			right_answer = ?
		WHERE id = ?`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("准备SQL语句失败: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			data.AIRes.Title,
			newType,
			int(data.Difficulty),
			string(answerJSON),
			string(rightJSON),
			id,
		)
	}

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据失败: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// 删除题目
func (s *StorageService) DeleteQuestions(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("没有指定要删除的题目ID")
	}

	// 构建占位符
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM questions WHERE id IN (%s)", strings.Join(placeholders, ","))

	result, err := s.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("删除题目失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("没有找到指定的题目")
	}

	return nil
}
