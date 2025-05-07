import { useEffect, useState } from 'react'
import { 
  Table, 
  Button, 
  Space, 
  Modal, 
  message, 
  Popconfirm, 
  Form,
  Input,
  Select,
  Card,
  Row,
  Col,
  Radio,
  Menu,
  Dropdown,
  InputNumber,
  Checkbox,
  Spin,
  Alert
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { getQuestionList, deleteQuestions, editQuestion, createQuestion, addQuestion } from '../services/api'
import type { QuestionData, QuestionRequest, AIResponse } from '../types'
import { QuestionType, QuestionDifficulty } from '../types'
import { DownOutlined } from '@ant-design/icons'
import '../styles/QuestionList.css'

const { Option } = Select

const QuestionList = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<QuestionData[]>([])
  const [total, setTotal] = useState(0)
  const [pageParams, setPageParams] = useState({
    page: 1,
    pageSize: 10,
    type: undefined as QuestionType | undefined,
    titleKeyword: undefined as string | undefined
  })
  const [editModalVisible, setEditModalVisible] = useState(false)
  const [aiGenerateModalVisible, setAiGenerateModalVisible] = useState(false)
  const [manualAddModalVisible, setManualAddModalVisible] = useState(false)
  const [currentQuestion, setCurrentQuestion] = useState<QuestionData | null>(null)
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([])
  const [errorMsg, setErrorMsg] = useState<string>('')
  const [currentType, setCurrentType] = useState<QuestionType>(QuestionType.SingleChoice)
  const [aiGenerateForm] = Form.useForm()
  const [manualAddForm] = Form.useForm()
  const [generatedQuestions, setGeneratedQuestions] = useState<AIResponse[]>([])
  const [generatingQuestions, setGeneratingQuestions] = useState(false)
  const [selectedQuestions, setSelectedQuestions] = useState<number[]>([])

  const fetchQuestions = async () => {
    setLoading(true)
    try {
      const res = await getQuestionList({
        page: pageParams.page,
        pageSize: pageParams.pageSize,
        type: pageParams.type,
        title: pageParams.titleKeyword
      })
      if (res.code === 0) {
        setData(res.list)
        setTotal(res.total)
      } else {
        message.error(res.msg || '获取题目列表失败')
      }
    } catch (err) {
      message.error('获取题目列表失败')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchQuestions()
  }, [pageParams])

  const handleDelete = async (ids: number[]) => {
    try {
      if (!ids || ids.length === 0) {
        message.warning('请选择要删除的题目')
        return
      }
      
      const validIds = ids.filter(id => id !== undefined && id !== null)
      
      if (validIds.length === 0) {
        message.warning('无有效的题目ID')
        return
      }
      
      console.log('准备删除题目IDs:', validIds)
      
      const res = await deleteQuestions({ ids: validIds })
      if (res.code === 0) {
        message.success('删除成功')
        fetchQuestions()
        setSelectedRowKeys([])
      } else {
        message.error(res.msg || '删除失败')
      }
    } catch (err) {
      console.error('删除题目时出错:', err)
      message.error('删除失败，请检查网络连接或服务器状态')
    }
  }

  const handleEdit = (record: QuestionData) => {
    setCurrentQuestion(record)
    setCurrentType(record.aiReq.type || QuestionType.SingleChoice)
    
    // 为不同类型的题目设置不同的初始表单值
    const isProgramming = record.aiReq.type === QuestionType.Programming;
    
    const options = !isProgramming ? record.aiRes.answer : [];
    const formValues: any = {
      title: record.aiRes.title,
      type: record.aiReq.type,
      language: record.aiReq.language ? [record.aiReq.language] : [],
      difficulty: record.difficulty || QuestionDifficulty.Medium, 
    };
    
    // 如果是选择题，设置选项
    if (!isProgramming) {
      if (options.length >= 1) formValues.optionA = options[0];
      if (options.length >= 2) formValues.optionB = options[1];
      if (options.length >= 3) formValues.optionC = options[2];
      if (options.length >= 4) formValues.optionD = options[3];
      
      if (record.aiRes.right && record.aiRes.right.length > 0) {
        if (record.aiReq.type === QuestionType.SingleChoice && record.aiRes.right.length > 0) {
          formValues.answers = record.aiRes.right[0];
        } else if (record.aiReq.type === QuestionType.MultiChoice) {
          formValues.answers = Array.isArray(record.aiRes.right) ? record.aiRes.right : [record.aiRes.right[0]];
        }
      }
    }
    
    console.log("设置表单值:", formValues);
    
    // 延迟设置表单值，防止状态更新导致的渲染问题
    setTimeout(() => {
      form.setFieldsValue(formValues);
    }, 100);
    
    setEditModalVisible(true);
  }

  const handleEditSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (!currentQuestion || !currentQuestion.id) {
        message.error('编辑失败：无效的题目')
        return
      }

      const isProgramming = values.type === QuestionType.Programming;
      
      let answer: string[] = [];
      let right: number[] = [];
      
      if (!isProgramming) {
        // 处理选择题选项
        const options = [];
        if (values.optionA) options.push(values.optionA);
        if (values.optionB) options.push(values.optionB);
        if (values.optionC) options.push(values.optionC);
        if (values.optionD) options.push(values.optionD);
        
        answer = options;
        right = Array.isArray(values.answers) ? values.answers : [values.answers];
      }
      
      const language = Array.isArray(values.language) && values.language.length > 0
        ? values.language[0]
        : values.language;

      const updatedQuestion: QuestionData = {
        ...currentQuestion,
        difficulty: values.difficulty,
        aiReq: {
          ...currentQuestion.aiReq,
          type: values.type,
          language: language,
          difficulty: values.difficulty
        },
        aiRes: {
          ...currentQuestion.aiRes,
          title: values.title,
          answer: answer,
          right: right,
          code: ''
        }
      }

      const res = await editQuestion(currentQuestion.id, updatedQuestion)
      if (res.code === 0) {
        message.success('编辑成功')
        setEditModalVisible(false)
        fetchQuestions()
      } else {
        message.error(res.msg || '编辑失败')
      }
    } catch (err) {
      message.error('编辑失败')
      console.error(err)
    }
  }

  const handlePageChange = (page: number, pageSize?: number) => {
    setPageParams(prev => ({
      ...prev,
      page,
      pageSize: pageSize || prev.pageSize
    }))
  }

  const handleTitleSearch = (value: string) => {
    setPageParams(prev => ({
      ...prev,
      page: 1,
      titleKeyword: value
    }))
  }

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的题目')
      return
    }
    
    const ids = selectedRowKeys.map(key => Number(key))
    await handleDelete(ids)
  }

  const getQuestionTypeText = (type: QuestionType) => {
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

  const columns: ColumnsType<QuestionData> = [
    {
      title: '题目标题',
      dataIndex: ['aiRes', 'title'],
      key: 'title',
      ellipsis: true
    },
    {
      title: '题目类型',
      dataIndex: ['aiReq', 'type'],
      key: 'type',
      width: 120,
      render: (type: QuestionType) => getQuestionTypeText(type)
    },
    {
      title: '题目难度',
      dataIndex: 'difficulty',
      key: 'difficulty',
      width: 100,
      render: (difficulty: QuestionDifficulty) => {
        switch (difficulty) {
          case QuestionDifficulty.Easy:
            return <span style={{ color: 'green' }}>简单</span>
          case QuestionDifficulty.Medium:
            return <span style={{ color: 'orange' }}>中等</span>
          case QuestionDifficulty.Hard:
            return <span style={{ color: 'red' }}>困难</span>
          default:
            return <span style={{ color: 'orange' }}>中等</span>
        }
      }
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space size="small">
          <Button type="link" onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个题目吗？"
            onConfirm={() => handleDelete([record.id!])}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys)
  }

  const handleAIGenerate = async (values: QuestionRequest) => {
    setGeneratingQuestions(true)
    setErrorMsg('')
    
    try {
      const requestParams = {
        ...values,
        // 取语言数组的第一个值作为请求参数
        language: Array.isArray(values.language) && values.language.length > 0 
          ? values.language[0].toLowerCase() 
          : (values.language || "go"),
        difficulty: values.difficulty
      }
      
      console.log('准备生成题目，参数:', requestParams)
      const res = await createQuestion(requestParams)
      
      if (res.code === 0) {
        message.success('生成题目成功')
        
        let aiResData: AIResponse[] = [];
        
        // 处理不同的响应格式
        if (res.aiRes) {
          aiResData = Array.isArray(res.aiRes) ? res.aiRes : [res.aiRes];
        } else if (res.questions && res.questions.length > 0) {
          aiResData = res.questions.map(q => q.aiRes);
        }
        
        if (aiResData.length === 0) {
          setErrorMsg('未找到生成的题目数据')
          message.error('未找到生成的题目数据')
          return
        }
        
        // 设置生成的题目，并按题目数量创建唯一索引
        setGeneratedQuestions(aiResData)
        setSelectedQuestions(Array.from(new Set(aiResData.map((_, index) => index))))
      } else {
        setErrorMsg(res.msg || '生成题目失败，请稍后重试')
        message.error(res.msg || '生成题目失败')
      }
    } catch (err) {
      console.error('生成题目错误:', err)
      setErrorMsg('生成题目失败，请检查网络连接或服务器状态')
      message.error('生成题目失败，请稍后重试')
    } finally {
      setGeneratingQuestions(false)
    }
  }

  const handleSaveGeneratedQuestions = async () => {
    setLoading(true)
    try {
      const aiGenerateValues = await aiGenerateForm.validateFields()
      
      // 取第一个语言作为保存参数
      const language = Array.isArray(aiGenerateValues.language) && aiGenerateValues.language.length > 0
        ? aiGenerateValues.language[0]
        : (aiGenerateValues.language || "go")
      
      const selectedQuestionsData = selectedQuestions.map(index => {
        const aiRes = generatedQuestions[index]
        const currentTime = new Date().toISOString()
        return {
          aiStartTime: currentTime,
          aiEndTime: currentTime,
          aiCostTime: 0,
          aiReq: {
            type: aiGenerateValues.type,
            difficulty: aiGenerateValues.difficulty,
            language: language
          },
          aiRes: aiRes,
          difficulty: aiGenerateValues.difficulty,
          createdAt: currentTime
        }
      })

      let successCount = 0
      for (const questionData of selectedQuestionsData) {
        try {
          const res = await addQuestion(questionData)
          if (res.code === 0) {
            successCount++
          }
        } catch (error) {
          console.error('保存题目失败:', error)
        }
      }

      if (successCount > 0) {
        message.success(`成功保存 ${successCount} 道题目`)
        setAiGenerateModalVisible(false)
        setGeneratedQuestions([])
        setSelectedQuestions([])
        aiGenerateForm.resetFields()
        fetchQuestions()
      } else {
        message.error('所有题目保存失败')
      }
    } catch (error) {
      message.error('保存题目失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleManualAdd = async (values: any) => {
    setLoading(true)
    setErrorMsg('')
    
    try {
      const isProgramming = values.type === QuestionType.Programming;
      
      let answer: string[] = [];
      let right: number[] = [];
      
      if (!isProgramming) {
        const options = [];
        if (values.optionA) options.push(values.optionA);
        if (values.optionB) options.push(values.optionB);
        if (values.optionC) options.push(values.optionC);
        if (values.optionD) options.push(values.optionD);
        
        answer = options;
        right = Array.isArray(values.answers) ? values.answers : [values.answers];
        
        if (answer.length < 2) {
          setErrorMsg('至少需要输入两个选项')
          message.error('至少需要输入两个选项')
          setLoading(false)
          return
        }
        
        if (right.length === 0) {
          setErrorMsg('请选择至少一个正确答案')
          message.error('请选择至少一个正确答案')
          setLoading(false)
          return
        }
        
        if (right.some(idx => idx < 0 || idx >= answer.length)) {
          setErrorMsg('答案选择无效')
          message.error('答案选择无效')
          setLoading(false)
          return
        }
      }

      const language = Array.isArray(values.language) && values.language.length > 0
        ? values.language[0]
        : values.language;
      
      // 构建题目数据
      const currentTime = new Date().toISOString()
      const questionData: QuestionData = {
        aiStartTime: currentTime,
        aiEndTime: currentTime,
        aiCostTime: 0,
        aiReq: {
          type: values.type,
          difficulty: values.difficulty,
          language: language
        },
        aiRes: {
          title: values.title,
          answer: answer,
          right: right,
          code: ''
        },
        difficulty: values.difficulty,
        createdAt: currentTime
      }
      
      console.log('准备手动添加题目:', questionData)
      
      const res = await addQuestion(questionData)
      if (res.code === 0) {
        message.success('添加题目成功')
        manualAddForm.resetFields()
        setManualAddModalVisible(false)
        fetchQuestions()
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

  const handleQuestionSelect = (index: number) => {
    setSelectedQuestions(prev => {
      const isSelected = prev.includes(index)
      if (isSelected) {
        return prev.filter(i => i !== index)
      } else {
        return [...prev, index]
      }
    })
  }

  return (
    <div className="question-list">
 <Card className="filter-card">
  <Row gutter={16} align="middle" justify="space-between">
    <Col xs={24} sm={12} md={8} lg={8} xl={8}>
      <div className="filter-item">
        <span className="filter-label">选择类型：</span>
        <Radio.Group 
          options={[
            { label: '全部', value: undefined },
            { label: '单选题', value: QuestionType.SingleChoice },
            { label: '多选题', value: QuestionType.MultiChoice },
            { label: '编程题', value: QuestionType.Programming }
          ]}
          onChange={(e) => setPageParams(prev => ({
            ...prev,
            page: 1,
            type: e.target.value
          }))}
          value={pageParams.type}
          optionType="button"
          size="small"
        />
      </div>
    </Col>
    <Col xs={24} sm={12} md={16} lg={16} xl={16} style={{ display: 'flex', justifyContent: 'flex-end' }}>
      <div className="filter-item">
        <span className="filter-label">题目搜索：</span>
        <Input.Search
          placeholder="输入题目关键词"
          allowClear
          onSearch={handleTitleSearch}
          style={{ width: '300px', maxWidth: '100%' }}
        />
      </div>
    </Col>
  </Row>
  <Row justify="end" style={{ marginTop: 16 }}>
    <Col>
      <Space>
        <Dropdown
          overlay={
            <Menu>
              <Menu.Item key="ai" onClick={() => setAiGenerateModalVisible(true)}>
                AI出题
              </Menu.Item>
              <Menu.Item key="manual" onClick={() => setManualAddModalVisible(true)}>
                人工出题
              </Menu.Item>
            </Menu>
          }
          trigger={['click']}
        >
          <Button type="primary">
            出题 <DownOutlined />
          </Button>
        </Dropdown>
        <Popconfirm
          title="确定要删除选中的题目吗？"
          onConfirm={handleBatchDelete}
          okText="确定"
          cancelText="取消"
          disabled={selectedRowKeys.length === 0}
        >
          <Button danger disabled={selectedRowKeys.length === 0}>
            批量删除
          </Button>
        </Popconfirm>
      </Space>
    </Col>
  </Row>
</Card>

      <Table
        rowSelection={rowSelection}
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          current: pageParams.page,
          pageSize: pageParams.pageSize,
          total: total,
          onChange: handlePageChange,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`
        }}
        className="questions-table"
      />

      <Modal
        title={
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <span>AI生成试题</span>
            <img 
              src="https://platform.wps.cn/assets/logo_ai-80593444.svg" 
              alt="金山AI" 
              style={{ height: '24px', marginLeft: '8px' }} 
            />
          </div>
        }
        open={aiGenerateModalVisible}
        onCancel={() => {
          setAiGenerateModalVisible(false)
          setGeneratedQuestions([])
          setSelectedQuestions([])
          aiGenerateForm.resetFields()
        }}
        width={1000}
        footer={[
          <Button key="cancel" onClick={() => {
            setAiGenerateModalVisible(false)
            setGeneratedQuestions([])
            setSelectedQuestions([])
            aiGenerateForm.resetFields()
          }}>
            取消
          </Button>,
          <Button 
            key="save" 
            type="primary" 
            disabled={generatedQuestions.length === 0}
            onClick={handleSaveGeneratedQuestions}
            loading={loading}
          >
            保存到题库
          </Button>
        ]}
      >
        <Row gutter={24}>
          <Col span={8}>
            <Form 
              form={aiGenerateForm} 
              layout="vertical"
              initialValues={{
                type: QuestionType.SingleChoice,
                count: 1,
                language: []
              }}
            >
              <Form.Item
                name="type"
                label="题目类型"
                rules={[{ required: true, message: '请选择题目类型' }]}
              >
                <Radio.Group onChange={e => setCurrentType(e.target.value)}>
                  <Radio value={QuestionType.SingleChoice}>单选题</Radio>
                  <Radio value={QuestionType.MultiChoice}>多选题</Radio>
                  <Radio value={QuestionType.Programming}>编程题</Radio>
                </Radio.Group>
              </Form.Item>

              <Form.Item
                name="difficulty"
                label="题目难度"
                rules={[{ required: true, message: '请选择题目难度' }]}
              >
                <Radio.Group>
                  <Radio value={QuestionDifficulty.Easy}>简单</Radio>
                  <Radio value={QuestionDifficulty.Medium}>中等</Radio>
                  <Radio value={QuestionDifficulty.Hard}>困难</Radio>
                </Radio.Group>
              </Form.Item>

              <Form.Item
                name="count"
                label="生成数量"
                rules={[{ required: true, message: '请输入生成数量' }]}
              >
                <InputNumber min={1} max={10} />
              </Form.Item>

              <Form.Item
                name="language"
                label="编程语言"
                rules={[
                  { 
                    required: true, 
                    message: '请选择编程语言' 
                  }
                ]}
              >
                <Select mode="multiple">
                  <Option value="javascript">JavaScript</Option>
                  <Option value="python">Python</Option>
                  <Option value="java">Java</Option>
                  <Option value="c++">C++</Option>
                  <Option value="go">Go</Option>
                </Select>
              </Form.Item>

              <Form.Item>
                <Button 
                  type="primary" 
                  onClick={() => aiGenerateForm.validateFields().then(handleAIGenerate)}
                  loading={generatingQuestions}
                >
                  生成题目
                </Button>
              </Form.Item>

              {errorMsg && <Alert message={errorMsg} type="error" showIcon />}
            </Form>
          </Col>

          <Col span={16} style={{ borderLeft: '1px solid #f0f0f0', paddingLeft: '20px' }}>
            <div className="ai-generated-questions">
              <h3>生成结果</h3>
              {generatingQuestions ? (
                <div style={{ textAlign: 'center', padding: '30px' }}>
                  <Spin tip="正在生成题目..." />
                </div>
              ) : generatedQuestions.length > 0 ? (
                <div>
                  {generatedQuestions.map((question, index) => (
                    <Card 
                      key={index} 
                      style={{ marginBottom: '16px' }}
                      title={
                        <Checkbox 
                          checked={selectedQuestions.includes(index)}
                          onChange={() => handleQuestionSelect(index)}
                        >
                          题目 {index + 1}
                        </Checkbox>
                      }
                    >
                      <h4>{question.title}</h4>
                      {currentType !== QuestionType.Programming ? (
                        <div>
                          <p>选项：</p>
                          <ul>
                            {question.answer.map((option, optIndex) => (
                              <li key={optIndex} style={{ color: question.right.includes(optIndex) ? 'green' : 'inherit' }}>
                                {option} {question.right.includes(optIndex) && '(正确)'}
                              </li>
                            ))}
                          </ul>
                        </div>
                      ) : null}
                    </Card>
                  ))}
                </div>
              ) : (
                <div style={{ textAlign: 'center', padding: '30px', color: '#999' }}>
                  点击"生成题目"按钮开始生成
                </div>
              )}
            </div>
          </Col>
        </Row>
      </Modal>

      <Modal
        title="编辑题目"
        open={editModalVisible}
        onOk={handleEditSubmit}
        onCancel={() => setEditModalVisible(false)}
        width={700}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="type"
            label="题目类型"
            rules={[{ required: true, message: '请选择题目类型' }]}
          >
            <Radio.Group 
              onChange={e => {
                const newType = e.target.value;
                setCurrentType(newType);
                
                // 切换题目类型时，重置部分表单字段
                if (newType === QuestionType.Programming) {
                  form.setFieldsValue({
                    answers: undefined, 
                  });
                } else if (newType === QuestionType.SingleChoice) {
                  // 如果有多选题答案，取第一个作为单选题答案
                  const multiAnswers = form.getFieldValue('answers');
                  if (Array.isArray(multiAnswers) && multiAnswers.length > 0) {
                    form.setFieldsValue({
                      answers: multiAnswers[0]
                    });
                  }
                } else if (newType === QuestionType.MultiChoice) {
                  const singleAnswer = form.getFieldValue('answers');
                  if (!Array.isArray(singleAnswer) && singleAnswer !== undefined) {
                    form.setFieldsValue({
                      answers: [singleAnswer]
                    });
                  }
                }
              }}
            >
              <Radio value={QuestionType.SingleChoice}>单选题</Radio>
              <Radio value={QuestionType.MultiChoice}>多选题</Radio>
              <Radio value={QuestionType.Programming}>编程题</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            name="difficulty"
            label="题目难度"
            rules={[{ required: true, message: '请选择题目难度' }]}
          >
            <Radio.Group>
              <Radio value={QuestionDifficulty.Easy}>简单</Radio>
              <Radio value={QuestionDifficulty.Medium}>中等</Radio>
              <Radio value={QuestionDifficulty.Hard}>困难</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            name="language"
            label="编程语言"
            rules={[
              { 
                required: true, 
                message: '请选择编程语言' 
              }
            ]}
          >
            <Select mode="multiple">
              <Option value="javascript">JavaScript</Option>
              <Option value="python">Python</Option>
              <Option value="java">Java</Option>
              <Option value="c++">C++</Option>
              <Option value="go">Go</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="title"
            label="题目标题"
            rules={[{ required: true, message: '请输入题目标题' }]}
          >
            <Input.TextArea autoSize={{ minRows: 2, maxRows: 6 }} />
          </Form.Item>
          
          {currentType !== QuestionType.Programming ? (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="optionA"
                    label="A:"
                    rules={[{ required: true, message: '请输入选项A' }]}
                  >
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="optionB"
                    label="B:"
                    rules={[{ required: true, message: '请输入选项B' }]}
                  >
                    <Input />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="optionC"
                    label="C:"
                    rules={[{ required: true, message: '请输入选项C' }]}
                  >
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="optionD"
                    label="D:"
                    rules={[{ required: true, message: '请输入选项D' }]}
                  >
                    <Input />
                  </Form.Item>
                </Col>
              </Row>
              
              <Form.Item
                name="answers"
                label="正确答案"
                rules={[{ required: true, message: '请选择正确答案' }]}
              >
                {currentType === QuestionType.SingleChoice ? (
                  <Radio.Group>
                    <Radio value={0}>A</Radio>
                    <Radio value={1}>B</Radio>
                    <Radio value={2}>C</Radio>
                    <Radio value={3}>D</Radio>
                  </Radio.Group>
                ) : (
                  <Checkbox.Group>
                    <Checkbox value={0}>A</Checkbox>
                    <Checkbox value={1}>B</Checkbox>
                    <Checkbox value={2}>C</Checkbox>
                    <Checkbox value={3}>D</Checkbox>
                  </Checkbox.Group>
                )}
              </Form.Item>
            </>
          ) : null}
        </Form>
      </Modal>

      <Modal
        title="人工出题"
        open={manualAddModalVisible}
        onOk={() => manualAddForm.validateFields().then(handleManualAdd)}
        onCancel={() => {
          setManualAddModalVisible(false)
          manualAddForm.resetFields()
        }}
        width={700}
        okText="保存"
        cancelText="取消"
        confirmLoading={loading}
      >
        <Form form={manualAddForm} layout="vertical" initialValues={{ type: QuestionType.SingleChoice }}>
          <Form.Item
            name="type"
            label="题目类型"
            rules={[{ required: true, message: '请选择题目类型' }]}
          >
            <Radio.Group 
              onChange={e => {
                const newType = e.target.value;
                setCurrentType(newType);
                
                // 切换题目类型时，重置相关表单字段
                if (newType === QuestionType.Programming) {
                  manualAddForm.setFieldsValue({
                    answers: undefined 
                  });
                } else if (newType === QuestionType.SingleChoice) {
                  // 如果有多选题答案，取第一个作为单选题答案
                  const multiAnswers = manualAddForm.getFieldValue('answers');
                  if (Array.isArray(multiAnswers) && multiAnswers.length > 0) {
                    manualAddForm.setFieldsValue({
                      answers: multiAnswers[0]
                    });
                  }
                } else if (newType === QuestionType.MultiChoice) {
                  // 如果有单选题答案，转为数组
                  const singleAnswer = manualAddForm.getFieldValue('answers');
                  if (!Array.isArray(singleAnswer) && singleAnswer !== undefined) {
                    manualAddForm.setFieldsValue({
                      answers: [singleAnswer]
                    });
                  }
                }
              }}
            >
              <Radio value={QuestionType.SingleChoice}>单选题</Radio>
              <Radio value={QuestionType.MultiChoice}>多选题</Radio>
              <Radio value={QuestionType.Programming}>编程题</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            name="difficulty"
            label="题目难度"
            rules={[{ required: true, message: '请选择题目难度' }]}
          >
            <Radio.Group>
              <Radio value={QuestionDifficulty.Easy}>简单</Radio>
              <Radio value={QuestionDifficulty.Medium}>中等</Radio>
              <Radio value={QuestionDifficulty.Hard}>困难</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            name="language"
            label="编程语言"
            rules={[
              { 
                required: true, 
                message: '请选择编程语言' 
              }
            ]}
          >
            <Select mode="multiple">
              <Option value="javascript">JavaScript</Option>
              <Option value="python">Python</Option>
              <Option value="java">Java</Option>
              <Option value="c++">C++</Option>
              <Option value="go">Go</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="title"
            label="题目标题"
            rules={[{ required: true, message: '请输入题目标题' }]}
          >
            <Input.TextArea autoSize={{ minRows: 2, maxRows: 6 }} />
          </Form.Item>
          
          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) => prevValues.type !== currentValues.type}
          >
            {({ getFieldValue }) => {
              const type = getFieldValue('type');
              return type !== QuestionType.Programming ? (
                <>
                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item
                        name="optionA"
                        label="A:"
                        rules={[{ required: true, message: '请输入选项A' }]}
                      >
                        <Input />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item
                        name="optionB"
                        label="B:"
                        rules={[{ required: true, message: '请输入选项B' }]}
                      >
                        <Input />
                      </Form.Item>
                    </Col>
                  </Row>
                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item
                        name="optionC"
                        label="C:"
                        rules={[{ required: true, message: '请输入选项C' }]}
                      >
                        <Input />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item
                        name="optionD"
                        label="D:"
                        rules={[{ required: true, message: '请输入选项D' }]}
                      >
                        <Input />
                      </Form.Item>
                    </Col>
                  </Row>
                  
                  <Form.Item
                    name="answers"
                    label="正确答案"
                    rules={[{ required: true, message: '请选择正确答案' }]}
                  >
                    {type === QuestionType.SingleChoice ? (
                      <Radio.Group>
                        <Radio value={0}>A</Radio>
                        <Radio value={1}>B</Radio>
                        <Radio value={2}>C</Radio>
                        <Radio value={3}>D</Radio>
                      </Radio.Group>
                    ) : (
                      <Checkbox.Group>
                        <Checkbox value={0}>A</Checkbox>
                        <Checkbox value={1}>B</Checkbox>
                        <Checkbox value={2}>C</Checkbox>
                        <Checkbox value={3}>D</Checkbox>
                      </Checkbox.Group>
                    )}
                  </Form.Item>
                </>
              ) : null
            }}
          </Form.Item>
          
          {errorMsg && <Alert message={errorMsg} type="error" showIcon style={{ marginBottom: '16px' }} />}
        </Form>
      </Modal>
    </div>
  )
}

export default QuestionList