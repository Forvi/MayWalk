/**
 * Валидация форм — перенос src/auth-state.js из дизайн-макета.
 * Сообщения об ошибках сохранены дословно.
 */

export function isValidEmail(value: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(value.trim());
}

export type FieldErrors = Partial<Record<'login' | 'email' | 'password' | 'consent', string>>;

export type AuthFormValues = {
  login?: string;
  email?: string;
  password?: string;
  consent?: boolean;
};

/** Режим формы задаёт набор обязательных полей (как в макете). */
export type AuthMode = 'login' | 'register' | 'recovery' | 'profile';

export function validateAuthFields(values: AuthFormValues, mode: AuthMode): FieldErrors {
  const errors: FieldErrors = {};
  const { login, email, password, consent } = values;

  if (login !== undefined) {
    if (!login.trim()) {
      errors.login = 'Введите логин.';
    } else if (login.length > 50) {
      errors.login = 'Логин должен быть не длиннее 50 символов.';
    }
  }

  if (email !== undefined) {
    if (!email.trim()) {
      errors.email = 'Введите электронную почту.';
    } else if (!isValidEmail(email) || email.length > 100) {
      errors.email = 'Проверьте адрес электронной почты.';
    }
  }

  if (password !== undefined && !password) {
    errors.password = 'Введите пароль.';
  }

  if (mode === 'register' && consent !== undefined && !consent) {
    errors.consent = 'Для регистрации нужно ваше согласие.';
  }

  return errors;
}

/** Первое поле с ошибкой — на него переводится фокус после сабмита. */
export function firstErrorField(errors: FieldErrors): keyof FieldErrors | null {
  const order = ['login', 'email', 'password', 'consent'] as const;
  for (const field of order) {
    if (errors[field]) {
      return field;
    }
  }
  return null;
}
