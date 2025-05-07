import { useState, useEffect } from 'react'
import { Card, Typography, Spin, Alert } from 'antd'
import axios from 'axios'
import SimpleMarkdown from '../components/SimpleMarkdown'
import '../styles/StudyNotes.css'

const { Title } = Typography

// 从README.md内容中提取学习心得部分
const extractLearningNotes = (content: string): string => {
  const lines = content.split('\n')
  let isInLearningNotes = false
  let learningNotesContent: string[] = []
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i].trim()
    
    if (line === '## 学习心得') {
      isInLearningNotes = true
      continue 
    }
    
    if (isInLearningNotes && line.startsWith('## ') && line !== '## 学习心得') {
      break
    }
    
    if (isInLearningNotes) {
      learningNotesContent.push(lines[i])
    }
  }
  
  return learningNotesContent.join('\n').trim()
}

const StudyNotes = () => {
  const [readmeContent, setReadmeContent] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string>('')

  useEffect(() => {
    const fetchReadmeContent = async () => {
      try {
        // 添加随机查询参数，防止缓存
        const timestamp = new Date().getTime();
        const response = await axios.get(`/README.md?t=${timestamp}`, { 
          headers: { 'Content-Type': 'text/plain' },
          transformResponse: [(data) => data] 
        })
        
        const learningNotesPart = extractLearningNotes(response.data)
        
        if (!learningNotesPart) {
          setError('未找到学习心得部分，请确保README.md文件中包含"## 学习心得"部分。')
        } else {
          setReadmeContent(learningNotesPart)
        }
        
        setLoading(false)
      } catch (err) {
        console.error('获取README.md失败:', err)
        setError('无法加载学习心得内容，请确保README.md文件存在。')
        setLoading(false)
      }
    }

    fetchReadmeContent()
  }, [])

  if (loading) {
    return (
      <div className="study-notes-loading">
        <Spin tip="加载中..."/>
      </div>
    )
  }

  if (error) {
    return (
      <div className="study-notes-container">
        <Card className="study-notes-card">
          <Alert
            message="加载错误"
            description={error}
            type="error"
            showIcon
          />
        </Card>
      </div>
    )
  }

  return (
    <div className="study-notes-container">
      <Card className="study-notes-card">
        <Title level={2} className="study-notes-title">学习心得</Title>
        <div className="markdown-content">
          <SimpleMarkdown>{readmeContent}</SimpleMarkdown>
        </div>
      </Card>
    </div>
  )
}

export default StudyNotes