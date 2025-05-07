import axios from 'axios'
import type { 
  QuestionRequest, 
  QuestionData, 
  QuestionQueryRequest,
  QuestionListResponse,
  HTTPResponse,
  QuestionDeleteRequest
} from '../types'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 50000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 添加响应拦截器，统一处理错误
api.interceptors.response.use(
  response => response,
  error => {
    console.error('API请求错误:', error)
    if (error.response) {
      console.error('错误响应:', error.response.data)
    }
    return Promise.reject(error)
  }
)

// 生成题目
export const createQuestion = async (params: QuestionRequest): Promise<HTTPResponse> => {
  try {
    console.log('生成题目请求参数:', params)
    const response = await api.post<HTTPResponse>('/questions/create', params)
    console.log('生成题目响应:', response.data)
    return response.data
  } catch (error) {
    console.error('生成题目失败:', error)
    throw error
  }
}

// 获取题目列表
export const getQuestionList = async (params: QuestionQueryRequest): Promise<QuestionListResponse> => {
  try {
    console.log('获取题目列表请求参数:', params)
    const response = await api.get<QuestionListResponse>('/questions/list', { params })
    console.log('获取题目列表响应:', response.data)
    return response.data
  } catch (error) {
    console.error('获取题目列表失败:', error)
    throw error
  }
}

// 手动添加题目
export const addQuestion = async (data: QuestionData): Promise<HTTPResponse> => {
  try {
    console.log('添加题目请求数据:', data)
    const response = await api.post<HTTPResponse>('/questions/add', data)
    console.log('添加题目响应:', response.data)
    return response.data
  } catch (error) {
    console.error('添加题目失败:', error)
    throw error
  }
}

// 编辑题目
export const editQuestion = async (id: number, data: QuestionData): Promise<HTTPResponse> => {
  try {
    console.log(`编辑题目(ID:${id})请求数据:`, data)
    const response = await api.put<HTTPResponse>(`/questions/edit/${id}`, data)
    console.log('编辑题目响应:', response.data)
    return response.data
  } catch (error) {
    console.error(`编辑题目(ID:${id})失败:`, error)
    throw error
  }
}

// 删除题目
export const deleteQuestions = async (params: QuestionDeleteRequest): Promise<HTTPResponse> => {
  try {
    console.log('删除题目请求参数:', params)
    // 确保发送的是JSON格式
    const response = await api.delete<HTTPResponse>('/questions/delete', { 
      data: params,
      headers: {
        'Content-Type': 'application/json'
      }
    })
    console.log('删除题目响应:', response.data)
    return response.data
  } catch (error) {
    console.error('删除题目失败:', error)
    throw error
  }
}

export default api 