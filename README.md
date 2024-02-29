<div align="center">

## ⭐ White
🌺 GPT-3 on the terminal.

</div>

### 📥 Install
*Only for Linux*
1. Download the release file
2. Move the file to `/usr/local/bin`
3. Create a file with this path `/home/your-user/.white/config.json`
4. Paste there the template down below
5. You're ready to go.
```json
{
    "Key": "OpenAI Key",
    "MaxTokens": <max tokens int>
}
```

<br>

### ❔ Running
To use white you have 2 options
1. With default MaxTokens
```shell 
$ white -q "query"
```
2. Manually specifying MaxTokens
```shell
$ white -q "query" <MaxTokens here>
```

<br>


### 📩 Stats
If you want to check your stats run this
```shell
$ white -s
```

<br>

### 📤 Uninstall
1. Delete the file named `white` in `/usr/local/bin`
2. Delete the folder `/home/your-user/.white/`
