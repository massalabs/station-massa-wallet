/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BASE_APP: string;
  readonly VITE_BASE_API: string;
  readonly VITE_ENV: string;
  readonly VITE_CI_TEST: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
