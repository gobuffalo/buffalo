if exists("b:did_ftplugin")
    finish
endif
let b:did_ftplugin = 1

setlocal comments=s1:#
setlocal commentstring=#\ %s

let b:undo_ftplugin = "setl com< cms<"

" vim:ts=4:sw=4:et
