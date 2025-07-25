import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '@/services/api';
import formStyles from '@/styles/Form.module.css';
import styles from '@/styles/pages/AuthPages.module.css';

export function RegisterPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [passwordConfirm, setPasswordConfirm] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();

    async function handleSubmit(event: React.FormEvent) {
        event.preventDefault();
        setIsLoading(true);
        try {
            await api.post('/register', { email, password, passwordConfirm });
            toast.success('Usuário registrado com sucesso! Faça o login.');
            navigate('/login');
        } catch (error: any) {
            const message = error.response?.data?.message || 'Erro ao registrar. Tente novamente.';
            toast.error(message);
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className={styles.authContainer}>
            <div className={styles.authBox}>
                <h1>Criar Conta</h1>
                <form onSubmit={handleSubmit} className={formStyles.form}>
                    <label>Email:<input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required className={formStyles.input} /></label>
                    <label>Senha:<input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required className={formStyles.input} /></label>
                    <label>Confirmar Senha:<input type="password" value={passwordConfirm} onChange={(e) => setPasswordConfirm(e.target.value)} required className={formStyles.input} /></label>
                    <button type="submit" disabled={isLoading} className={formStyles.button}>
                        {isLoading ? 'Registrando...' : 'Registrar'}
                    </button>
                </form>
                <p className={styles.linkText}>
                    Já tem uma conta? <Link to="/login">Faça o login</Link>
                </p>
            </div>
        </div>
    );
}