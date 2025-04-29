import eslint from '@eslint/js';
import tsEslint from 'typescript-eslint';
import eslintCdkPlugin from 'eslint-cdk-plugin';

export default tsEslint.config(
  {
    ignores: ['cdk.out', 'node_modules', '**/*.js', '**/*.json', '**/*.d.ts'],
  },
  // ESLintの推奨設定を適用
  eslint.configs.recommended,
  ...tsEslint.configs.recommended,
  ...tsEslint.configs.stylistic,
  {
    // lib と bin ディレクトリ配下の typescript ファイルを対象にする
    files: ['lib/**/*.ts', 'bin/*.ts'],
    languageOptions: {
      parserOptions: {
        projectService: true,
        project: './tsconfig.json',
      },
    },
    plugins: {
      cdk: eslintCdkPlugin,
    },
    rules: {
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      ...eslintCdkPlugin.configs.recommended.rules,
    },
  },
);
