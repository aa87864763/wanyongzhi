import { useState, useEffect } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { 
  Card, 
  Typography, 
  Button, 
  Space, 
  Tag, 
  List, 
  Radio, 
  Checkbox, 
  message, 
  Divider,
  Empty,
  Alert
} from 'antd'
import type { AIResponse, QuestionRequest, QuestionData } from '../types'
import { QuestionType, QuestionDifficulty } from '../types'
import { addQuestion } from '../services/api'
import '../styles/AIPreview.css'

const { Title, Paragraph } = Typography

const AIPreview = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [aiResponses, setAiResponses] = useState<AIResponse[]>([])
  const [request, setRequest] = useState<QuestionRequest | null>(null)
  const [selectedQuestions, setSelectedQuestions] = useState<number[]>([])

  useEffect(() => {
    // 从location state获取数据
    const state = location.state as { aiRes: AIResponse | AIResponse[]; request: QuestionRequest } | null
    if (state?.aiRes) {
      // 处理单个题目或多个题目的情况
      if (Array.isArray(state.aiRes)) {
        setAiResponses(state.aiRes)
      } else {
        setAiResponses([state.aiRes])
      }
      setRequest(state.request)
    } else {
      message.error('没有需要预览的题目数据')
      navigate('/create')
    }
  }, [location, navigate])

  // 添加选中的题目到数据库
  const handleAddSelectedQuestions = async () => {
    if (!request || aiResponses.length === 0) {
      message.error('题目数据不完整')
      return
    }

    // 检查是否有选中的题目
    if (selectedQuestions.length === 0) {
      message.warning('请至少选择一个题目添加到题库')
      return
    }

    setLoading(true)
    try {
      const currentTime = new Date().toISOString()
      const promises = selectedQuestions.map(index => {
        const aiRes = aiResponses[index]
        // 构建题目数据
        const questionData: QuestionData = {
          aiStartTime: currentTime,
          aiEndTime: currentTime,
          aiCostTime: 0,
          aiReq: request,
          aiRes: aiRes,
          difficulty: request.difficulty || QuestionDifficulty.Medium,
          createdAt: currentTime
        }

        return addQuestion(questionData)
      })

      const results = await Promise.all(promises)
      const failures = results.filter(res => res.code !== 0)
      
      if (failures.length === 0) {
        message.success(`成功添加${selectedQuestions.length}个题目到题库`)
        navigate('/question-list')
      } else {
        message.error(`有${failures.length}个题目添加失败`)
      }
    } catch (err) {
      console.error('添加题目失败:', err)
      message.error('添加题目失败，请检查网络连接或服务器状态')
    } finally {
      setLoading(false)
    }
  }

  // 返回出题页面
  const handleBack = () => {
    navigate('/create')
  }

  // 处理题目选择状态变更
  const handleQuestionSelect = (index: number, checked: boolean) => {
    if (checked) {
      setSelectedQuestions(prev => [...prev, index])
    } else {
      setSelectedQuestions(prev => prev.filter(i => i !== index))
    }
  }

  // 渲染选择题选项
  const renderChoiceOptions = (answer: string[], right: number[]) => {
    return (
      <div className="options-container">
        <List
          dataSource={answer}
          renderItem={(item, index) => {
            const isCorrect = right.includes(index)
            return (
              <List.Item className={isCorrect ? 'correct-option' : ''}>
                <Space>
                  {request?.type === QuestionType.SingleChoice ? (
                    <Radio checked={isCorrect} disabled />
                  ) : (
                    <Checkbox checked={isCorrect} disabled />
                  )}
                  <div dangerouslySetInnerHTML={{ __html: item }} />
                </Space>
                {isCorrect && <Tag color="success">正确答案</Tag>}
              </List.Item>
            )
          }}
        />
      </div>
    )
  }

  const renderProgrammingCode = (code: string, languages: string[]) => {
    return (
      <div className="code-container">
        {languages.map((lang, index) => (
          <div key={index}>
            <Tag color="blue">{lang}</Tag>
            <pre className="code-block">
              <code>{code}</code>
            </pre>
          </div>
        ))}
      </div>
    )
  }

  // 获取题目类型文本
  const getQuestionTypeText = (type?: QuestionType) => {
    switch (type) {
      case QuestionType.SingleChoice:
        return '单选题'
      case QuestionType.MultiChoice:
        return '多选题'
      case QuestionType.Programming:
        return '编程题'
      default:
        return '未知类型'
    }
  }

  // 获取难度文本和颜色
  const getDifficultyInfo = (difficulty?: QuestionDifficulty) => {
    switch (difficulty) {
      case QuestionDifficulty.Easy:
        return { text: '简单', color: 'green' }
      case QuestionDifficulty.Medium:
        return { text: '中等', color: 'orange' }
      case QuestionDifficulty.Hard:
        return { text: '困难', color: 'red' }
      default:
        return { text: '中等', color: 'orange' }
    }
  }

  if (aiResponses.length === 0) {
    return <Empty description="没有题目数据" />
  }

  const difficultyInfo = getDifficultyInfo(request?.difficulty)

  return (
    <div className="ai-preview">
      <div className="preview-header">
        <Space style={{ marginBottom: 16 }}>
          <Tag color="blue">{getQuestionTypeText(request?.type)}</Tag>
          <Tag color={difficultyInfo.color}>{difficultyInfo.text}</Tag>
        </Space>
        <Alert
          message="题目预览"
          description={`共生成了 ${aiResponses.length} 道题目，请选择要添加到题库的题目`}
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
        />
      </div>

      {aiResponses.map((aiRes, index) => {
        const isProgramming = request?.type === QuestionType.Programming || aiRes.code;
        return (
          <Card
            key={index}
            title={`题目 ${index + 1}`}
            className="preview-card"
            extra={
              <Checkbox 
                checked={selectedQuestions.includes(index)}
                onChange={(e) => handleQuestionSelect(index, e.target.checked)}
              >
                选择此题
              </Checkbox>
            }
            style={{ marginBottom: 16 }}
          >
            <div className="question-content">
              <Title level={4}>题目</Title>
              <Paragraph>
                <div dangerouslySetInnerHTML={{ __html: aiRes.title.replace(/\n/g, '<br/>') }} />
              </Paragraph>

              <Divider />

              <Title level={4}>答案</Title>
              {isProgramming ? (
                renderProgrammingCode(aiRes.code || '', request?.language ? [request.language] : [])
              ) : (
                renderChoiceOptions(aiRes.answer, aiRes.right)
              )}
            </div>
          </Card>
        )
      })}

      <div className="action-buttons">
        <Space>
          <Button 
            type="primary" 
            onClick={handleAddSelectedQuestions} 
            loading={loading}
            disabled={selectedQuestions.length === 0}
          >
            添加已选题目({selectedQuestions.length})到题库
          </Button>
          <Button onClick={handleBack}>
            返回
          </Button>
        </Space>
      </div>
    </div>
  )
}

export default AIPreview