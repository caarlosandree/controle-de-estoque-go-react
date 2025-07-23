import styles from '@/styles/components/Pagination.module.css';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
    // Não renderiza nada se houver apenas uma página ou nenhuma
    if (totalPages <= 1) {
        return null;
    }

    return (
        <div className={styles.paginationContainer}>
            <button
                onClick={() => onPageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className={styles.button}
            >
                Anterior
            </button>

            <span className={styles.pageInfo}>
        Página {currentPage} de {totalPages}
      </span>

            <button
                onClick={() => onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className={styles.button}
            >
                Próxima
            </button>
        </div>
    );
}