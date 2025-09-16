[![buymeacoffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-dumindu-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=ffffff&labelColor=333333)](https://www.buymeacoffee.com/dumindu)

# 🦄 E25DX

Your next documentation theme...

https://github.com/dumindu/E25DX/assets/1280428/00935567-6d50-4e1b-bb0c-cbd159ad37f0

## 🌱 Inspiration

- Build a simple, modular but customizable technical documentation.
- Build the website without using node/npm or any external frameworks(CSS, JS, icon, font).
- Demo
  - [![learning-rust.github.io](https://img.shields.io/github/stars/learning-rust/learning-rust.github.io?style=for-the-badge&logo=rust&label=learning-rust.github.io&labelColor=333333&color=F46623)](https://learning-rust.github.io)
  - [![learning-cloud-native-go.github.io](https://img.shields.io/github/stars/learning-cloud-native-go/learning-cloud-native-go.github.io?style=for-the-badge&logo=go&logoColor=ffffff&label=learning-cloud-native-go.github.io&labelColor=333333&color=00ADD8)](https://learning-cloud-native-go.github.io)

## ✅ Features

- [x] Responsive and adaptive layouts.
- [x] Built-in light and dark modes.
- [x] Customizable sidebars using Hugo data templates.
- [ ] Support for multiple documentation sets.
- [ ] Implement a menu via Hugo configs.
- [ ] Integrate search or AI assistant.

## 🚀 Getting Started

1. [Install Hugo](https://gohugo.io/installation/).
    - Use prebuilt binaries via https://github.com/gohugoio/hugo/releases/latest
    - Verify the installation via `hugo version`

2. Create a new Hugo site with the E25DX theme.
    1. Foundation
        ```shell
        hugo new site newsite
        cd newsite

        git init
        git submodule add https://github.com/dumindu/E25DX.git themes/E25DX

        rm hugo.toml && cp themes/E25DX/hugo.yaml .
        cp themes/E25DX/static/sw.js ./static/
        cp themes/E25DX/static/manifest.json ./static/
        ```
       > 💡 We don't need generated `archetypes`, `assets`, `i18n`, `layouts` folders in the `newsite` folder

   2. Configuration
      - `newsite/hugo.yaml`
        - languages: en:title and en:description
        - params: project, author and home
      - `newsite/static/sw.js`
        - cacheName prefix: ex. todo-
      - `newsite/static/manifest.json`
        - name, short_name, description

   3. Sample content and sidebar
      - add `newsite/content/en/docs/_overview.md` for overview
        ```markdown
        ---
        title: Overview
        url: "docs/overview"
        aliases:
        - "/docs"
        ---
        ```
      - add `newsite/content/en/docs/a1.hello-world.md` as the first page of first section
        ```markdown
        ---
        title: Hello World
        slug: hello-world
        ---
        ```
      - add `newsite/content/en/docs/b1.modules.md` as the first page of second section
        ```markdown
        ---
        title: Modules
        slug: modules
        ---
        ```
      - add `newsite/data/en/docs/sidebar.yml` for section titles and page titles
        ```markdown
        - title: Documentation
          pages:
            - title: Overview

        - title: Basics
          pages:
            - title: Hello World

        - title: Beyond The Basics
          pages:
            - title: Modules
        ```

   4. Run `hugo server`
      > 💡 Change `newsite/hugo.yaml` -> `enableGitInfo: false`, if you want to run `hugo server` before commit the changes.
