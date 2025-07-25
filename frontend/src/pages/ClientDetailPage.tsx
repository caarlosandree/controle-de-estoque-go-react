import { useState, useEffect, lazy, Suspense, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '@/services/api';
import { Client } from '@/types/Client';
import { ClientStock } from '@/types/ClientStock';
import styles from '@/styles/pages/DetailPage.module.css';
import tableStyles from '@/styles/Table.module.css';
import formStyles from '@/styles/Form.module.css';
import { Modal } from '@/components/Modal';

const TransferStockForm = lazy(() => import('@/components/TransferStockForm'));

export function ClientDetailPage() {
    const { clientID } = useParams<{ clientID: string }>();
    const [client, setClient] = useState<Client | null>(null);
    const [stock, setStock] = useState<ClientStock[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Estado para controlar o modal de transferência
    const [isTransferModalOpen, setIsTransferModalOpen] = useState(false);

    // `useCallback` para evitar recriar a função a cada renderização
    const fetchClientDetails = useCallback(async () => {
        if (!clientID) return;
        setLoading(true);
        setError(null);
        try {
            const [clientResponse, stockResponse] = await Promise.all([
                api.get<Client>(`/clients/${clientID}`),
                api.get<ClientStock[]>(`/clients/${clientID}/stock`),
            ]);
            setClient(clientResponse.data);
            setStock(stockResponse.data);
        } catch (err) {
            setError('Não foi possível carregar os dados do cliente.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    }, [clientID]);

    useEffect(() => {
        fetchClientDetails();
    }, [fetchClientDetails]);

    // Função chamada pelo formulário em caso de sucesso
    function handleTransferSuccess() {
        setIsTransferModalOpen(false);
        // Recarrega todos os dados da página para refletir as mudanças
        fetchClientDetails();
    }

    if (loading) return <div>Carregando...</div>;
    if (error) return <div>{error}</div>;
    if (!client) return <div>Cliente não encontrado.</div>;

    return (
        <div>
            <div className={styles.header}>
                <Link to="/clients" className={styles.backLink}>&larr; Voltar para Clientes</Link>
                <h1>{client.name}</h1>
                <p className={styles.subHeader}>{client.email || 'Sem email'} | {client.phone || 'Sem telefone'}</p>
            </div>

            <div className={styles.contentBox}>
                <div className={styles.contentHeader}>
                    <h2>Estoque do Cliente</h2>
                    {/* Botão para abrir o modal de transferência */}
                    <button onClick={() => setIsTransferModalOpen(true)} className={formStyles.button}>
                        Transferir Produto
                    </button>
                </div>

                <div className={tableStyles.tableContainer}>
                    <table className={tableStyles.table}>
                        <thead>
                        <tr>
                            <th>Produto</th>
                            <th>Quantidade em Estoque</th>
                            <th className={tableStyles.actionsCell}>Ações</th>
                        </tr>
                        </thead>
                        <tbody>
                        {stock.length > 0 ? (
                            stock.map(item => (
                                <tr key={item.productId}>
                                    <td>{item.productName}</td>
                                    <td>{item.quantity}</td>
                                    <td className={tableStyles.actionsCell}>
                                        {/* Futuramente, botões para ajustar ou devolver estoque */}
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={3} className={tableStyles.emptyState}>
                                    Nenhum produto alocado para este cliente.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Modal de Transferência */}
            <Modal open={isTransferModalOpen} onOpenChange={setIsTransferModalOpen} title="Transferir Estoque">
                <Suspense fallback={<div>Carregando formulário...</div>}>
                    {clientID && (
                        <TransferStockForm
                            clientId={clientID}
                            onSuccess={handleTransferSuccess}
                            onCancel={() => setIsTransferModalOpen(false)}
                        />
                    )}
                </Suspense>
            </Modal>
        </div>
    );
}