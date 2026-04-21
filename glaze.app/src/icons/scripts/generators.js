import { writeFileSync } from "fs";

import { toPascalCase } from "./utils.js";

export const generateIconComponentFile = (icons, { dir }) => {
  const indexContent = [
    "// @ts-nocheck;",
    "",
    "import Icon from './Icon';",
    "",
    "type Props = {",
    "  name?: string;",
    "  size?: number;",
    "  color?: string;",
    "  stroke?: string;",
    "  onClick?: () => void",
    "  cursor?: string;",
    "  className?: string;",
    "};",
    "",
    icons
      ?.map(
        (icon) =>
          `export const ${toPascalCase(
            icon,
          )} = (props: Props) => <Icon {...props} name="${icon}" />;`,
      )
      .join("\n"),
  ].join("\n");

  writeFileSync(`${dir}/index.tsx`, indexContent);
  console.log(`Icon component file created! ✅`);
};

export const generateWebIconMap = (icons, { dir }) => {
  const iconMapContent = [
    icons
      ?.map(
        (icon) =>
          `// @ts-ignore;\nimport ${toPascalCase(icon)} from './${icon}.svg?react';`,
      )
      .join("\n"),
    "",
    "export default {",
    icons?.map((icon) => `"${icon}": ${toPascalCase(icon)}, `).join("\n"),
    "};",
  ].join("\n");

  writeFileSync(`${dir}/IconMap.tsx`, iconMapContent);
  console.log("Web icon map created! ✅");
};
