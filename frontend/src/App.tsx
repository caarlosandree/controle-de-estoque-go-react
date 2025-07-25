import { Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { ProtectedRoute } from './components/ProtectedRoute';
import { ProductListPage } from './pages/ProductListPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { ClientListPage } from './pages/ClientListPage';
import { ClientDetailPage } from './pages/ClientDetailPage'; // Importando a nova página

function App() {
    return (
        <Routes>
            {/* Rotas Públicas */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />

            {/* Rotas Privadas (Protegidas) */}
            <Route element={<ProtectedRoute />}>
                <Route path="/" element={<Layout />}>
                    <Route index element={<ProductListPage />} />
                    <Route path="/clients" element={<ClientListPage />} />
                    <Route path="/clients/:clientID" element={<ClientDetailPage />} /> {/* Nova rota de detalhes */}
                </Route>
            </Route>
        </Routes>
    );
}

export default App