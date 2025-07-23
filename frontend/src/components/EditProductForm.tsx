import { useState, useEffect } from 'react';
import api from '@/services/api';
import toast from 'react-hot-toast';
import { Product } from '@/types/Product';
import formStyles from '@/styles/Form.module.css';

interface EditProductFormProps {
    product: Product;
    onSuccess: (updatedProduct: Product) => void;
    onCancel: () => void;
}

export default function EditProductForm({ product, onSuccess, onCancel }: EditProductFormProps) {
    const [name, setName] = useState(product.name);
    const [description, setDescription] = useState(product.description);
    const [price, setPrice] = useState((product.price_in_cents / 100).toFixed(2));
    const [quantity, setQuantity] = useState(product.quantity.toString());
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        setName(product.name);
        setDescription(product.description);
        setPrice((product.price_in_cents / 100).toFixed(2));
        setQuantity(product.quantity.toString());
    }, [product]);

    async function handleSubmit(event: React.FormEvent) {
        event.preventDefault();
        if (!name || !price || !quantity) {
            toast.error('Nome, Preço e Quantidade são obrigatórios.');
            return;
        }
        const payload = { name, description, price_in_cents: Math.round(parseFloat(price) * 100), quantity: parseInt(quantity, 10) };
        setIsLoading(true);
        try {
            const response = await api.put(`/products/${product.id}`, payload);
            toast.success('Produto atualizado com sucesso!');
            onSuccess(response.data);
        } catch (err) {
            toast.error('Ocorreu um erro ao atualizar o produto.');
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form onSubmit={handleSubmit} className={formStyles.form}>
            <label>
                Nome do Produto:
                <input type="text" value={name} onChange={(e) => setName(e.target.value)} className={formStyles.input} />
            </label>
            <label>
                Descrição:
                <textarea value={description} onChange={(e) => setDescription(e.target.value)} className={formStyles.textarea} />
            </label>
            <label>
                Preço (ex: 29.99):
                <input type="number" step="0.01" value={price} onChange={(e) => setPrice(e.target.value)} className={formStyles.input} />
            </label>
            <label>
                Quantidade em Estoque:
                <input type="number" value={quantity} onChange={(e) => setQuantity(e.target.value)} className={formStyles.input} />
            </label>

            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end', marginTop: '1rem' }}>
                <button type="button" onClick={onCancel} className={formStyles.button} style={{backgroundColor: '#6c757d'}}>Cancelar</button>
                <button type="submit" disabled={isLoading} className={formStyles.button}>
                    {isLoading ? 'Salvando...' : 'Salvar Alterações'}
                </button>
            </div>
        </form>
    );
}