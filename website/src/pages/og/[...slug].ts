import { getCollection } from "astro:content";
import { OGImageRoute } from "astro-og-canvas";

const docs = await getCollection("docs");
const docPages = Object.fromEntries(docs.map(({ data, id }) => [id, { data }]));

const pages = {
  ...docPages,
  home: {
    data: {
      title: "templ-components",
      description:
        "94 server-rendered Go UI components. templ + HTMX + Tailwind v4. Zero JavaScript framework.",
    },
  },
};

export const { getStaticPaths, GET } = await OGImageRoute({
  pages,
  param: "slug",
  getImageOptions: (_path, page) => ({
    title: page.data.title,
    description: page.data.description,
    bgGradient: [[15, 23, 42]],
    border: { color: [37, 99, 235], width: 4 },
    padding: 80,
  }),
});
