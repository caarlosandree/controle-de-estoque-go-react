import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { Pagination } from './Pagination';

describe('Pagination Component', () => {
    it('should render correctly', () => {
        render(<Pagination currentPage={1} totalPages={10} onPageChange={() => {}} />);

        expect(screen.getByRole('button', { name: /anterior/i })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /próxima/i })).toBeInTheDocument();
        expect(screen.getByText('Página 1 de 10')).toBeInTheDocument();
    });

    it('should disable the "previous" button on the first page', () => {
        render(<Pagination currentPage={1} totalPages={10} onPageChange={() => {}} />);

        expect(screen.getByRole('button', { name: /anterior/i })).toBeDisabled();
        expect(screen.getByRole('button', { name: /próxima/i })).not.toBeDisabled();
    });

    it('should disable the "next" button on the last page', () => {
        render(<Pagination currentPage={10} totalPages={10} onPageChange={() => {}} />);

        expect(screen.getByRole('button', { name: /anterior/i })).not.toBeDisabled();
        expect(screen.getByRole('button', { name: /próxima/i })).toBeDisabled();
    });

    it('should call onPageChange with the correct page number when clicking next', async () => {
        const onPageChangeMock = vi.fn();
        render(<Pagination currentPage={5} totalPages={10} onPageChange={onPageChangeMock} />);

        const nextButton = screen.getByRole('button', { name: /próxima/i });
        await userEvent.click(nextButton);

        expect(onPageChangeMock).toHaveBeenCalledTimes(1);  // Corrigido aqui
        expect(onPageChangeMock).toHaveBeenCalledWith(6);
    });

    it('should not render if there is only one page', () => {
        const { container } = render(<Pagination currentPage={1} totalPages={1} onPageChange={() => {}} />);

        expect(container).toBeEmptyDOMElement();
    });
});
