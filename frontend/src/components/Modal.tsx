import * as Dialog from '@radix-ui/react-dialog';
import styles from '@/styles/components/Modal.module.css';

// Tipos para as propriedades do nosso modal
interface ModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  children: React.ReactNode; // Conteúdo do modal
}

export function Modal({ open, onOpenChange, title, children }: ModalProps) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        {/* O Overlay é o fundo escurecido */}
        <Dialog.Overlay className={styles.overlay} />
        {/* O Content é a caixa do modal em si */}
        <Dialog.Content className={styles.content}>
          <Dialog.Title className={styles.title}>{title}</Dialog.Title>
          {/* O children permite passar qualquer conteúdo para dentro do modal */}
          <div>{children}</div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
