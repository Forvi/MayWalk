import { useRef, useState } from 'react';
import { Link } from 'react-router-dom';

import { PassportLayout } from '../PassportLayout';
import { useAuth } from '../AuthContext';
import { validateAuthFields } from '../validation';

/**
 * «Восстановление пароля» — форма из auth-views.js дизайн-макета (режим recovery).
 * Демо-режим без бэкенда: адрес проверяется валидацией, «письмо» отправляется всегда.
 */
export function RecoveryPage() {
  const { recoverAccess } = useAuth();

  const [email, setEmail] = useState('');
  const [errors, setErrors] = useState<{ email?: string }>({});
  const [status, setStatus] = useState('');
  const [vkDialogOpen, setVkDialogOpen] = useState(false);

  const emailInputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus('');
    const nextErrors = validateAuthFields({ email }, 'recovery');
    setErrors(nextErrors);
    if (nextErrors.email) {
      emailInputRef.current?.focus();
      return;
    }
    const result = recoverAccess({ email });
    setStatus(result.message);
  };

  return (
    <PassportLayout
      activeTab="login"
      vkDialogOpen={vkDialogOpen}
      onCloseVkDialog={() => setVkDialogOpen(false)}
    >
      <div className="form-heading">
        <h1 id="form-title" tabIndex={-1}>
          Восстановить пароль
        </h1>
        <p>Укажите почту, которую использовали при регистрации.</p>
      </div>

      <form noValidate onSubmit={handleSubmit}>
        <div className="field">
          <label htmlFor="auth-recovery-email">Электронная почта</label>
          <input
            ref={emailInputRef}
            id="auth-recovery-email"
            name="email"
            type="email"
            autoComplete="email"
            maxLength={100}
            placeholder="you@example.ru"
            required
            autoCapitalize="none"
            spellCheck={false}
            inputMode="email"
            value={email}
            aria-invalid={errors.email ? true : undefined}
            aria-describedby={errors.email ? 'auth-recovery-email-error' : undefined}
            onChange={(event) => {
              setEmail(event.target.value);
              setErrors({});
              setStatus('');
            }}
          />
          {errors.email ? (
            <p id="auth-recovery-email-error" className="field-error" role="alert">
              {errors.email}
            </p>
          ) : null}
        </div>

        <button className="primary-button" type="submit">
          Отправить письмо
        </button>

        {status ? (
          <p className="form-status" role="status" aria-live="polite" aria-atomic="true">
            {status}
          </p>
        ) : null}
      </form>

      <p className="signup-prompt">
        <Link className="text-button" to="/login">
          Вернуться ко входу
        </Link>
      </p>
    </PassportLayout>
  );
}
