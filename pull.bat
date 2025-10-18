set /p firstline=<"../secret_key/GITHUBSK.txt"
git pull https://mchhui:%firstline%@github.com/mchhui/mchhui.github.io.git master:master
pause