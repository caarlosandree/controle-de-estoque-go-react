import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import NewProductForm from './NewProductForm';
import { Toaster } from 'react-hot-toast'; // Importando o Toaster

// Mock da API
vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn(),
    },
}));
import api from '@/services/api';

describe('NewProductForm Component', () => {

    afterEach(() => {
        vi.clearAllMocks();
    });

    it('should render all form fields and buttons', () => {
        render(<NewProductForm onSuccess={() => {}} onCancel={() => {}} />);

        expect(screen.getByLabelText(/nome do produto/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/preço/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/quantidade em estoque/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /salvar produto/i })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /cancelar/i })).toBeInTheDocument();
    });

    it('should display validation error if required fields are empty on submit', async () => {
        // Correção: Renderizamos o Toaster junto com o formulário
        render(
            <>
                <Toaster />
                <NewProductForm onSuccess={() => {}} onCancel={() => {}} />
            </>
        );

        const saveButton = screen.getByRole('button', { name: /salvar produto/i });
        await userEvent.click(saveButton);

        expect(await screen.findByText(/nome, preço e quantidade são obrigatórios/i)).toBeInTheDocument();
    });

    it('should call onCancel when cancel button is clicked', async () => {
        const onCancelMock = vi.fn();
        render(<NewProductForm onSuccess={() => {}} onCancel={onCancelMock} />);

        const cancelButton = screen.getByRole('button', { name: /cancelar/i });
        await userEvent.click(cancelButton);

        expect(onCancelMock).toHaveBeenCalledOnce();
    });

    it('should call onSuccess with correct data on successful submission', async () => {
        const onSuccessMock = vi.fn();

        const mockNewProduct = { id: 'uuid-123', name: 'Teclado Mecânico', description: '', price_in_cents: 19990, quantity: 50 };
        (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockNewProduct });

        // Correção: Renderizamos o Toaster também neste teste
        render(
            <>
                <Toaster />
                <NewProductForm onSuccess={onSuccessMock} onCancel={() => {}} />
            </>
        );

        await userEvent.type(screen.getByLabelText(/nome do produto/i), 'Teclado Mecânico');
        await userEvent.type(screen.getByLabelText(/preço/i), '199.90');
        await userEvent.type(screen.getByLabelText(/quantidade em estoque/i), '50');

        const saveButton = screen.getByRole('button', { name: /salvar produto/i });
        await userEvent.click(saveButton);

        await waitFor(() => {
            expect(api.post).toHaveBeenCalledWith('/products', {
                name: 'Teclado Mecânico',
                description: '',
                price_in_cents: 19990,
                quantity: 50,
            });
            expect(onSuccessMock).toHaveBeenCalledWith(mockNewProduct);
        });
    });
});