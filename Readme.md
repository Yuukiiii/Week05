# Week 05

## Homework

1. 参考 Hystrix 实现一个滑动窗口计数器。

## 思路

仅实现请求数量统计功能。统计成功量，失败量，错误率等等可以在这个上面加。

启动一个简单的 HTTP server，以 100ms 为窗口大小，10 为窗口数量，统计每个窗口内的请求数量。

使用 wrk 工具模拟并发请求

