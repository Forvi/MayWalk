import { useState } from 'react';
import { Link } from 'react-router-dom';

import { useAuth } from '../auth/AuthContext';
import { firstErrorField, validateAuthFields } from '../auth/validation';
import { DialogCloseButton, ModalDialog } from '../../shared/components/ModalDialog';
import { BrandMark, FlowerMark } from '../../shared/components/BrandMark';
import { PageShell } from '../../shared/components/PageShell';
import { PasswordField } from '../../shared/components/PasswordField';

/**
 * Лендинг с регистрацией — перенос festival-register.html дизайн-макета:
 * разметка общая для сцены, оформлена оригинальной (зелёной) концепцией.
 * Вход в VK-кабинет демонстрируется диалогом (внешнего бэкенда в макете нет).
 */
export function LandingPage() {
  const { user, register } = useAuth();

  const [values, setValues] = useState({ login: '', email: '', password: '', consent: false });
  const [errors, setErrors] = useState<{
    login?: string;
    email?: string;
    password?: string;
    consent?: string;
  }>({});
  const [feedback, setFeedback] = useState('');
  const [vkDialogOpen, setVkDialogOpen] = useState(false);

  const clearError = (field: 'login' | 'email' | 'password' | 'consent') => {
    setErrors((current) => (current[field] ? { ...current, [field]: undefined } : current));
    setFeedback('');
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFeedback('');
    const nextErrors = validateAuthFields(values, 'register');
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) {
      const first = firstErrorField(nextErrors);
      document.getElementById(`v4-${first}`)?.focus();
      return;
    }
    const result = register({
      login: values.login,
      email: values.email,
      password: values.password,
    });
    if (!result.ok) {
      setFeedback(result.error);
      document.getElementById('v4-login')?.focus();
      return;
    }
    // Успех: register() открывает сессию, и вместо формы сразу
    // показывается блок «Вы уже с нами» с явным переходом в кабинет
    // (preview_is_explicit_separate_action из макета — без автоперехода).
  };

  return (
    <PageShell page="registration">
      <header className="global-header shell">
        <Link className="brand" to="/" aria-label="Майская прогулка — главная">
          <BrandMark />
          <span className="brand-name">
            майская
            <br />
            прогулка.
          </span>
        </Link>
        <div className="header-links">
          <a href="https://mayprogulka.ru/faq" target="_blank" rel="noopener noreferrer">
            Помощь
          </a>
          {user ? (
            <Link to="/cabinet/journey">Личный кабинет</Link>
          ) : (
            <Link to="/login">Войти</Link>
          )}
        </div>
      </header>

      <main>
        <section className="registration-lead shell" aria-labelledby="registration-title">
          <h1 id="registration-title">
            <span>ТВОЯ</span>{' '}
            <span className="title-tail">
              ПРОГУЛ<span className="word-ending">КА</span>
            </span>
          </h1>
          <FlowerMark className="lead-mark" size={105} />
        </section>

        <section className="registration-workspace shell" aria-labelledby="form-title">
          {user ? (
            /*
             * Вошедший участник не видит форму повторно: вместо него — короткий
             * блок возврата в кабинет (классы и токены дизайн-макета).
             */
            <>
              <div className="registration-workspace-head">
                <h2 id="form-title">Вы уже с нами</h2>
                <p className="return-login">
                  С возвращением, <strong>{user.login}</strong>!
                </p>
              </div>
              <div className="registration-columns">
                <div className="email-option">
                  <p className="option-label">Личный кабинет</p>
                  <p>Ваш маршрут и черновик анкеты ждут в личном кабинете.</p>
                  <div className="form-actions">
                    <Link className="preview-link" to="/cabinet/journey">
                      Перейти в личный кабинет
                    </Link>
                  </div>
                </div>
              </div>
            </>
          ) : (
            <>
              <div className="registration-workspace-head">
                <h2 id="form-title">Создать аккаунт</h2>
                <p className="return-login">
                  Уже с нами? <Link to="/login">Войти</Link>
                </p>
              </div>

          <div className="registration-columns">
            <aside className="social-option">
              <p className="option-label">Через ВКонтакте</p>
              <button className="vk-button" type="button" onClick={() => setVkDialogOpen(true)}>
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path
                    fill="currentColor"
                    d="M12.63 17.7C6.8 17.7 3.47 13.7 3.33 7.03h2.92c.1 4.9 2.26 6.98 3.97 7.41V7.03h2.75v4.23c1.68-.18 3.44-2.11 4.04-4.23h2.75c-.46 2.61-2.36 4.54-3.71 5.33 1.35.65 3.5 2.34 4.32 5.34h-3.03c-.65-2-2.26-3.54-4.37-3.75v3.75z"
                  />
                </svg>
                Войти через ВК
              </button>
              <p>Начните с аккаунта. Свой маршрут выберете в личном кабинете.</p>
            </aside>

            <div className="email-option">
              <p className="option-label">Или заполните форму</p>
              <form className="auth-form" noValidate onSubmit={handleSubmit}>
                <div className="field">
                  <label htmlFor="v4-login">Логин</label>
                  <input
                    id="v4-login"
                    name="login"
                    type="text"
                    required
                    maxLength={50}
                    autoComplete="username"
                    autoCapitalize="none"
                    spellCheck={false}
                    placeholder="Придумайте логин"
                    aria-describedby="v4-login-error"
                    value={values.login}
                    onChange={(event) => {
                      setValues((current) => ({ ...current, login: event.target.value }));
                      clearError('login');
                    }}
                  />
                  {errors.login ? (
                    <p className="field-error" id="v4-login-error" role="alert">
                      {errors.login}
                    </p>
                  ) : null}
                </div>

                <div className="field">
                  <label htmlFor="v4-email">Электронная почта</label>
                  <input
                    id="v4-email"
                    name="email"
                    type="email"
                    required
                    maxLength={100}
                    autoComplete="email"
                    inputMode="email"
                    autoCapitalize="none"
                    spellCheck={false}
                    placeholder="you@example.ru"
                    aria-describedby="v4-email-error"
                    value={values.email}
                    onChange={(event) => {
                      setValues((current) => ({ ...current, email: event.target.value }));
                      clearError('email');
                    }}
                  />
                  {errors.email ? (
                    <p className="field-error" id="v4-email-error" role="alert">
                      {errors.email}
                    </p>
                  ) : null}
                </div>

                <div className="field full">
                  <PasswordField
                    id="v4-password"
                    label="Пароль"
                    containerClassName="password-box"
                    placeholder="Придумайте пароль"
                    autoComplete="new-password"
                    value={values.password}
                    onChange={(password) => {
                      setValues((current) => ({ ...current, password }));
                      clearError('password');
                    }}
                    errorId="v4-password-error"
                  />
                  {errors.password ? (
                    <p className="field-error" id="v4-password-error" role="alert">
                      {errors.password}
                    </p>
                  ) : null}
                </div>

                <label className="consent" htmlFor="v4-consent">
                  <input
                    type="checkbox"
                    name="consent"
                    id="v4-consent"
                    required
                    checked={values.consent}
                    aria-describedby="v4-consent-error"
                    onChange={(event) => {
                      setValues((current) => ({ ...current, consent: event.target.checked }));
                      clearError('consent');
                    }}
                  />
                  <span>
                    Я ознакомлен(а) с{' '}
                    <a
                      href="https://mayprogulka.ru/about/pdpol/"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      политикой конфиденциальности
                    </a>{' '}
                    и даю{' '}
                    <a
                      href="https://mayprogulka.ru/about/pdsogl/"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      согласие на обработку персональных данных
                    </a>
                    {errors.consent ? (
                      <span id="v4-consent-error" className="field-error" role="alert">
                        {errors.consent}
                      </span>
                    ) : null}
                  </span>
                </label>

                <div className="form-actions">
                  <button className="primary-button" type="submit">
                    Зарегистрироваться
                  </button>
                </div>

                {feedback ? (
                  <p className="form-feedback" role="status">
                    {feedback}
                  </p>
                ) : null}
              </form>
            </div>
          </div>
            </>
          )}
        </section>

        <section className="registration-outro shell">
          <p className="outro-title">Каждый идёт своим ходом</p>
          <p>Лесные тропы, городские маршруты и день, который хочется провести на ногах.</p>
        </section>
      </main>

      <footer className="site-footer shell">
        <p>Майская прогулка · Екатеринбург</p>
        <div className="site-footer-links">
          <a href="https://mayprogulka.ru/about/pdpol/" target="_blank" rel="noopener noreferrer">
            Политика конфиденциальности
          </a>
          <a href="https://mayprogulka.ru/about/pdsogl/" target="_blank" rel="noopener noreferrer">
            Обработка персональных данных
          </a>
        </div>
      </footer>

      <ModalDialog
        className="dialog vk-dialog"
        open={vkDialogOpen}
        onClose={() => setVkDialogOpen(false)}
        titleId="vk-title"
        descriptionId="vk-description"
      >
        <DialogCloseButton label="Закрыть окно" onClick={() => setVkDialogOpen(false)} />
        <h2 id="vk-title">Вход через ВКонтакте</h2>
        <p id="vk-description">
          Для входа в личный кабинет используйте свой аккаунт ВКонтакте.
        </p>
        <div className="dialog-actions">
          <button
            className="primary-button"
            type="button"
            onClick={() => setVkDialogOpen(false)}
          >
            Понятно
          </button>
        </div>
      </ModalDialog>
      {/* Навигация в кабинет остаётся явным действием (как в макете:
          preview_is_explicit_separate_action) — без автоперехода. */}
    </PageShell>
  );
}