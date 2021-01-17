## Introduce

管理V2ray工具

- 支持订阅链接
- 自动生成v2ray配置文件

## Install

[releases](https://github.com/vczyh/v2rayT/releases)

## Usage

```bash
v2rayT -url 订阅链接 -v2ray v2ray可执行文件路径
```

| 参数        | 说明                                                | 默认                            |
| :---------- | --------------------------------------------------- | ------------------------------- |
| url         | 订阅链接，支持vmess                                 | 必须                            |
| v2ray       | v2ray可执行文件路径                                 | 必须                            |
| v2rayConfig | v2ray配置文件路径，你不需要创建它，v2rayT会自动创建 | $HOME/.config/v2ray/config.json |
| logPath     | 日志目录，例如`/var/log/v2ray`                      | 当前目录                        |
| socksPort   | socks5端口                                          | 1080                            |

除此之外，可以使用交互模式。

```shell
v2rayT
请输入订阅链接：xxx
请输入v2ray可执行文件路径：xxx
```

这样只会提示输入必要的参数，非必要参数不会提示，所以你可以混合使用。

```bash
v2rayT -v2rayConfig ./config.json -logPath /var/log/v2ray
请输入订阅链接：xxx
请输入v2ray可执行文件路径：xxx
```

## Other

仅供交流学习。