import { useState } from 'react';
import toast from 'react-hot-toast';
import api from '@/services/api';
import { Client } from '@/types/Client';
import formStyles from '@/styles/Form.module.css';

interface NewClientFormProps {
    onSuccess: (newClient: Client) => void;
    onCancel: () => void;
}

export default function NewClientForm({ onSuccess, onCancel }: NewClientFormProps) {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    async function handleSubmit(event: React.FormEvent) {
        event.preventDefault();
        if (!name) {
            toast.error('O nome do cliente é obrigatório.');
            return;
        }
        setIsLoading(true);
        try {
            const response = await api.post('/clients', { name, email, phone });
            toast.success('Cliente criado com sucesso!');
            onSuccess(response.data);
        } catch (error) {
            toast.error('Ocorreu um erro ao criar o cliente.');
            console.error(error);
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form onSubmit={handleSubmit} className={formStyles.form}>
            <label>Nome do Cliente:<input type="text" value={name} onChange={(e) => setName(e.target.value)} className={formStyles.input} /></label>
            <label>Email:<input type="email" value={email} onChange={(e) => setEmail(e.target.value)} className={formStyles.input} /></label>
            <label>Telefone:<input type="tel" value={phone} onChange={(e) => setPhone(e.target.value)} className={formStyles.input} /></label>
            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end', marginTop: '1rem' }}>
                <button type="button" onClick={onCancel} className={formStyles.button} style={{backgroundColor: '#6c757d'}}>Cancelar</button>
                <button type="submit" disabled={isLoading} className={formStyles.button}>{isLoading ? 'Salvando...' : 'Salvar Cliente'}</button>
            </div>
        </form>
    );
}