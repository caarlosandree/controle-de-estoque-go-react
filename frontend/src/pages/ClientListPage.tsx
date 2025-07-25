import { useState, useEffect, useMemo, lazy, Suspense } from 'react';
import { Link } from 'react-router-dom'; // Importando Link para navegação
import toast from 'react-hot-toast';
import api from '@/services/api';
import { Client } from '@/types/Client';
import styles from '@/styles/pages/ListPage.module.css';
import tableStyles from '@/styles/Table.module.css';
import formStyles from '@/styles/Form.module.css';
import { Modal } from '@/components/Modal';
import { FiEdit, FiTrash2 } from 'react-icons/fi';

const NewClientForm = lazy(() => import('@/components/NewClientForm'));
const EditClientForm = lazy(() => import('@/components/EditClientForm'));

export function ClientListPage() {
    const [clients, setClients] = useState<Client[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [clientToEdit, setClientToEdit] = useState<Client | null>(null);
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [clientToDelete, setClientToDelete] = useState<Client | null>(null);

    useEffect(() => {
        async function fetchClients() {
            setLoading(true);
            try {
                const response = await api.get<Client[]>('/clients');
                setClients(response.data);
            } catch (err) {
                setError('Não foi possível carregar os clientes.');
            } finally {
                setLoading(false);
            }
        }
        fetchClients();
    }, []);

    const filteredClients = useMemo(() => {
        if (!searchTerm) return clients;
        return clients.filter(client =>
            client.name.toLowerCase().includes(searchTerm.toLowerCase())
        );
    }, [clients, searchTerm]);

    function handleCreateSuccess(newClient: Client) {
        setClients(current => [newClient, ...current]);
        setIsCreateModalOpen(false);
    }

    function handleUpdateSuccess(updatedClient: Client) {
        setClients(current => current.map(c => (c.id === updatedClient.id ? updatedClient : c)));
        setIsEditModalOpen(false);
    }

    function handleOpenEditModal(client: Client) {
        setClientToEdit(client);
        setIsEditModalOpen(true);
    }

    function handleOpenDeleteModal(client: Client) {
        setClientToDelete(client);
        setIsDeleteModalOpen(true);
    }

    async function handleDeleteConfirm() {
        if (!clientToDelete) return;
        try {
            await api.delete(`/clients/${clientToDelete.id}`);
            setClients(current => current.filter(c => c.id !== clientToDelete.id));
            toast.success('Cliente deletado com sucesso!');
        } catch (err) {
            toast.error('Erro ao deletar o cliente.');
        } finally {
            setIsDeleteModalOpen(false);
            setClientToDelete(null);
        }
    }

    if (loading) return <div>Carregando clientes...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div>
            <div className={styles.header}>
                <h1>Clientes</h1>
                <button onClick={() => setIsCreateModalOpen(true)} className={formStyles.button}>
                    Novo Cliente
                </button>
            </div>

            <div className={styles.filterContainer}>
                <input
                    type="text"
                    placeholder="Buscar por nome..."
                    className={formStyles.input}
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />
            </div>

            <div className={tableStyles.tableContainer}>
                <table className={tableStyles.table}>
                    <thead>
                    <tr>
                        <th>Nome</th>
                        <th>Email</th>
                        <th>Telefone</th>
                        <th className={tableStyles.actionsCell}>Ações</th>
                    </tr>
                    </thead>
                    <tbody>
                    {filteredClients.length > 0 ? (
                        filteredClients.map(client => (
                            <tr key={client.id}>
                                <td>
                                    <Link to={`/clients/${client.id}`} className={tableStyles.tableLink}>
                                        {client.name}
                                    </Link>
                                </td>
                                <td>{client.email || '-'}</td>
                                <td>{client.phone || '-'}</td>
                                <td className={tableStyles.actionsCell}>
                                    <button
                                        onClick={() => handleOpenEditModal(client)}
                                        className={tableStyles.iconButton}
                                        title="Editar"
                                    >
                                        <FiEdit />
                                    </button>
                                    <button
                                        onClick={() => handleOpenDeleteModal(client)}
                                        className={tableStyles.iconButton}
                                        title="Deletar"
                                    >
                                        <FiTrash2 />
                                    </button>
                                </td>
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan={4} className={tableStyles.emptyState}>
                                Nenhum cliente encontrado.
                            </td>
                        </tr>
                    )}
                    </tbody>
                </table>
            </div>

            {/* Modais */}
            <Modal open={isCreateModalOpen} onOpenChange={setIsCreateModalOpen} title="Criar Novo Cliente">
                <Suspense fallback={<div>Carregando...</div>}>
                    <NewClientForm onSuccess={handleCreateSuccess} onCancel={() => setIsCreateModalOpen(false)} />
                </Suspense>
            </Modal>

            {clientToEdit && (
                <Modal open={isEditModalOpen} onOpenChange={setIsEditModalOpen} title={`Editar ${clientToEdit.name}`}>
                    <Suspense fallback={<div>Carregando...</div>}>
                        <EditClientForm
                            client={clientToEdit}
                            onSuccess={handleUpdateSuccess}
                            onCancel={() => setIsEditModalOpen(false)}
                        />
                    </Suspense>
                </Modal>
            )}

            <Modal open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen} title="Confirmar Exclusão">
                <p>
                    Você tem certeza que deseja deletar o cliente <strong>"{clientToDelete?.name}"</strong>?
                </p>
                <div
                    style={{
                        display: 'flex',
                        gap: '1rem',
                        justifyContent: 'flex-end',
                        marginTop: '1.5rem',
                    }}
                >
                    <button
                        onClick={() => setIsDeleteModalOpen(false)}
                        className={formStyles.button}
                        style={{ backgroundColor: '#6c757d' }}
                    >
                        Cancelar
                    </button>
                    <button onClick={handleDeleteConfirm} className={`${styles.button} ${styles.deleteButton}`}>
                        Sim, deletar
                    </button>
                </div>
            </Modal>
        </div>
    );
}
