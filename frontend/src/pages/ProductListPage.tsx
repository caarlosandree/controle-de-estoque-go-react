import {useState, useEffect, useCallback, lazy} from 'react';
import api from '@/services/api';
import toast from 'react-hot-toast';
import { Product } from '@/types/Product';
import { Metadata, PaginatedResponse } from '@/types/Api';
import styles from '@/styles/pages/ProductListPage.module.css';
import tableStyles from '@/styles/Table.module.css';
import formStyles from '@/styles/Form.module.css';
import { Modal } from '@/components/Modal';
import { Pagination } from '@/components/Pagination';
import { FiEdit, FiTrash2 } from 'react-icons/fi';

const EditProductForm = lazy(() => import('@/components/EditProductForm'));
const NewProductForm = lazy(() => import('@/components/NewProductForm'));

export function ProductListPage() {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [debouncedSearchTerm, setDebouncedSearchTerm] = useState('');
    const [metadata, setMetadata] = useState<Metadata | null>(null);
    const [currentPage, setCurrentPage] = useState(1);
    const PAGE_LIMIT = 10;
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [productToDelete, setProductToDelete] = useState<Product | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [productToEdit, setProductToEdit] = useState<Product | null>(null);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

    useEffect(() => {
        const timerId = setTimeout(() => { setDebouncedSearchTerm(searchTerm); setCurrentPage(1); }, 500);
        return () => { clearTimeout(timerId); };
    }, [searchTerm]);

    const fetchProducts = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await api.get<PaginatedResponse<Product>>('/products', { params: { page: currentPage, limit: PAGE_LIMIT, search: debouncedSearchTerm } });
            setProducts(response.data.data);
            setMetadata(response.data.metadata);
        } catch (err) {
            setError('Não foi possível carregar os produtos.');
        } finally {
            setLoading(false);
        }
    }, [currentPage, debouncedSearchTerm]);

    useEffect(() => { fetchProducts(); }, [fetchProducts]);

    function handleOpenDeleteModal(product: Product) { setProductToDelete(product); setIsDeleteModalOpen(true); }

    async function handleDeleteConfirm() {
        if (!productToDelete) return;
        try {
            await api.delete(`/products/${productToDelete.id}`);
            toast.success('Produto deletado com sucesso!'); // Toast de sucesso
            fetchProducts();
        } catch (err) {
            toast.error('Erro ao deletar o produto.'); // Toast de erro
        } finally {
            setIsDeleteModalOpen(false);
            setProductToDelete(null);
        }
    }

    function handleOpenEditModal(product: Product) { setProductToEdit(product); setIsEditModalOpen(true); }
    function handleUpdateSuccess(updatedProduct: Product) { setProducts(p => p.map(prod => prod.id === updatedProduct.id ? updatedProduct : prod)); setIsEditModalOpen(false); }
    function handleCreateSuccess() { setCurrentPage(1); fetchProducts(); setIsCreateModalOpen(false); }

    const modalActionsStyle: React.CSSProperties = { display: 'flex', gap: '1rem', justifyContent: 'flex-end', marginTop: '1.5rem' };

    return (
        <div>
            <div className={styles.header}>
                <h1>Produtos</h1>
                <button onClick={() => setIsCreateModalOpen(true)} className={`${formStyles.button}`}>
                    Novo Produto
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

            {loading && <div>Carregando...</div>}
            {!loading && error && <div>Erro: {error}</div>}

            {!loading && !error && (
                <>
                    <div className={tableStyles.tableContainer}>
                        <table className={tableStyles.table}>
                            <thead>
                            <tr>
                                <th>Nome</th>
                                <th>Descrição</th>
                                <th>Preço</th>
                                <th>Quantidade</th>
                                <th className={tableStyles.actionsCell}>Ações</th>
                            </tr>
                            </thead>
                            <tbody>
                            {products.length > 0 ? (
                                products.map((product) => (
                                    <tr key={product.id}>
                                        <td>{product.name}</td>
                                        <td>{product.description || '-'}</td>
                                        <td>R$ {(product.price_in_cents / 100).toFixed(2)}</td>
                                        <td>{product.quantity}</td>
                                        <td className={tableStyles.actionsCell}>
                                            <button onClick={() => handleOpenEditModal(product)} className={tableStyles.iconButton} title="Editar">
                                                <FiEdit />
                                            </button>
                                            <button onClick={() => handleOpenDeleteModal(product)} className={tableStyles.iconButton} title="Deletar">
                                                <FiTrash2 />
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan={5} className={tableStyles.emptyState}>
                                        Nenhum produto encontrado.
                                    </td>
                                </tr>
                            )}
                            </tbody>
                        </table>
                    </div>

                    {metadata && <Pagination currentPage={metadata.current_page} totalPages={metadata.total_pages} onPageChange={setCurrentPage} />}
                </>
            )}

            {/* Modais */}
            <Modal open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen} title="Confirmar Exclusão">
                <p>Você tem certeza que deseja deletar o produto <strong>"{productToDelete?.name}"</strong>?</p>
                <div style={modalActionsStyle}>
                    <button onClick={() => setIsDeleteModalOpen(false)} className={formStyles.button} style={{backgroundColor: '#6c757d'}}>Cancelar</button>
                    <button onClick={handleDeleteConfirm} className={`${formStyles.button} ${styles.deleteButton}`}>Sim, deletar</button>
                </div>
            </Modal>

            {productToEdit && (
                <Modal open={isEditModalOpen} onOpenChange={setIsEditModalOpen} title={`Editar ${productToEdit.name}`}>
                    <EditProductForm product={productToEdit} onSuccess={handleUpdateSuccess} onCancel={() => setIsEditModalOpen(false)} />
                </Modal>
            )}

            <Modal open={isCreateModalOpen} onOpenChange={setIsCreateModalOpen} title="Criar Novo Produto">
                <NewProductForm onSuccess={handleCreateSuccess} onCancel={() => setIsCreateModalOpen(false)} />
            </Modal>
        </div>
    );
}