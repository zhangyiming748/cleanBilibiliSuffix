<div align="center">

# 🧹 cleanBilibiliSuffix

**一键清理哔哩哔哩下载视频文件名中的多余前缀**

[![Release](https://img.shields.io/github/v/release/zhangyiming748/cleanBilibiliSuffix?label=Release&logo=github)](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases)
[![Go Release](https://img.shields.io/github/actions/workflow/status/zhangyiming748/cleanBilibiliSuffix/gorelease.yml?label=Go%20Release&logo=githubactions&logoColor=white)](https://github.com/zhangyiming748/cleanBilibiliSuffix/actions/workflows/gorelease.yml)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache--2.0-blue)](LICENSE)

[功能特性](#✨-功能特性) · [快速下载](#📥-快速下载) · [使用方法](#🚀-使用方法) · [项目结构](#📁-项目结构) · [从源码构建](#🔨-从源码构建)

</div>

---

## 💡 简介

哔哩哔哩客户端下载的视频，文件名经常带有与视频内容无关的画质描述前缀：

```text
4K修复 1080P高清修复 4K120帧 …
```

`cbs`（**c**lean**B**ili**S**uffix）会递归扫描指定目录下的所有视频文件，自动剥离这些前缀，让文件名只保留真正的视频标题：

```text
4K修复 1080P高清修复 某视频标题.mp4   →   某视频标题.mp4
```

单个静态二进制文件，无运行时依赖，Linux / macOS / Windows 开箱即用。

## ✨ 功能特性

| 特性 | 说明 |
| ---- | ---- |
| 🔍 **递归扫描** | 自动遍历目标目录的所有子目录，不遗漏任何层级的视频文件 |
| 🎬 **真视频识别** | 通过文件头魔数检测视频，不依赖扩展名；rmvb/rm/vob/flv/ts/m2ts 等格式用扩展名兜底 |
| 🧩 **多重前缀处理** | 循环剥离关键词前缀，能处理叠加的多层前缀 |
| 🪶 **自动去空格** | 删除前缀后，开头的空格（含全角空格）一并删除 |
| 🛡️ **安全保护** | 目标文件已存在时拒绝覆盖；清理后文件名为空时跳过，不会误删 |
| ⚙️ **关键词可配置** | 要删除的前缀写在文本文件里，每行一个，随需增减 |

## 📥 快速下载

前往 [Releases 页面](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases) 下载最新版本，或直接按下表获取：

| 平台 | 架构 | 下载链接 |
|------|------|----------|
| Linux | amd64 | [cbs_linux_amd64](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_linux_amd64) |
| Linux | arm64 | [cbs_linux_arm64](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_linux_arm64) |
| macOS | amd64 | [cbs_darwin_amd64](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_darwin_amd64) |
| macOS | arm64 | [cbs_darwin_arm64](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_darwin_arm64) |
| Windows | amd64 | [cbs_windows_amd64.exe](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_windows_amd64.exe) |
| Windows | arm64 | [cbs_windows_arm64.exe](https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_windows_arm64.exe) |

**一键下载命令：**

```bash
# Linux/macOS
wget https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/') -O cbs && chmod +x cbs

# Windows PowerShell
Invoke-WebRequest -Uri "https://github.com/zhangyiming748/cleanBilibiliSuffix/releases/latest/download/cbs_windows_amd64.exe" -OutFile "cbs.exe"
```

## 🚀 使用方法

### 1. 准备关键词文件

创建一个文本文档（如 `keywords.txt`），每行写一个要删除的前缀关键词：

```text
4K修复
4K120帧
1080P高清修复
```

空行会被自动忽略。仓库自带的 [keywords.txt](keywords.txt) 可以直接作为模板使用。

### 2. 运行

```bash
cbs --keyword keywords.txt --dir /path/to/your/videos
```

### 3. 查看结果

程序边扫描边处理，实时输出每个文件的处理结果：

```text
重命名成功: 4K修复 1080P高清修复 某视频标题.mp4 -> 某视频标题.mp4
无需修改: 本来就是干净的名字.mp4
```

### 命令行参数

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--keyword` | `-k` | ✅ | 关键词文本文档的路径，每行一个关键词 |
| `--dir` | `-d` | ✅ | 要处理的视频文件所在根目录 |
| `--version` | `-v` | — | 打印版本号（包含构建时间和 commit 信息） |

## 📁 项目结构

```text
cleanBilibiliSuffix/
├── core/
│   ├── files.go          # 目录递归遍历与视频文件识别
│   └── rename.go         # 关键词加载与文件名清洗/重命名
├── main.go               # 命令行入口（cobra）
├── keywords.txt          # 默认关键词模板
└── .goreleaser.yml       # 多平台发布配置
```

## 🔨 从源码构建

需要 Go 1.26+：

```bash
git clone https://github.com/zhangyiming748/cleanBilibiliSuffix.git
cd cleanBilibiliSuffix
go build -o cbs .
```

正式发布由 GitHub Actions + [GoReleaser](https://goreleaser.com) 自动完成：推送 `vX.Y.Z` 格式的 tag 即触发多平台构建并发布到 Releases，版本号、构建时间、commit 哈希会通过编译参数注入二进制。

## ⚠️ 注意事项

- 重命名是真实执行的（直接改名，不可撤销），建议首次使用时先在小范围目录上验证关键词配置
- 只处理视频文件，其他类型的文件不会被触碰
- 关键词只匹配文件名的**开头**，不会误删文件名中间或结尾的相同字样

## 🤝 贡献

欢迎提交 Issue 和 Pull Request：

1. Fork 本仓库并创建特性分支：`git checkout -b feature/xxx`
2. 提交改动：`git commit -m 'feat: add xxx'`
3. 推送并发起 PR：`git push origin feature/xxx`

## 📄 许可证

本项目基于 [Apache License 2.0](LICENSE) 开源。

---

<div align="center">

如果这个项目对你有帮助，欢迎 ⭐ Star 支持一下！

</div>
