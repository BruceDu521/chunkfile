# ChunkFile - 文件分块与合并工具

[English](README.md) | [简体中文](README_CN.md)

一个用 Go 语言编写的命令行工具，用于将大文件分割成多个小块，以及将这些小块重新合并成原始文件。特别适用于需要将超过 1GB 的大文件上传到有大小限制的云存储服务，然后在另一台设备上重新组合的场景。

## 功能特点

- 将大文件分割成指定大小的多个小块
- 将分割后的小块重新合并成原始文件
- 动态调整分块文件命名，支持任意数量的分块
- 自定义分块大小，支持多种单位（B、KB、MB、GB）
- 可选择在合并成功后自动清理分块文件
- 支持相对路径和绝对路径
- 跨平台兼容性（Windows、Linux、macOS）

## 安装方法

### 系统要求

- Go 1.18 或更高版本
- 项目使用 Go 1.23 工具链（如果可用，将自动选择）

### 使用 Go Install（推荐）

如果您的系统上已安装 Go，可以直接使用以下命令安装 ChunkFile：

```bash
go install github.com/BruceDu521/chunkfile/cmd/chunkfile@latest
```

确保您的 Go bin 目录在系统 PATH 中：
- Windows: `%USERPROFILE%\go\bin`
- Linux/macOS: `~/go/bin`

### 手动安装

1. 克隆仓库：
   ```bash
   git clone https://github.com/BruceDu521/chunkfile.git
   cd chunkfile
   ```

2. 构建可执行文件：
   ```bash
   go build -o chunkfile ./cmd/chunkfile
   ```

3. 将可执行文件移动到 PATH 目录中，或直接在当前位置使用。

## 使用方法

ChunkFile 提供两个主要命令：`split` 和 `merge`。

### 分割文件

```bash
chunkfile split --path <文件路径> [--size <大小>] [--unit <单位>]
```

或使用短标志：

```bash
chunkfile split -p <文件路径> [-s <大小>] [-u <单位>]
```

参数说明：
- `--path, -p`：要分割的文件路径（必需）
- `--size, -s`：每个分块的大小（默认：400）
- `--unit, -u`：大小单位（B、KB、MB、GB，不区分大小写，默认：MB）

例如，将一个大文件分割成每块 500MB 的小块：

```bash
chunkfile split --path "large_file.zip" --size 500 --unit MB
```

或分割成 1GB 的块：

```bash
chunkfile split -p "large_file.zip" -s 1 -u GB
```

这将生成一系列分块文件，如：`large_file.zip.chunk.0001`、`large_file.zip.chunk.0002` 等。

### 合并文件

```bash
chunkfile merge --path <分块文件前缀> [--clear]
```

或使用短标志：

```bash
chunkfile merge -p <分块文件前缀> [-c]
```

参数说明：
- `--path, -p`：分块文件的路径前缀（必需）
- `--clear, -c`：合并成功后删除分块文件（可选）

例如，合并之前分割的文件：

```bash
chunkfile merge --path "large_file.zip"
```

或合并并清理分块文件：

```bash
chunkfile merge -p "large_file.zip" -c
```

这将查找所有匹配模式 `large_file.zip.chunk.*` 的文件，正确排序后合并它们，生成原始文件 `large_file.zip`。

## 注意事项

- 合并时，程序会自动查找并排序所有匹配的分块文件
- 分块文件使用 `.chunk.XXXX` 作为后缀，其中 XXXX 是从 0001 开始的序号
- 分块后缀中的位数会根据总分块数动态确定
- 确保有足够的磁盘空间用于存储分块文件或合并后的文件
- 程序会自动将相对路径转换为绝对路径进行处理
- 兼容 Windows、Linux 和 macOS 系统

## 获取帮助

获取更多关于可用命令和选项的信息：

```bash
chunkfile --help
chunkfile split --help
chunkfile merge --help
```

## 许可证

[MIT 许可证](LICENSE) 