# Resources for ModelHelper

https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet
https://gowebexamples.com/templates/
https://blog.gopheracademy.com/advent-2017/using-go-templates/
https://golang.org/pkg/text/template/
https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff

Testing
https://github.com/onsi/ginkgo
https://github.com/onsi/gomega


pluralisering etc
js: https://github.com/plurals/pluralize
go: https://github.com/gertd/go-pluralize
https://pkg.go.dev/golang.org/x/text/feature/plural

Hent inn pluralize- filen i prosjektet (+ testene)

Må finne ut av denne: https://pkg.go.dev/mod/golang.org/x/text?tab=overview

StripPrefix

CLI UI
https://github.com/manifoldco/promptui
https://github.com/rivo/tview
https://github.com/gdamore/tcell
https://github.com/gizak/termui
https://github.com/wong2/pick
https://manifold.co/blog/improve-your-command-line-go-application-with-promptui-6c4e6fb5a1bc
https://github.com/asaskevich/govalidator
https://appliedgo.net/tui/


Progress
https://github.com/cheggaaa/pb
https://github.com/schollz/progressbar
https://github.com/gosuri/uiprogress


Use Io.Reader and Io.Writer when possible

https://www.youtube.com/watch?v=29LLRKIL_TI

- pointers vs values
It's not a question of performance (generally), but one of shared access
If you vant to share the value with a function or method, then use a pointer
If you dont want to share it, then use a value (copy)

## Pointer receviers

If you want to share a value with its method, use a pointer receiver
Since methods commonly manage state, this is the common usage
Not safe for concurrent access

## Value receivers
If you want the value copied (not shared), use values
If the type is an emptty struct (stateless, just behavior, then just use value)
Safe for concurrent access

## Error
error is an interface

type Error struct {
    Code ErrorCode
    Message string
    Detail interface{}
}

func (e Error) Error() string {
    return fmt.Sprintf("....")
}