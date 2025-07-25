import { expect, afterEach, vi } from 'vitest';
import { cleanup } from '@testing-library/react';
import '@testing-library/jest-dom';

// --- ADICIONE ESTE BLOCO DE CÓDIGO ---
// Mock para a função window.matchMedia que não existe no JSDOM.
// Isso é necessário para bibliotecas como react-hot-toast.
Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(), // Depreciado, mas algumas bibliotecas ainda usam
        removeListener: vi.fn(), // Depreciado, mas algumas bibliotecas ainda usam
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
});
// ------------------------------------

// Limpa o DOM do JSDOM após cada teste
afterEach(() => {
    cleanup();
});