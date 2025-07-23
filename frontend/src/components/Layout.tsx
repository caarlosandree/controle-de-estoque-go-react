import { Outlet } from 'react-router-dom';
import { Navbar } from './Navbar';
import { Footer } from './Footer';
import styles from '@/styles/components/Layout.module.css';
import { Toaster } from 'react-hot-toast'; // Importando o Toaster

export function Layout() {
    return (
        <div className={styles.layout}>
            {/* O Toaster renderiza as notificações.
        Configuramos para aparecer no canto superior direito
        e com um estilo que combina com nosso tema.
      */}
            <Toaster
                position="top-right"
                toastOptions={{
                    style: {
                        background: 'var(--color-surface)',
                        color: 'var(--color-text-primary)',
                        border: '1px solid var(--color-border)',
                    },
                    success: {
                        iconTheme: {
                            primary: 'var(--color-success)',
                            secondary: 'white',
                        },
                    },
                    error: {
                        iconTheme: {
                            primary: 'var(--color-danger)',
                            secondary: 'white',
                        },
                    },
                }}
            />

            <Navbar />
            <main className={styles.mainContent}>
                <Outlet />
            </main>
            <Footer />
        </div>
    )
}