import { useState } from 'react';

import { useAuth } from '../../../features/auth/AuthContext';
import { isValidEmail } from '../../../features/auth/validation';
import { CABINET_COPY } from '../content';

/**
 * «Профиль»: просмотр и правка логина с электронной почтой.
 * Валидация и сообщения перенесены из src/v4/cabinet-state.js дизайн-макета.
 */
export function ProfilePage() {
  const { user, updateProfile } = useAuth();

  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState({ login: '', email: '' });
  const [errors, setErrors] = useState<{ login?: string; email?: string }>({});
  const [feedback, setFeedback] = useState('');

  if (!user) {
    return null; // ProtectedRoute не пропустит сюда без входа
  }

  const startEditing = () => {
    setDraft({ login: user.login, email: user.email });
    setErrors({});
    setFeedback('');
    setEditing(true);
  };

  const cancelEditing = () => {
    setDraft({ login: user.login, email: user.email });
    setErrors({});
    setFeedback('');
    setEditing(false);
  };

  const changeDraft = (field: 'login' | 'email', value: string) => {
    setDraft((current) => ({ ...current, [field]: value }));
    setErrors((current) => ({ ...current, [field]: undefined }));
    setFeedback(CABINET_COPY.profile.dirty);
  };

  const saveProfile = () => {
    const nextErrors: { login?: string; email?: string } = {};
    if (!draft.login.trim()) {
      nextErrors.login = CABINET_COPY.profile.invalidLogin;
    } else if (draft.login.length > 50) {
      nextErrors.login = 'Логин должен быть не длиннее 50 символов.';
    }
    if (!isValidEmail(draft.email)) {
      nextErrors.email = CABINET_COPY.profile.invalidEmail;
    } else if (draft.email.length > 100) {
      nextErrors.email = CABINET_COPY.profile.invalidEmail;
    }
    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors);
      setFeedback('');
      return;
    }
    const result = updateProfile({ login: draft.login, email: draft.email });
    if (!result.ok) {
      setErrors({ email: result.error });
      setFeedback('');
      return;
    }
    setErrors({});
    setEditing(false);
    setFeedback(CABINET_COPY.profile.saved);
  };

  return (
    <>
      <div className="page-intro">
        <div>
          <h1>{CABINET_COPY.page.profileTitle}</h1>
          <p className="page-context">{CABINET_COPY.page.context}</p>
        </div>
        <p className="participant">
          {CABINET_COPY.page.participantLabel} <br />
          <strong>{user.login}</strong>
        </p>
      </div>

      <div id="cabinet-content" role="tabpanel" aria-labelledby="tab-profile">
        <section className="profile-surface" aria-labelledby="profile-data-title">
          <div className="profile-aside">
            <h2 id="profile-data-title">{CABINET_COPY.profile.title}</h2>
            <p>{CABINET_COPY.profile.description}</p>
          </div>
          <div className="profile-body">
            {editing ? (
              <form
                className="profile-form"
                noValidate
                onSubmit={(event) => {
                  event.preventDefault();
                  saveProfile();
                }}
              >
                <div className="field">
                  <label htmlFor="profile-login">{CABINET_COPY.profile.loginLabel}</label>
                  <input
                    id="profile-login"
                    name="login"
                    type="text"
                    required
                    maxLength={50}
                    autoComplete="username"
                    autoCapitalize="none"
                    spellCheck={false}
                    placeholder={CABINET_COPY.profile.loginPlaceholder}
                    value={draft.login}
                    aria-invalid={errors.login ? true : undefined}
                    aria-describedby={errors.login ? 'profile-login-error' : undefined}
                    onChange={(event) => changeDraft('login', event.target.value)}
                  />
                  {errors.login ? (
                    <p className="field-error" id="profile-login-error" role="alert">
                      {errors.login}
                    </p>
                  ) : null}
                </div>
                <div className="field">
                  <label htmlFor="profile-email">{CABINET_COPY.profile.emailLabel}</label>
                  <input
                    id="profile-email"
                    name="email"
                    type="email"
                    required
                    maxLength={100}
                    autoComplete="email"
                    inputMode="email"
                    autoCapitalize="none"
                    spellCheck={false}
                    placeholder={CABINET_COPY.profile.emailPlaceholder}
                    value={draft.email}
                    aria-invalid={errors.email ? true : undefined}
                    aria-describedby={errors.email ? 'profile-email-error' : undefined}
                    onChange={(event) => changeDraft('email', event.target.value)}
                  />
                  {errors.email ? (
                    <p className="field-error" id="profile-email-error" role="alert">
                      {errors.email}
                    </p>
                  ) : null}
                </div>
                <div className="form-actions">
                  <button type="submit" className="primary-button">
                    {CABINET_COPY.actions.saveProfile}
                  </button>
                  <button type="button" className="secondary-button" onClick={cancelEditing}>
                    {CABINET_COPY.actions.cancel}
                  </button>
                </div>
              </form>
            ) : (
              <>
                <dl className="profile-details">
                  <div>
                    <dt>{CABINET_COPY.profile.loginLabel}</dt>
                    <dd>{user.login}</dd>
                  </div>
                  <div>
                    <dt>{CABINET_COPY.profile.emailLabel}</dt>
                    <dd>{user.email}</dd>
                  </div>
                </dl>
                <button className="primary-button" type="button" onClick={startEditing}>
                  {CABINET_COPY.actions.editProfile}
                </button>
              </>
            )}
            <p className="profile-feedback" role="status">
              {feedback}
            </p>
          </div>
        </section>
      </div>
    </>
  );
}
