import { NavLink } from 'react-router-dom'; // Usando NavLink para estilo ativo
import styles from '@/styles/components/Navbar.module.css';
import { useAuth } from '@/contexts/AuthContext';

export function Navbar() {
    const { isAuthenticated, user, logout } = useAuth();

    return (
        <header className={styles.navbar}>
            <div className={styles.container}>
                <NavLink to="/" className={styles.brand}>
                    Controle de Estoque
                </NavLink>
                <nav className={styles.navLinks}>
                    {isAuthenticated ? (
                        <>
                            {/* Links de Navegação Principal */}
                            <NavLink to="/" className={({ isActive }) => isActive ? `${styles.navLink} ${styles.active}` : styles.navLink}>
                                Produtos
                            </NavLink>
                            <NavLink to="/clients" className={({ isActive }) => isActive ? `${styles.navLink} ${styles.active}` : styles.navLink}>
                                Clientes
                            </NavLink>

                            <div className={styles.separator}></div>

                            {/* Informações do Usuário */}
                            <span className={styles.userInfo}>{user?.email}</span>
                            <button onClick={logout} className={styles.logoutButton}>
                                Sair
                            </button>
                        </>
                    ) : (
                        <NavLink to="/login" className={styles.loginLink}>Login</NavLink>
                    )}
                </nav>
            </div>
        </header>
    );
}