if exists("b:current_syntax")
  finish
endif

syn case match

syn keyword     ankoDirective         module
syn keyword     ankoDeclaration       var

hi def link     ankoDirective         Statement
hi def link     ankoDeclaration       Type

syn keyword     ankoStatement         return break continue throw
syn keyword     ankoConditional       if else switch try catch finally
syn keyword     ankoLabel             case default
syn keyword     ankoRepeat            for range

hi def link     ankoStatement         Statement
hi def link     ankoConditional       Conditional
hi def link     ankoLabel             Label
hi def link     ankoRepeat            Repeat

syn match       ankoDeclaration       /\<func\>/
syn match       ankoDeclaration       /^func\>/

syn keyword     ankoCast              bytes runes string

hi def link     ankoCast              Type

syn keyword     ankoBuiltins          keys len
syn keyword     ankoBuiltins          println printf print
syn keyword     ankoConstants         true false nil

hi def link     ankoBuiltins          Keyword
hi def link     ankoConstants         Keyword

" Comments; their contents
syn keyword     ankoTodo              contained TODO FIXME XXX BUG
syn cluster     ankoCommentGroup      contains=ankoTodo
syn region      ankoComment           start="#" end="$" contains=@ankoCommentGroup,@Spell

hi def link     ankoComment           Comment
hi def link     ankoTodo              Todo

" anko escapes
syn match       ankoEscapeOctal       display contained "\\[0-7]\{3}"
syn match       ankoEscapeC           display contained +\\[abfnrtv\\'"]+
syn match       ankoEscapeX           display contained "\\x\x\{2}"
syn match       ankoEscapeU           display contained "\\u\x\{4}"
syn match       ankoEscapeBigU        display contained "\\U\x\{8}"
syn match       ankoEscapeError       display contained +\\[^0-7xuUabfnrtv\\'"]+

hi def link     ankoEscapeOctal       ankoSpecialString
hi def link     ankoEscapeC           ankoSpecialString
hi def link     ankoEscapeX           ankoSpecialString
hi def link     ankoEscapeU           ankoSpecialString
hi def link     ankoEscapeBigU        ankoSpecialString
hi def link     ankoSpecialString     Special
hi def link     ankoEscapeError       Error

" Strings and their contents
syn cluster     ankoStringGroup       contains=ankoEscapeOctal,ankoEscapeC,ankoEscapeX,ankoEscapeU,ankoEscapeBigU,ankoEscapeError
syn region      ankoString            start=+"+ skip=+\\\\\|\\"+ end=+"+ contains=@ankoStringGroup
syn region      ankoRawString         start=+`+ end=+`+

hi def link     ankoString            String
hi def link     ankoRawString         String

" Characters; their contents
syn cluster     ankoCharacterGroup    contains=ankoEscapeOctal,ankoEscapeC,ankoEscapeX,ankoEscapeU,ankoEscapeBigU
syn region      ankoCharacter         start=+'+ skip=+\\\\\|\\'+ end=+'+ contains=@ankoCharacterGroup

hi def link     ankoCharacter         Character

" Regions
syn region      ankoBlock             start="{" end="}" transparent fold
syn region      ankoParen             start='(' end=')' transparent

" Integers
syn match       ankoDecimalInt        "\<\d\+\([Ee]\d\+\)\?\>"
syn match       ankoHexadecimalInt    "\<0x\x\+\>"
syn match       ankoOctalInt          "\<0\o\+\>"
syn match       ankoOctalError        "\<0\o*[89]\d*\>"

hi def link     ankoDecimalInt        Integer
hi def link     ankoHexadecimalInt    Integer
hi def link     ankoOctalInt          Integer
hi def link     Integer             Number

" Floating point
syn match       ankoFloat             "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>"
syn match       ankoFloat             "\<\.\d\+\([Ee][-+]\d\+\)\?\>"
syn match       ankoFloat             "\<\d\+[Ee][-+]\d\+\>"

hi def link     ankoFloat             Float
hi def link     ankoImaginary         Number

syn sync minlines=500

let b:current_syntax = "anko"
