import { useEffect, useState } from 'react'
import { Routes, Route, useLocation, useNavigate } from 'react-router-dom'
import { Layout, Menu, ConfigProvider, Button } from 'antd'
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons'
import zhCN from 'antd/lib/locale/zh_CN'
import QuestionList from './pages/QuestionList'
import CreateQuestion from './pages/CreateQuestion'
import AIPreview from './pages/AIPreview'
import StudyNotes from './pages/StudyNotes'
import './App.css'

const { Header, Content, Sider } = Layout

function App() {
  const location = useLocation()
  const navigate = useNavigate()
  const [selectedKey, setSelectedKey] = useState('1')
  const [collapsed, setCollapsed] = useState(false)

  useEffect(() => {
    const path = location.pathname
    if (path === '/') {
      setSelectedKey('1')
    } else if (path === '/create' || path === '/question-list' || path.includes('/preview')) {
      setSelectedKey('2')
    }
  }, [location])

  const menuItems = [
    {
      key: '1',
      label: '学习心得',
      onClick: () => navigate('/')
    },
    {
      key: '2',
      label: '题库管理',
      onClick: () => navigate('/question-list')
    }
  ]

  const toggleCollapsed = () => {
    setCollapsed(!collapsed)
  }

  return (
    <ConfigProvider locale={zhCN}>
      <Layout className="app-layout">
        <Sider 
          trigger={null} 
          collapsible 
          collapsed={collapsed}
          width={200}
          className="site-layout-background custom-sider"
          style={{ 
            overflow: 'auto',
            height: '100vh',
            background: '#f0f5fa' 
          }}
        >
          <div className="logo-container">
            <div className="logo">
              <img 
                src="https://volcengine-kdocs-cache.wpscdn.cn/kdocs/img/logo.5c78b00f.svg" 
                alt="金山办公" 
                style={{ 
                  width: '100%', 
                  height: 'auto'
                }} 
              />
            </div>
          </div>
          <Menu
            theme="light"
            mode="inline"
            selectedKeys={[selectedKey]}
            items={menuItems}
            className="side-menu"
          />
          <Button 
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={toggleCollapsed}
            className="collapse-button"
            style={{ color: '#000' }}
          />
        </Sider>
        <Layout>
          <Header className="app-header custom-header">
            <div className="header-title">武汉科技大学 万永智 大作业</div>
          </Header>
          <Content className="app-content">
            <div className="app-container">
              
              <Routes>
                <Route path="/" element={<StudyNotes />} />
                <Route path="/question-list" element={<QuestionList />} />
                <Route path="/create" element={<CreateQuestion />} />
                <Route path="/preview" element={<AIPreview />} />
              </Routes>
            </div>
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  )
}

export default App