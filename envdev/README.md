# Dev Basic

开发的基础环境搭建

## go-kratos 开发

go-kratos 文档: https://go-kratos.dev/docs/

vscode 插件: `Makefile Tools`

nodejs 插件: `Prettier - Code formatter`, `ESLint`

## python-miniconda 开发

conda 常用命令
```bash
conda env list
```

## 备注

vscode 若不能将 docker目录下的文件当作 Dockerfile，请将以下代码复制到 .vscode/settings.json
```json
"files.associations": {
  "**/Dockerfile*": "dockerfile"
}
```