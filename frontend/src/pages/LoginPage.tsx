import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '@/services/api';
import { useAuth } from '@/contexts/AuthContext';
import formStyles from '@/styles/Form.module.css';
import styles from '@/styles/pages/AuthPages.module.css';

export function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();
    const { login } = useAuth();

    async function handleSubmit(event: React.FormEvent) {
        event.preventDefault();
        setIsLoading(true);
        try {
            const response = await api.post('/login', { email, password });
            await login(response.data.token);
            navigate('/');
            toast.success('Login realizado com sucesso!');
        } catch (error) {
            toast.error('Credenciais inválidas. Tente novamente.');
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className={styles.authContainer}>
            <div className={styles.authBox}>
                <h1>Login</h1>
                <form onSubmit={handleSubmit} className={formStyles.form}>
                    <label>Email:<input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required className={formStyles.input} /></label>
                    <label>Senha:<input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required className={formStyles.input} /></label>
                    <button type="submit" disabled={isLoading} className={formStyles.button}>
                        {isLoading ? 'Entrando...' : 'Entrar'}
                    </button>
                </form>
                <p className={styles.linkText}>
                    Não tem uma conta? <Link to="/register">Registre-se</Link>
                </p>
            </div>
        </div>
    );
}