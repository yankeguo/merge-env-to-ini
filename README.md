# merge-env-to-ini

A tool to make modifications of ini file from environment variables

[简体中文](README.zh.md)

## Usage

Assuming we have environment variables

```
DEMO_section1__1="key_a=val_a"
DEMO_section2__1="key_b=val_b"
DEMO_section3__1="key_b=val_b"
```

And a file `somefile.ini`

```ini
[section1]
key_a=val_0

[section2]
key_a=val_0
```

By executing command

```shell
merge-env-to-ini --from DEMO_ --to somefile.ini
```

File content will be updated to

```ini
[section1]
key_a=val_a ; value overrided

[section2]
key_a=val_0
key_b=val_b ; value added

[section3]  ; section created
key_b=val_b
```

**Be advised，suffixes `__` in environment variable keys will be ignored, it's for avoiding conflicts**

## Credits

GUO YANKE, MIT License
