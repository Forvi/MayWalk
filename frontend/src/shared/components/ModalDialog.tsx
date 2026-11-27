import { useEffect, useRef } from 'react';
import type { ReactNode } from 'react';

type ModalDialogProps = {
  /** Открыт ли диалог. */
  open: boolean;
  /** Вызывается при закрытии любым путём: Esc, клик по подложке, кнопки. */
  onClose: () => void;
  /** Заголовок диалога (id будет передан в aria-labelledby). */
  titleId: string;
  /** Дополнительное описание для aria-describedby. */
  descriptionId?: string;
  /** Класс нативного <dialog> из дизайн-макета: dialog, vk-dialog, route-dialog… */
  className?: string;
  children: ReactNode;
};

/**
 * Обёртка над нативным <dialog> из дизайн-макета.
 *
 * Поведение перенесено из src/v4/runtime.js дизайнера:
 *  - открытие через showModal() (ловит фокус, блокирует остальную страницу);
 *  - Esc и клик по подложке закрывают диалог (событие cancel перехватывается);
 *  - фокус возвращается на кнопку, которая диалог открыла.
 */
export function ModalDialog({
  open,
  onClose,
  titleId,
  descriptionId,
  className = 'dialog',
  children,
}: ModalDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const returnFocusRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) {
      return;
    }

    if (open && !dialog.open) {
      returnFocusRef.current =
        document.activeElement instanceof HTMLElement ? document.activeElement : null;
      dialog.showModal();
    }
    if (!open && dialog.open) {
      dialog.close();
      returnFocusRef.current?.focus({ preventScroll: true });
    }
  }, [open]);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) {
      return;
    }

    const handleCancel = (event: Event) => {
      event.preventDefault();
      onClose();
    };

    const handleClick = (event: MouseEvent) => {
      // клик по самой подложке (элемент совпадает с <dialog>)
      if (event.target !== dialog) {
        return;
      }
      const rect = dialog.getBoundingClientRect();
      const outside =
        event.clientX < rect.left ||
        event.clientX > rect.right ||
        event.clientY < rect.top ||
        event.clientY > rect.bottom;
      if (outside) {
        onClose();
      }
    };

    dialog.addEventListener('cancel', handleCancel);
    dialog.addEventListener('click', handleClick);
    return () => {
      dialog.removeEventListener('cancel', handleCancel);
      dialog.removeEventListener('click', handleClick);
    };
  }, [onClose]);

  return (
    <dialog
      ref={dialogRef}
      className={className}
      aria-labelledby={titleId}
      aria-describedby={descriptionId}
      data-open={open ? '' : undefined}
    >
      {children}
    </dialog>
  );
}

/** Кнопка «крестик» в шапке диалога — общая для всех модальных окон макета. */
export function DialogCloseButton({ label, onClick }: { label: string; onClick: () => void }) {
  return (
    <button type="button" className="dialog-close" aria-label={label} onClick={onClick}>
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path
          d="m6 6 12 12M18 6 6 18"
          stroke="currentColor"
          strokeWidth="1.8"
          strokeLinecap="round"
        />
      </svg>
    </button>
  );
}
