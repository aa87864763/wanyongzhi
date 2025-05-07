// 题目类型
export enum QuestionType {
  SingleChoice = 1,
  MultiChoice = 2,
  Programming = 3
}

// 题目难度
export enum QuestionDifficulty {
  Easy = 1,
  Medium = 2,
  Hard = 3
}

// 模型提供商
export type ModelProvider = 'deepseek' | 'tongyi'

// 编程语言
export type ProgrammingLanguage = 'go' | 'java' | 'python' | 'c++' | 'javascript'

// 题目请求
export interface QuestionRequest {
  model?: ModelProvider
  language?: ProgrammingLanguage
  type?: QuestionType
  difficulty?: QuestionDifficulty
  count?: number
}

// AI生成的题目响应
export interface AIResponse {
  title: string
  answer: string[]
  right: number[]
  code?: string
}

// 完整的题目数据
export interface QuestionData {
  id?: number
  aiStartTime: string
  aiEndTime: string
  aiCostTime: number
  aiReq: QuestionRequest
  aiRes: AIResponse
  difficulty: QuestionDifficulty
  createdAt: string
}

// HTTP响应
export interface HTTPResponse {
  code: number
  msg: string
  aiRes?: AIResponse | AIResponse[]
  questions?: QuestionData[]
  totalCount?: number
}

// 题目查询请求
export interface QuestionQueryRequest {
  page: number
  pageSize: number
  type?: QuestionType
  title?: string
}

// 题目列表响应
export interface QuestionListResponse {
  code: number
  msg: string
  total: number
  list: QuestionData[]
}

// 题目删除请求
export interface QuestionDeleteRequest {
  ids: number[]
} 