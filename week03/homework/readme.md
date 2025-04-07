此目录存放本周课后作业，可以在此文件添加作业设计思路和流程图等

本次作业根据Json文件中的格式将数据库(data.db)中的表格式划分为word(id,word);translations(id,word_id,translation,type);phrases(id,word_id,phrase.translation)。

其中具体的表构建sql语句为：

```sql
('''
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY,
    word TEXT UNIQUE NOT NULL
)
''')
('''
CREATE TABLE IF NOT EXISTS translations (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    translation TEXT NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words (id),
    UNIQUE(word_id, translation, type) ON CONFLICT IGNORE
)
''')
('''
CREATE TABLE IF NOT EXISTS phrases (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    phrase TEXT NOT NULL,
    translation TEXT NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words (id),
    UNIQUE(word_id, phrase, translation) ON CONFLICT IGNORE
)
''')

```