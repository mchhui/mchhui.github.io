hugo
xcopy public docs /E /I /Y
git add .
git commit -m "auto commit"
set /p firstline=<"../secret_key/GITHUBSK.txt"
git push https://mchhui:%firstline%@github.com/mchhui/mchhui.github.io.git master:master
pause