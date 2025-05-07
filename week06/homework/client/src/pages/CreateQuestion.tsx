import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { 
  Card, 
  Form, 
  Select, 
  Button, 
  InputNumber, 
  Radio, 
  Space, 
  message,
  Input,
  Modal,
  Alert
} from 'antd'
import { QuestionType, QuestionDifficulty } from '../types'
import type { QuestionRequest, QuestionData, AIResponse } from '../types'
import { createQuestion, addQuestion } from '../services/api'
import '../styles/CreateQuestion.css'

const { Option } = Select
const { TextArea } = Input

const CreateQuestion = () => {
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [manualAdd, setManualAdd] = useState(false)
  const [currentType, setCurrentType] = useState<QuestionType>(QuestionType.SingleChoice)
  const [errorMsg, setErrorMsg] = useState<string>('')

  // 创建题目（AI生成）
  const handleCreateQuestion = async (values: QuestionRequest) => {
    setLoading(true)
    setErrorMsg('')
    
    try {
      console.log('准备生成题目，参数:', values)
      const res = await createQuestion(values)
      
      if (res.code === 0) {
        message.success('生成题目成功')
        
        let aiResData: AIResponse[] = [];
        
        // 处理不同的响应格式
        if (res.aiRes) {
          // 处理单个题目或多个题目的aiRes字段
          aiResData = Array.isArray(res.aiRes) ? res.aiRes : [res.aiRes];
        } else if (res.questions && res.questions.length > 0) {
          // 处理批量生成的questions字段
          aiResData = res.questions.map(q => q.aiRes);
        }
        
        if (aiResData.length === 0) {
          setErrorMsg('未找到生成的题目数据')
          message.error('未找到生成的题目数据')
          setLoading(false)
          return
        }
        
        // 如果请求了多题但返回单题，记录警告
        if (values.count && values.count > 1 && aiResData.length === 1) {
          console.warn('服务器只返回了一个题目，但请求了多个题目');
        }
        
        // 导航到预览页面，传递生成的题目数据
        navigate('/preview', { 
          state: { 
            aiRes: aiResData, 
            request: values 
          } 
        })
      } else {
        setErrorMsg(res.msg || '生成题目失败，请稍后重试')
        message.error(res.msg || '生成题目失败')
      }
    } catch (err) {
      console.error('生成题目错误:', err)
      setErrorMsg('生成题目失败，请检查网络连接或服务器状态')
      message.error('生成题目失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  // 手动添加题目
  const handleManualAdd = async (values: any) => {
    setLoading(true)
    setErrorMsg('')
    
    try {
      // 处理选项和正确答案
      let answer: string[] = []
      let right: number[] = []
      
      if (values.type !== QuestionType.Programming) {
        // 处理选择题数据
        answer = values.options.split('\n').filter(Boolean)
        right = values.rightAnswers.split(',').map(Number)
        
        // 验证答案是否有效
        if (right.some(idx => idx < 0 || idx >= answer.length)) {
          setErrorMsg(`正确答案索引无效，请确保索引在0到${answer.length - 1}之间`)
          message.error(`正确答案索引无效，请确保索引在0到${answer.length - 1}之间`)
          setLoading(false)
          return
        }
      }
      
      // 构建题目数据
      const currentTime = new Date().toISOString()
      const questionData: QuestionData = {
        aiStartTime: currentTime,
        aiEndTime: currentTime,
        aiCostTime: 0,
        aiReq: {
          type: values.type,
          difficulty: values.difficulty,
          language: values.language
        },
        aiRes: {
          title: values.title,
          answer: answer,
          right: right,
          code: values.type === QuestionType.Programming ? values.code : undefined
        },
        difficulty: values.difficulty,
        createdAt: currentTime
      }
      
      console.log('准备手动添加题目:', questionData)
      
      const res = await addQuestion(questionData)
      if (res.code === 0) {
        message.success('添加题目成功')
        form.resetFields()
        setManualAdd(false)
        // 重定向到题目列表
        navigate('/question-list')
      } else {
        setErrorMsg(res.msg || '添加题目失败')
        message.error(res.msg || '添加题目失败')
      }
    } catch (err) {
      console.error('手动添加题目错误:', err)
      setErrorMsg('添加题目失败，请检查网络连接或服务器状态')
      message.error('添加题目失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }
  
  // 处理表单提交
  const handleFinish = (values: any) => {
    if (manualAdd) {
      handleManualAdd(values)
    } else {
      handleCreateQuestion(values)
    }
  }
  
  // 切换题目类型
  const handleTypeChange = (value: QuestionType) => {
    setCurrentType(value)
  }
  
  const renderAIGenerateForm = () => (
    <Form.Provider>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleFinish}
        initialValues={{
          type: QuestionType.SingleChoice,
          difficulty: QuestionDifficulty.Medium,
          count: 1,
          languages: ["go"]
        }}
      >
        {errorMsg && (
          <Alert
            message="操作错误"
            description={errorMsg}
            type="error"
            showIcon
            closable
            onClose={() => setErrorMsg('')}
            style={{ marginBottom: 16 }}
          />
        )}

        <Form.Item
          name="type"
          label="题目类型"
          rules={[{ required: true, message: '请选择题目类型' }]}
        >
          <Radio.Group onChange={(e) => handleTypeChange(e.target.value)}>
            <Radio.Button value={QuestionType.SingleChoice}>单选题</Radio.Button>
            <Radio.Button value={QuestionType.MultiChoice}>多选题</Radio.Button>
            <Radio.Button value={QuestionType.Programming}>编程题</Radio.Button>
          </Radio.Group>
        </Form.Item>

        <Form.Item
          name="difficulty"
          label="难度"
          rules={[{ required: true, message: '请选择难度' }]}
        >
          <Radio.Group>
            <Radio.Button value={QuestionDifficulty.Easy}>简单</Radio.Button>
            <Radio.Button value={QuestionDifficulty.Medium}>中等</Radio.Button>
            <Radio.Button value={QuestionDifficulty.Hard}>困难</Radio.Button>
          </Radio.Group>
        </Form.Item>

        <Form.Item
          name="languages"
          label="编程语言"
          rules={[{ required: true, message: '请选择至少一种编程语言' }]}
          tooltip="选择相关领域的编程语言，这将影响题目内容领域"
        >
          <Select mode="multiple" placeholder="选择编程语言">
            <Option value="go">Go</Option>
            <Option value="java">Java</Option>
            <Option value="python">Python</Option>
            <Option value="c++">C++</Option>
            <Option value="javascript">JavaScript</Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="count"
          label="生成数量"
          rules={[{ required: true, message: '请输入数量' }]}
        >
          <InputNumber min={1} max={10} />
        </Form.Item>

        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>
              AI生成题目
            </Button>
            <Button onClick={() => setManualAdd(true)}>
              手工出题
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Form.Provider>
  )
  
  // 手动添加题目表单
  const renderManualAddForm = () => (
    <Modal
      title="手工出题"
      open={manualAdd}
      onCancel={() => setManualAdd(false)}
      footer={null}
      width={700}
    >
      <Form
        layout="vertical"
        onFinish={handleFinish}
        initialValues={{
          type: QuestionType.SingleChoice,
          difficulty: QuestionDifficulty.Medium,
          language: "go"
        }}
      >
        {errorMsg && (
          <Alert
            message="操作错误"
            description={errorMsg}
            type="error"
            showIcon
            closable
            onClose={() => setErrorMsg('')}
            style={{ marginBottom: 16 }}
          />
        )}
        
        <Form.Item
          name="type"
          label="题目类型"
          rules={[{ required: true, message: '请选择题目类型' }]}
        >
          <Radio.Group onChange={(e) => handleTypeChange(e.target.value)}>
            <Radio.Button value={QuestionType.SingleChoice}>单选题</Radio.Button>
            <Radio.Button value={QuestionType.MultiChoice}>多选题</Radio.Button>
            <Radio.Button value={QuestionType.Programming}>编程题</Radio.Button>
          </Radio.Group>
        </Form.Item>
        
        <Form.Item
          name="difficulty"
          label="难度"
          rules={[{ required: true, message: '请选择难度' }]}
        >
          <Radio.Group>
            <Radio.Button value={QuestionDifficulty.Easy}>简单</Radio.Button>
            <Radio.Button value={QuestionDifficulty.Medium}>中等</Radio.Button>
            <Radio.Button value={QuestionDifficulty.Hard}>困难</Radio.Button>
          </Radio.Group>
        </Form.Item>
        
        <Form.Item
          name="language"
          label="编程语言"
          rules={[{ required: true, message: '请选择编程语言' }]}
          initialValue="go"
          tooltip={currentType !== QuestionType.Programming ? "选择相关领域的编程语言，这将影响选择题的内容领域" : "选择编程题使用的语言"}
        >
          <Select>
            <Option value="go">Go</Option>
            <Option value="java">Java</Option>
            <Option value="python">Python</Option>
            <Option value="c++">C++</Option>
            <Option value="javascript">JavaScript</Option>
          </Select>
        </Form.Item>
        
        <Form.Item
          name="title"
          label="题目标题"
          rules={[{ required: true, message: '请输入题目标题' }]}
        >
          <TextArea rows={4} placeholder="请输入题目内容" />
        </Form.Item>
        
        {currentType !== QuestionType.Programming ? (
          <>
            <Form.Item
              name="options"
              label="选项"
              rules={[{ required: true, message: '请输入选项，每行一个' }]}
              extra="每行输入一个选项"
            >
              <TextArea 
                rows={6} 
                placeholder={`输入选项，每行一项，例如：\nA. 选项1\nB. 选项2\nC. 选项3\nD. 选项4`} 
              />
            </Form.Item>
            
            <Form.Item
              name="rightAnswers"
              label="正确答案序号"
              rules={[{ required: true, message: '请输入正确答案序号' }]}
              extra="多个答案用逗号分隔，从0开始，例如：0,2 表示第1个和第3个选项"
            >
              <Input placeholder="例如：0 或 0,2,3" />
            </Form.Item>
          </>
        ) : (
          <Form.Item
            name="code"
            label="代码"
            rules={[{ required: true, message: '请输入代码' }]}
          >
            <TextArea rows={12} placeholder="请输入代码" />
          </Form.Item>
        )}
        
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>
              保存题目
            </Button>
            <Button onClick={() => setManualAdd(false)}>
              取消
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
  
  return (
    <div className="create-question">
      <Card title="题目生成">
        {renderAIGenerateForm()}
        {renderManualAddForm()}
      </Card>
    </div>
  )
}

export default CreateQuestion