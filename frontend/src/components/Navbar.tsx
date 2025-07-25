import { Link } from 'react-router-dom';
import styles from '@/styles/components/Navbar.module.css';
import { useAuth } from '@/contexts/AuthContext';

export function Navbar() {
    const { isAuthenticated, user, logout } = useAuth();

    return (
        <header className={styles.navbar}>
            <div className={styles.container}>
                <Link to="/" className={styles.brand}>
                    Controle de Estoque
                </Link>
                <nav className={styles.navLinks}>
                    {isAuthenticated ? (
                        <>
                            <span className={styles.userInfo}>{user?.email}</span>
                            <button onClick={logout} className={styles.logoutButton}>
                                Sair
                            </button>
                        </>
                    ) : (
                        <Link to="/login" className={styles.loginLink}>Login</Link>
                    )}
                </nav>
            </div>
        </header>
    );
}