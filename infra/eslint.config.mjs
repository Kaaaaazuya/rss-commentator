import tsEslint from "typescript-eslint";
import tsParser from "@typescript-eslint/parser";
import eslintCdkPlugin from "eslint-cdk-plugin";

export default tsEslint.config({
	files: ["lib/**/*.ts", "bin/*.ts"],
	ignores: ["**/*.d.ts"],
	languageOptions: {
		parser: tsParser,
		parserOptions: {
			project: "./tsconfig.json",
		},
	},
	plugins: {
		cdk: eslintCdkPlugin,
	},
	rules: {
		...eslintCdkPlugin.configs.recommended.rules,
		"cdk/no-public-class-fields": "warn",
		"cdk/no-variable-construct-id": "warn",
	},
});
