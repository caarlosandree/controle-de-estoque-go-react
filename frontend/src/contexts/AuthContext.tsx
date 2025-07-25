import { createContext, useState, useContext, ReactNode, useEffect } from 'react';
import api from '@/services/api';

interface User {
    id: string;
    email: string;
}

interface AuthContextType {
    isAuthenticated: boolean;
    user: User | null;
    login: (token: string) => Promise<void>;
    logout: () => void;
    isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true); // Para checar o token inicial

    useEffect(() => {
        async function loadUserFromStorage() {
            const token = localStorage.getItem('authToken');
            if (token) {
                try {
                    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
                    const response = await api.get('/me');
                    setUser(response.data);
                } catch (error) {
                    console.error("Token inválido, fazendo logout.", error);
                    logout();
                }
            }
            setIsLoading(false);
        }
        loadUserFromStorage();
    }, []);

    const login = async (token: string) => {
        localStorage.setItem('authToken', token);
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        const response = await api.get('/me');
        setUser(response.data);
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem('authToken');
        delete api.defaults.headers.common['Authorization'];
    };

    return (
        <AuthContext.Provider value={{ isAuthenticated: !!user, user, login, logout, isLoading }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}