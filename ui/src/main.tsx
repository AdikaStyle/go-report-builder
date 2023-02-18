import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'
import 'primeflex/primeflex.min.css'
import 'primereact/resources/themes/mdc-dark-deeppurple/theme.css'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
