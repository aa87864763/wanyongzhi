import { ReactNode } from 'react'

interface SimpleMarkdownProps {
  children: string
}

const SimpleMarkdown = ({ children }: SimpleMarkdownProps): ReactNode => {
  if (!children) return null

  // 将文本分割成行
  const lines = children.split('\n')
  
  // 存储处理后的HTML元素
  const elements: ReactNode[] = []
  
  // 跟踪我们是否在一个列表中
  let inList = false
  let listItems: ReactNode[] = []
  
  // 处理每一行
  lines.forEach((line, index) => {
    const trimmedLine = line.trim()
    
    // 空行结束当前段落或列表
    if (!trimmedLine) {
      if (inList) {
        elements.push(<ul key={`list-${index}`}>{listItems}</ul>)
        listItems = []
        inList = false
      }
      return
    }
    
    // 处理标题
    if (trimmedLine.startsWith('# ')) {
      elements.push(<h1 key={index} dangerouslySetInnerHTML={{ __html: trimmedLine.substring(2) }} />)
    } else if (trimmedLine.startsWith('## ')) {
      elements.push(<h2 key={index} dangerouslySetInnerHTML={{ __html: trimmedLine.substring(3) }} />)
    } else if (trimmedLine.startsWith('### ')) {
      elements.push(<h3 key={index} dangerouslySetInnerHTML={{ __html: trimmedLine.substring(4) }} />)
    } 
    // 处理列表项
    else if (trimmedLine.startsWith('- ') || trimmedLine.startsWith('* ')) {
      inList = true
      listItems.push(<li key={`item-${index}`} dangerouslySetInnerHTML={{ __html: trimmedLine.substring(2) }} />)
    }
    // 处理代码块
    else if (trimmedLine.startsWith('```') && trimmedLine.length > 3) {
      elements.push(<pre key={index}><code>{trimmedLine.substring(3)}</code></pre>)
    }
    // 常规段落
    else {
      if (!trimmedLine.startsWith('```') && !trimmedLine.endsWith('```')) {
        elements.push(<p key={index} dangerouslySetInnerHTML={{ __html: trimmedLine }} />)
      }
    }
  })
  
  // 如果结束时仍在列表中，添加列表
  if (inList && listItems.length > 0) {
    elements.push(<ul key="final-list">{listItems}</ul>)
  }
  
  return <div>{elements}</div>
}

export default SimpleMarkdown 