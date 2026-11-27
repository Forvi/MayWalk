import js from '@eslint/js';
import globals from 'globals';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  // Глобальные игноры
  {
    ignores: ['dist/', 'node_modules/', 'src/assets/'],
  },

  // Базовые + TypeScript-правила
  js.configs.recommended,
  ...tseslint.configs.recommended,

  // Правила для React-кода
  reactHooks.configs['flat/recommended'],
  reactRefresh.configs.vite,

  // Доступность: публичный сайт, разметка дизайнера aria-friendly
  jsxA11y.flatConfigs.recommended,

  // Общие настройки
  {
    files: ['**/*.{js,ts,jsx,tsx}'],
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.serviceworker,
      },
    },
    rules: {
      'no-console': 'warn',
      'no-debugger': 'warn',
      'prefer-const': 'error',
      eqeqeq: ['error', 'always', { null: 'ignore' }],

      // TypeScript
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
      ],
      '@typescript-eslint/consistent-type-imports': [
        'error',
        { prefer: 'type-imports', fixStyle: 'inline-type-imports' },
      ],
      '@typescript-eslint/no-import-type-side-effects': 'error',

      // React
      'react-hooks/exhaustive-deps': 'error',
      'react-hooks/set-state-in-effect': 'error',
      'react-refresh/only-export-components': 'off',

      // Доступность
      'jsx-a11y/anchor-is-valid': 'off', // допустимы href="#" и внешние ссылки дизайн-макета
    },
  },

  // Конфиги — среда Node
  {
    files: ['vite.config.ts', 'eslint.config.js'],
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
  },
);
