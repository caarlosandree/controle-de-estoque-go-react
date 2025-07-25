import { useState, useEffect } from 'react';
import toast from 'react-hot-toast';
import api from '@/services/api';
import { Product } from '@/types/Product'; // Reutilizamos o tipo Product
import { PaginatedResponse } from '@/types/Api'; // Para buscar a lista completa
import formStyles from '@/styles/Form.module.css';

interface TransferStockFormProps {
    clientId: string;
    onSuccess: () => void; // Apenas notifica o sucesso, a página pai recarregará
    onCancel: () => void;
}

export default function TransferStockForm({ clientId, onSuccess, onCancel }: TransferStockFormProps) {
    const [products, setProducts] = useState<Product[]>([]);
    const [selectedProductId, setSelectedProductId] = useState('');
    const [quantity, setQuantity] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    // Busca todos os produtos para preencher o <select>
    useEffect(() => {
        async function fetchAllProducts() {
            try {
                // Pedimos um limite alto para buscar todos os produtos de uma vez
                const response = await api.get<PaginatedResponse<Product>>('/products', {
                    params: { limit: 1000 },
                });
                setProducts(response.data.data);
            } catch (error) {
                toast.error('Não foi possível carregar o catálogo de produtos.');
            }
        }
        fetchAllProducts();
    }, []);

    async function handleSubmit(event: React.FormEvent) {
        event.preventDefault();
        if (!selectedProductId || !quantity || parseInt(quantity, 10) <= 0) {
            toast.error('Por favor, selecione um produto e insira uma quantidade válida.');
            return;
        }

        setIsLoading(true);
        try {
            await api.post(`/products/${selectedProductId}/transfer`, {
                clientId: clientId,
                quantity: parseInt(quantity, 10),
            });
            toast.success('Estoque transferido com sucesso!');
            onSuccess();
        } catch (error: any) {
            const message = error.response?.data?.message || 'Erro ao transferir estoque.';
            toast.error(message);
        } finally {
            setIsLoading(false);
        }
    }

    const selectedProduct = products.find(p => p.id === selectedProductId);

    return (
        <form onSubmit={handleSubmit} className={formStyles.form}>
            <label>
                Selecione um Produto:
                <select
                    value={selectedProductId}
                    onChange={(e) => setSelectedProductId(e.target.value)}
                    className={formStyles.input} // Reutilizando o estilo do input
                    required
                >
                    <option value="" disabled>Selecione...</option>
                    {products.map(product => (
                        <option key={product.id} value={product.id}>
                            {product.name} (Estoque Global: {product.quantity})
                        </option>
                    ))}
                </select>
            </label>

            <label>
                Quantidade a Transferir:
                <input
                    type="number"
                    value={quantity}
                    onChange={(e) => setQuantity(e.target.value)}
                    className={formStyles.input}
                    required
                    min="1"
                    max={selectedProduct?.quantity} // Impede de inserir mais do que o disponível
                />
            </label>

            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end', marginTop: '1rem' }}>
                <button type="button" onClick={onCancel} className={formStyles.button} style={{backgroundColor: '#6c757d'}}>Cancelar</button>
                <button type="submit" disabled={isLoading || !selectedProductId} className={formStyles.button}>
                    {isLoading ? 'Transferindo...' : 'Confirmar Transferência'}
                </button>
            </div>
        </form>
    );
}