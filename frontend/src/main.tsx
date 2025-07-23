import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'; // Importando o BrowserRouter
import App from '@/App'
import '@/styles/global.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        {/* Envolvemos nosso App com o BrowserRouter para habilitar o roteamento */}
        <BrowserRouter>
            <App />
        </BrowserRouter>
    </React.StrictMode>,
)