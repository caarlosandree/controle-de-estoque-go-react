import { Link } from 'react-router-dom';
import styles from '@/styles/components/Navbar.module.css';

export function Navbar() {
    return (
        <header className={styles.navbar}>
            <div className={styles.container}>
                <Link to="/" className={styles.brand}>
                    Controle de Estoque
                </Link>
                {/* Futuramente, podemos adicionar outros links de navegação aqui */}
                <nav>
                    {/* Ex: <Link to="/dashboard">Dashboard</Link> */}
                </nav>
            </div>
        </header>
    );
}