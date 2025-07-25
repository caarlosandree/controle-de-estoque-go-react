import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';

export function ProtectedRoute() {
    const { isAuthenticated, isLoading } = useAuth();

    // Enquanto verifica a autenticação, não renderiza nada ou um spinner
    if (isLoading) {
        return <div>Carregando...</div>;
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    // Se autenticado, renderiza a rota filha
    return <Outlet />;
}