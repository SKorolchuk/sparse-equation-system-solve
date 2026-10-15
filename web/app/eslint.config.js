import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([
  globalIgnores(['dist', 'public/wasm_exec.js', 'public/solver.wasm']),
  
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
    ],
    plugins: {
      'react-refresh': reactRefresh,
    },
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.es2022,
        // Allowing standard WebAssembly execution context global types
        WebAssembly: 'readonly',
      },
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/no-unused-vars': ['error', { 
        'argsIgnorePattern': '^_',
        'varsIgnorePattern': '^_' 
      }],
      '@typescript-eslint/explicit-function-return-type': ['warn', {
        'allowExpressions': true,
        'allowTypedFunctionExpressions': true
      }],
      '@typescript-eslint/no-empty-interface': 'error',
      '@typescript-eslint/no-misused-new': 'error',
      '@typescript-eslint/consistent-type-imports': ['error', { 'prefer': 'type-imports' }],
      
      // React Hooks exhaustive dependencies checks to avoid stale closures in component lifecycle
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'error',
      
      // Strict layout scoping for hot module reloading (HMR) within components
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
      
      'no-console': ['warn', { allow: ['warn', 'error', 'info'] }],
      'no-debugger': 'error',
      'no-duplicate-imports': 'error',
      'eqeqeq': ['error', 'always', { 'null': 'ignore' }],
    },
  },
])
