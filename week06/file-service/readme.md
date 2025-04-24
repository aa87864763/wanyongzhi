# 设计思路

总体思路为：使用Gin框架的RESTful API接口，然后对文件进行处理，之后再对SQLite数据库进行操作。

该项目一共有两个存储内容。第一个为数据库(元数据)，第二个为实际文件，存储在单独的文件夹中。

其中数据库的表结构为：文件唯一id、文件名、文件类型、文件大小、存储路径。具体的建表语句为：

```sql
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE,
    original_name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    storage_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESAMP
)
```

其中文件上传的流程为：发送POST请求，然后验证文件详细信息，再生成UUID和存储路径，然后将这些内容记录为元数据并存入数据库，然后返回上传结果。

文件下载/预览的流程为：发送GET请求(.../files/{uuid})，然后查询文件元数据，返回一个存储路径，再读取内容。

其中api设置为：

文件上传：POST /upload

文件列表：GET /files

文件下载/预览： GET /files/:uuid

统计信息：GET /stats

文件删除 DELETE /files/:uuid

启动端口为8080

将数据库元数据用结构体表现为：

```go
type FileInfo struct {
    UUID string
    OriginalName string
    File
}
```
