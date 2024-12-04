## 功能

- `/v1/chat/completions`对话接口
- `/v1/models`接口
- 负载均衡
- 定时检查，检查模型可配置

## 运行

- `/cmd/api_app/`代码主体
- `/cmd/http_server`Mock,前期用来学习，了解接口的

## config讲解

`model_mapping`：模型映射，实现调用方一个模型到提供方多个模型的映射

举例：调用`4o-mini`时，所有提供`4o-mini-xxxx`、`4o-mini-yyy`、`4o-mini`的提供者可以被轮询访问

```yaml
model_mapping:
  4o-mini: # 无须增加gpt-4o-mini到gpt-4o-mini本身的映射，程序会自动添加
    - 4o-mini-xxxx
    - 4o-mini-yyyy
```

---

`check_models`，需要定时检测的模型。

举例：传入`4o-mini`，`4o-mini-xxxx`、`4o-mini-yyy`、`4o-mini`等都会检测

> 模型那么多，自己用到的，其实也就几个，没必要全部都做检测

---

`provider`，后边跟的是type

- 目前`openai`，`siliconflow`，`openaiCompatible`是本质是一样的，会调用模型列表，最终模型列表是模型接口+手动输入
- 接口获取`openaiCompatibleNoFetchModel`不会自动获取模型，以手动填入为主

```yaml
name: provider1 # 随便写，没影响
url: https://api.provider1.com # 到域名这里截断一般就可以了
models: # 手动填入的模型
- model1_provider1
- model2_provider1
api_keys: # 密钥列表
- your_api_key_here
- your_api_key_here_2
```

