import { useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { PassportLayout, VkButton } from '../PassportLayout';
import { useAuth } from '../AuthContext';
import { firstErrorField, validateAuthFields } from '../validation';
import { PasswordField } from '../../../shared/components/PasswordField';

/**
 * «Вход» — форма из auth-views.js дизайн-макета (режим login).
 * После успешного входа возвращает на адрес, с которого отправили к логину.
 */
export function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [values, setValues] = useState({ login: '', password: '' });
  const [errors, setErrors] = useState<{ login?: string; password?: string }>({});
  const [status, setStatus] = useState('');
  const [vkDialogOpen, setVkDialogOpen] = useState(false);

  const loginInputRef = useRef<HTMLInputElement>(null);
  const passwordInputRef = useRef<HTMLInputElement>(null);

  const focusField = (field: string | null) => {
    if (field === 'login') {
      loginInputRef.current?.focus();
    } else if (field === 'password') {
      passwordInputRef.current?.focus();
    }
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus('');
    const nextErrors = validateAuthFields(
      { login: values.login, password: values.password },
      'login',
    );
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) {
      focusField(firstErrorField(nextErrors));
      return;
    }
    const result = login({ login: values.login, password: values.password });
    if (!result.ok) {
      setStatus(result.error);
      focusField('login');
      return;
    }
    const from = (location.state as { from?: string } | null)?.from ?? '/cabinet/journey';
    navigate(from, { replace: true });
  };

  return (
    <PassportLayout
      activeTab="login"
      vkDialogOpen={vkDialogOpen}
      onCloseVkDialog={() => setVkDialogOpen(false)}
    >
      <div className="form-heading">
        <h1 id="form-title" tabIndex={-1}>
          С возвращением!
        </h1>
        <p>Войдите, чтобы выбрать свой маршрут.</p>
      </div>

      <VkButton onClick={() => setVkDialogOpen(true)} />

      <div className="divider">
        <span>или по логину</span>
      </div>

      <form noValidate onSubmit={handleSubmit}>
        <div className="field">
          <label htmlFor="auth-login">Логин</label>
          <input
            ref={loginInputRef}
            id="auth-login"
            name="login"
            type="text"
            autoComplete="username"
            maxLength={50}
            placeholder="Введите логин"
            required
            autoCapitalize="none"
            spellCheck={false}
            value={values.login}
            aria-invalid={errors.login ? true : undefined}
            aria-describedby={errors.login ? 'auth-login-error' : undefined}
            onChange={(event) => {
              setValues((current) => ({ ...current, login: event.target.value }));
              setErrors((current) => ({ ...current, login: undefined }));
            }}
          />
          {errors.login ? (
            <p id="auth-login-error" className="field-error" role="alert">
              {errors.login}
            </p>
          ) : null}
        </div>

        <div className="field">
          <PasswordField
            id="auth-password"
            label="Пароль"
            containerClassName="password-field"
            placeholder="Введите пароль"
            autoComplete="current-password"
            value={values.password}
            onChange={(password) => {
              setValues((current) => ({ ...current, password }));
              setErrors((current) => ({ ...current, password: undefined }));
            }}
            labelRowExtra={
              <Link className="text-button" to="/recovery">
                Забыли пароль?
              </Link>
            }
            errorId="auth-password-error"
            inputRef={passwordInputRef}
          />
          {errors.password ? (
            <p id="auth-password-error" className="field-error" role="alert">
              {errors.password}
            </p>
          ) : null}
        </div>

        <button className="primary-button" type="submit">
          Войти
        </button>

        {status ? (
          <p className="form-status" role="alert">
            {status}
          </p>
        ) : null}
      </form>

      <p className="signup-prompt">
        Впервые с нами?{' '}
        <Link className="text-button" to="/register">
          Создать аккаунт
        </Link>
      </p>
    </PassportLayout>
  );
}
