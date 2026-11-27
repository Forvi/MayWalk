/// <reference types="vite/client" />

interface ImportMetaEnv {
  /**
   * Базовый URL backend-интерфейса MayWalk (см. servers в api/openapi.yml).
   * Локально по умолчанию http://localhost:8080; в docker-сети — http://backend:8080.
   */
  readonly VITE_API_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
