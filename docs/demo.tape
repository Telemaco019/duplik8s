Output demo.gif

Set Width 1400
Set LineHeight 1.8

Set Shell zsh

Type "kubectl duplicate deploy -n my-namespace"
Enter
Sleep 1s
Type "j"
Sleep 300ms
Type "j"
Sleep 1s
Enter
Type "kubectl get deploy -n my-namespace"
Enter
Sleep 1.5s

Type "kubectl duplicate pod -n my-namespace"
Enter
Sleep 1s
Type "j"
Sleep 1s
Enter
Type "kubectl get pods -n my-namespace"
Enter
Sleep 1.5s

Type "kubectl duplicate list -n my-namespace"
Enter
Sleep 3s
