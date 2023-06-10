VERSION 0.6

IMPORT github.com/prelift/earthly-udcs/go:4c1571fb564d582eabf0920f9f9e3e778a3755bc

module:
    DO go+MODULE --BUILDER_TAG=1.20
    DO go+MODULE --BUILDER_TAG=1.19 \
        --TIDY=false # Go 1.19 did not support the -x flag in go tidy