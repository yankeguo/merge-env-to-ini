# merge-env-to-ini

从环境变量往 INI 格式的文件内合并键值对，当前用于 acicn/php 系列镜像

## 用法

假设已有环境变量

```
DEMO_section1__1="key_a=val_a"
DEMO_section2__1="key_b=val_b"
DEMO_section3__1="key_b=val_b"
```

已有文件 `somefile.ini`

```ini
[section1]
key_a=val_0

[section2]
key_a=val_0
```

执行命令

```shell
merge-env-to-ini --from DEMO_ --to somefile.ini
```

会将文件内容改写为

```ini
[section1]
key_a=val_a ; 覆盖现有值

[section2]
key_a=val_0
key_b=val_b ; 新增一个值

[section3]  ; 新增一个分区
key_b=val_b
```

**注意，环境变量中 `__` 后缀会被舍弃，只是为了防止键名冲突**

## 许可证

Guo Y.K., MIT License
