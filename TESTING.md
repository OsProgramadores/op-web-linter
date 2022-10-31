# Testing information for developers

## Get list of languages

```
curl -v http://localhost:10000/getlanguages
```

## Make a lint request

```
curl -v --trace-ascii - --json '{ "lang":"xxx", "text":"program_text_here" }' http://localhost:10000/lint
```
