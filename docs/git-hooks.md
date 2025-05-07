**Installing pre-commit**

references:

https://pre-commit.com/

https://goangle.medium.com/golang-improving-your-go-project-with-pre-commit-hooks-a265fad0e02f

we use pre-commit package
```
pip install pre-commit
```

create a file named .pre-commit-config.yaml
```
touch .pre-commit-config.yaml
```

general hooks:
- id: trailing-whitespace
- id: end-of-file-fixer
- id: check-yaml
- id: check-added-large-files

trailing-whitespace to handle any whitespace of the end of line and the new line.

end-of-file-fixer to remove EOF of your whole files project.

check-yaml to fix yaml format file.

check-added-large-files this job will let you to know which file has large file size.

**Running pre-commit**
```
pre-commit run --all-files
```

## for install dependencies:
```shell
pre-commit clean
&&
pre-commit install
&&
pre-commit install --hook-type commit-msg
```

## for applying new added rule:
```shell
pre-commit clean &&
pre-commit install Or pre-commit install --hook-type commit-msg
```

## running pre-commit for test
```shell
pre-commit run --all-files
```

## testing gitlint
```shell
echo "bad commit message" > test_commit_msg.txt
&&
gitlint --msg-filename test_commit_msg.txt
```

## Or:
```shell
pre-commit run gitlint --hook-stage commit-msg --commit-msg-filename .git/COMMIT_EDITMSG
```
