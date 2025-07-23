import styles from '@/styles/components/Footer.module.css';

export function Footer() {
    const currentYear = new Date().getFullYear();

    return (
        <footer className={styles.footer}>
            <p>&copy; {currentYear} Controle de Estoque. Todos os direitos reservados.</p>
        </footer>
    );
}