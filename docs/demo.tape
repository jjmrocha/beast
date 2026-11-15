# Records docs/demo.gif. Start the playground app from the manual's worked example
# on 127.0.0.1:8080, then run from the repository root: vhs docs/demo.tape

Output docs/demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 760
Set Height 860
Set Padding 16
Set Framerate 12
Set TypingSpeed 40ms
Set Theme "Catppuccin Mocha"

Hide
Type "make build > /dev/null && export PATH=$PWD/bin:$PATH && cd $(mktemp -d) && clear"
Enter
Sleep 3s
Show

Type "beast template http://127.0.0.1:8080/apps/101 get.yaml"
Sleep 300ms
Enter
Sleep 1500ms

Type "clear && beast run -n 20000 -c 10 get.yaml"
Sleep 300ms
Enter
Sleep 8s
