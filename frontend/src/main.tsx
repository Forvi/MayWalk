import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import App from './App';
import './styles/global.css';
import './styles/controls.css';
import './styles/fonts.css';
import './styles/pages/stage-original.css';
import './styles/pages/auth-original.css';

const container = document.getElementById('root');

if (!container) {
  throw new Error('Не найден контейнер #root для монтирования приложения');
}

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
