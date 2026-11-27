import { useState } from 'react';
import type { ReactNode, Ref } from 'react';

type PasswordFieldProps = {
  /** Идентификатор поля — печатается в label htmlFor и aria-controls. */
  id: string;
  /** Текст подписи поля. */
  label: string;
  /** Класс-обёртка из макета: password-field (паспорт) или password-box (лендинг). */
  containerClassName: string;
  placeholder: string;
  autoComplete: 'current-password' | 'new-password';
  value: string;
  onChange: (value: string) => void;
  /** Ссылка-действие рядом с подписью (например «Забыли пароль?»). */
  labelRowExtra?: ReactNode;
  /** Ссылка на input для программного фокуса. */
  inputRef?: Ref<HTMLInputElement>;
  /** Идентификатор блока ошибки для aria-describedby. */
  errorId?: string;
  required?: boolean;
};

/**
 * Поле пароля с кнопкой показа — перенос password-toggle из дизайн-макета
 * (переключение типа, aria-pressed, иконка глаза со слэшем в скрытом состоянии).
 */
export function PasswordField({
  id,
  label,
  containerClassName,
  placeholder,
  autoComplete,
  value,
  onChange,
  labelRowExtra,
  inputRef,
  errorId,
  required,
}: PasswordFieldProps) {
  const [visible, setVisible] = useState(false);

  const labelElement = (
    <label htmlFor={id}>{label}</label>
  );

  return (
    <>
      {labelRowExtra ? (
        <div className="label-row">
          {labelElement}
          {labelRowExtra}
        </div>
      ) : (
        labelElement
      )}
      <div className={containerClassName}>
        <input
          ref={inputRef}
          id={id}
          name="password"
          type={visible ? 'text' : 'password'}
          autoComplete={autoComplete}
          placeholder={placeholder}
          value={value}
          required={required}
          aria-describedby={errorId}
          onChange={(event) => onChange(event.target.value)}
        />
        <button
          type="button"
          className="password-toggle"
          aria-label={visible ? 'Скрыть пароль' : 'Показать пароль'}
          aria-controls={id}
          aria-pressed={visible}
          onClick={() => setVisible((current) => !current)}
        >
          <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path
              d="M2.5 12S6 6.5 12 6.5s9.5 5.5 9.5 5.5-3.5 5.5-9.5 5.5S2.5 12 2.5 12Z"
              stroke="currentColor"
              strokeWidth="1.5"
            />
            <circle cx="12" cy="12" r="2.5" stroke="currentColor" strokeWidth="1.5" />
            {!visible ? (
              <path className="eye-slash" d="m5 4 14 16" stroke="currentColor" strokeWidth="1.5" />
            ) : null}
          </svg>
        </button>
      </div>
    </>
  );
}
