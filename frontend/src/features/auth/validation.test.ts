import { describe, expect, it } from 'vitest';

import { firstErrorField, isValidEmail, validateAuthFields } from './validation';

describe('isValidEmail', () => {
  it('принимает корректные адреса', () => {
    expect(isValidEmail('walker@example.ru')).toBe(true);
    expect(isValidEmail('имя.фамилия@maywalk.ru')).toBe(true);
  });

  it('отвергает мусор', () => {
    expect(isValidEmail('')).toBe(false);
    expect(isValidEmail('no-at-sign')).toBe(false);
    expect(isValidEmail('a@')).toBe(false);
    expect(isValidEmail('a b@example.ru')).toBe(false);
  });
});

describe('validateAuthFields: логин', () => {
  it('пустой логин — ошибка дизайн-макета', () => {
    expect(validateAuthFields({ login: '   ' }, 'login').login).toBe('Введите логин.');
  });

  it('логин длиннее 50 символов — ошибка дизайн-макета', () => {
    expect(validateAuthFields({ login: 'a'.repeat(51) }, 'login').login).toBe(
      'Логин должен быть не длиннее 50 символов.',
    );
  });

  it('корректный логин не даёт ошибок', () => {
    expect(validateAuthFields({ login: 'aleksandra' }, 'login')).toEqual({});
  });
});

describe('validateAuthFields: почта', () => {
  it('пустая почта', () => {
    expect(validateAuthFields({ email: '' }, 'recovery').email).toBe(
      'Введите электронную почту.',
    );
  });

  it('некорректный формат или длиннее 100 символов', () => {
    expect(validateAuthFields({ email: 'abc' }, 'recovery').email).toBe(
      'Проверьте адрес электронной почты.',
    );
    expect(validateAuthFields({ email: 'a'.repeat(95) + '@ex.ru' }, 'recovery').email).toBe(
      'Проверьте адрес электронной почты.',
    );
  });
});

describe('validateAuthFields: пароль и согласие', () => {
  it('пустой пароль', () => {
    expect(validateAuthFields({ password: '' }, 'login').password).toBe('Введите пароль.');
  });

  it('согласие обязательно только при регистрации', () => {
    expect(validateAuthFields({ consent: false }, 'register').consent).toBe(
      'Для регистрации нужно ваше согласие.',
    );
    expect(validateAuthFields({ consent: false }, 'login').consent).toBeUndefined();
  });

  it('не проверяет поля, которых нет в форме', () => {
    // у формы входа нет почты — валидатор получает только логин и пароль
    expect(validateAuthFields({ login: 'aleksandra', password: 'secret' }, 'login')).toEqual({});
  });
});

describe('firstErrorField', () => {
  it('возвращает первое поле в порядке формы', () => {
    expect(firstErrorField({ password: 'Введите пароль.', login: 'Введите логин.' })).toBe(
      'login',
    );
  });

  it('null, когда ошибок нет', () => {
    expect(firstErrorField({})).toBeNull();
  });
});
