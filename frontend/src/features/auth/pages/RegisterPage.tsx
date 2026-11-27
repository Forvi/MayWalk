import { useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { PassportLayout, VkButton } from '../PassportLayout';
import { useAuth } from '../AuthContext';
import { firstErrorField, validateAuthFields } from '../validation';
import { PasswordField } from '../../../shared/components/PasswordField';

/**
 * «Регистрация» — форма из auth-views.js дизайн-макета (режим register).
 * Создаёт аккаунт, входит и переносит в личный кабинет.
 */
export function RegisterPage() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [values, setValues] = useState({ login: '', email: '', password: '', consent: false });
  const [errors, setErrors] = useState<{
    login?: string;
    email?: string;
    password?: string;
    consent?: string;
  }>({});
  const [status, setStatus] = useState('');
  const [vkDialogOpen, setVkDialogOpen] = useState(false);

  const fieldRefs = {
    login: useRef<HTMLInputElement>(null),
    email: useRef<HTMLInputElement>(null),
    password: useRef<HTMLInputElement>(null),
  };

  const focusField = (field: string | null) => {
    if (field === 'login' || field === 'email' || field === 'password') {
      fieldRefs[field].current?.focus();
    }
  };

  const clearError = (field: 'login' | 'email' | 'password' | 'consent') => {
    setErrors((current) => (current[field] ? { ...current, [field]: undefined } : current));
    setStatus('');
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus('');
    const nextErrors = validateAuthFields(values, 'register');
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) {
      focusField(firstErrorField(nextErrors));
      return;
    }
    const result = register({
      login: values.login,
      email: values.email,
      password: values.password,
    });
    if (!result.ok) {
      setStatus(result.error);
      focusField('login');
      return;
    }
    navigate('/cabinet/journey', { replace: true });
  };

  return (
    <PassportLayout
      activeTab="register"
      vkDialogOpen={vkDialogOpen}
      onCloseVkDialog={() => setVkDialogOpen(false)}
    >
      <div className="form-heading">
        <h1 id="form-title" tabIndex={-1}>
          Начнём прогулку?
        </h1>
        <p>Создайте аккаунт, чтобы оформить участие.</p>
      </div>

      <VkButton onClick={() => setVkDialogOpen(true)} />

      <div className="divider">
        <span>или создайте аккаунт</span>
      </div>

      <form noValidate onSubmit={handleSubmit}>
        <div className="field">
          <label htmlFor="auth-register-login">Логин</label>
          <input
            ref={fieldRefs.login}
            id="auth-register-login"
            name="login"
            type="text"
            autoComplete="username"
            maxLength={50}
            placeholder="Придумайте логин"
            required
            autoCapitalize="none"
            spellCheck={false}
            value={values.login}
            aria-invalid={errors.login ? true : undefined}
            aria-describedby={errors.login ? 'auth-register-login-error' : undefined}
            onChange={(event) => {
              setValues((current) => ({ ...current, login: event.target.value }));
              clearError('login');
            }}
          />
          {errors.login ? (
            <p id="auth-register-login-error" className="field-error" role="alert">
              {errors.login}
            </p>
          ) : null}
        </div>

        <div className="field">
          <label htmlFor="auth-register-email">Электронная почта</label>
          <input
            ref={fieldRefs.email}
            id="auth-register-email"
            name="email"
            type="email"
            autoComplete="email"
            maxLength={100}
            placeholder="you@example.ru"
            required
            autoCapitalize="none"
            spellCheck={false}
            inputMode="email"
            value={values.email}
            aria-invalid={errors.email ? true : undefined}
            aria-describedby={errors.email ? 'auth-register-email-error' : undefined}
            onChange={(event) => {
              setValues((current) => ({ ...current, email: event.target.value }));
              clearError('email');
            }}
          />
          {errors.email ? (
            <p id="auth-register-email-error" className="field-error" role="alert">
              {errors.email}
            </p>
          ) : null}
        </div>

        <div className="field">
          <PasswordField
            id="auth-register-password"
            label="Пароль"
            containerClassName="password-field"
            placeholder="Придумайте пароль"
            autoComplete="new-password"
            value={values.password}
            onChange={(password) => {
              setValues((current) => ({ ...current, password }));
              clearError('password');
            }}
            errorId="auth-register-password-error"
            inputRef={fieldRefs.password}
          />
          {errors.password ? (
            <p id="auth-register-password-error" className="field-error" role="alert">
              {errors.password}
            </p>
          ) : null}
        </div>

        <label className="consent" htmlFor="auth-register-consent">
          <input
            type="checkbox"
            name="consent"
            id="auth-register-consent"
            required
            checked={values.consent}
            aria-describedby={errors.consent ? 'auth-register-consent-error' : undefined}
            onChange={(event) => {
              setValues((current) => ({ ...current, consent: event.target.checked }));
              clearError('consent');
            }}
          />
          <span>
            Я ознакомлен(а) с{' '}
            <a href="https://mayprogulka.ru/about/pdpol/" target="_blank" rel="noopener noreferrer">
              политикой конфиденциальности
            </a>{' '}
            и даю{' '}
            <a href="https://mayprogulka.ru/about/pdsogl/" target="_blank" rel="noopener noreferrer">
              согласие на обработку персональных данных
            </a>
            {errors.consent ? (
              <span id="auth-register-consent-error" className="field-error" role="alert">
                {errors.consent}
              </span>
            ) : null}
          </span>
        </label>

        <button className="primary-button" type="submit">
          Зарегистрироваться
        </button>

        {status ? (
          <p className="form-status" role="alert">
            {status}
          </p>
        ) : null}
      </form>

      <p className="signup-prompt">
        Есть аккаунт?{' '}
        <Link className="text-button" to="/login">
          Войти
        </Link>
      </p>
    </PassportLayout>
  );
}
