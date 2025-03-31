import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import ListLogs from './views/list_logs/list_logs.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ListLogs />
  </StrictMode>,
)
