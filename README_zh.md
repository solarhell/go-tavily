# Go Tavily Client

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/solarhell/go-tavily?style=for-the-badge)](https://goreportcard.com/report/github.com/solarhell/go-tavily)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/solarhell/go-tavily)
[![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/solarhell/go-tavily?utm_source=oss&utm_medium=github&utm_campaign=solarhell%2Fgo-tavily&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)](https://coderabbit.ai)

[English](README.md) | 中文

轻量、类型安全的 [Tavily API](https://docs.tavily.com) Go 客户端，基于 Go 1.26+。

## 安装

```bash
go get github.com/solarhell/go-tavily
```

## 快速开始

```go
client := tavily.New("tvly-your-api-key")

resp, err := client.Search(ctx, &tavily.SearchParams{
    Query: "Go programming language",
})
```

## 搜索

```go
resp, err := client.Search(ctx, &tavily.SearchParams{
    Query:         "AI news",
    SearchDepth:   tavily.SearchDepthAdvanced,
    Topic:         tavily.TopicNews,
    TimeRange:     tavily.TimeRangeWeek,
    MaxResults:    10,
    IncludeAnswer: new(tavily.IncludeAnswerAdvanced),
})
```

零值字段会从请求中省略，API 使用其服务端默认值。

## 国家筛选

使用 ISO 3166-1 alpha-2 国家代码来提升特定国家的搜索结果权重。国家筛选仅在 `topic` 为 `general`（或未设置）时可用。

```go
resp, err := client.Search(ctx, &tavily.SearchParams{
    Query:   "latest tech news",
    Country: new(tavily.CountryUS),
})
```

这个设计是为 LLM 工具调用场景考虑的：当你把搜索能力作为工具暴露给 LLM 时，工具的 schema 只需要声明"ISO 3166-1 alpha-2 国家代码"即可，无需枚举 160 多个国家名称字符串，因为 LLM 天然了解 ISO 标准。

> **注意：** Tavily 支持约 160 个国家，是完整 ISO 3166-1 alpha-2 标准（249 个代码）的子集。不支持的代码会在调用时被拒绝。

## 内容提取

```go
resp, err := client.Extract(ctx, &tavily.ExtractParams{
    URLs:   []string{"https://example.com"},
    Format: tavily.FormatMarkdown,
})
```

## 客户端配置

```go
client := tavily.New("tvly-your-api-key",
    tavily.WithBaseURL("https://custom.api.com"),
    tavily.WithHTTPClient(&http.Client{Timeout: 45 * time.Second}),
)
```

如果 API key 为空，会自动读取 `TAVILY_API_KEY` 环境变量。

## 错误处理

```go
resp, err := client.Search(ctx, &tavily.SearchParams{Query: "test"})
if err != nil {
    var apiErr *tavily.APIError
    if errors.As(err, &apiErr) {
        switch {
        case apiErr.IsUnauthorized():
            // API key 无效 (401)
        case apiErr.IsRateLimit():
            // 触发限流 (429)
        case apiErr.IsPlanLimitExceeded():
            // 套餐用量超限 (432)
        case apiErr.IsPayGoLimitExceeded():
            // 按量付费额度超限 (433)
        case apiErr.IsBadRequest():
            // 请求参数无效 (400)
        }
    }
}
```

## 测试

```bash
go test -v -race ./...
```

## 链接

- [Tavily API 文档](https://docs.tavily.com)
