# 数据库

本次作业根据Json文件中的格式将数据库(data.db)划分为三张表：

- **word**(id, word)
- **translations**(id, word_id, translation, type)
- **phrases**(id, word_id, phrase, translation)

这些表通过外键关联，其中 `word_id` 是 `translations` 和 `phrases` 表的外键，指向 `words` 表的主键 `id`。

其中关于表构建的sql语句为：

```sql
'''
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY,
    word TEXT UNIQUE NOT NULL
)
'''
'''
CREATE TABLE IF NOT EXISTS translations (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    translation TEXT NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words (id),
    UNIQUE(word_id, translation, type) ON CONFLICT IGNORE
)
'''
'''
CREATE TABLE IF NOT EXISTS phrases (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    phrase TEXT NOT NULL,
    translation TEXT NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words (id),
    UNIQUE(word_id, phrase, translation) ON CONFLICT IGNORE
)
'''
```

# 代码实现流程

