# 数据库

本次作业根据Json文件中的格式将数据库(data.db)划分为三张表：

- **word**(id, word)

`id`:主键，自增

`word`:单词名称，唯一且不为空

- **translations**(id, word_id, translation, type)

`id`:主键，自增

`word_id`: 外键，指向 words 表的 id

`translation`: 单词翻译内容，不为空

`type`: 单词翻译的词性，不为空

- **phrases**(id, word_id, phrase, translation)

`id`:主键，自增

`word_id`: 外键，指向 words 表的 id

`phrase`: 单词的相关短语内容，不为空

`translation`: 单词的相关短语翻译，不为空

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

