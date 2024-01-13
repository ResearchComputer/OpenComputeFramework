import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Open Compute Framework',
			social: {
				github: 'https://github.com/autoai-org/opencomputeframework',
				discord: 'https://discord.gg/PgGb4z4Jve',
			},
			sidebar: [
				{
					label: 'Guides',
					items: [
						{ label: "Introduction", link: "/guides/intro/" },
						{ label: 'ML Inference', link: '/guides/ml_inference/' },
					],
				}

			],
		}),
	],
});
