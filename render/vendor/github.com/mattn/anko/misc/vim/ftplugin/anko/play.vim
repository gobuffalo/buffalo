scriptencoding utf-8

function! s:play()
  let code = join(getline(1, '$'), "\n")
  let res = webapi#http#post("http://play-anko.appspot.com/api/play", {"code": code})
  if res.status == "200"
    echo iconv(res.content, "utf-8", &encoding)
  else
    for line in split(res.content, "\n")
      echohl Error | echomsg iconv(line, "utf-8", &encoding) | echohl None
	endfor
  endif
endfunction

command! -buffer PlayAnko call s:play()
