source_filename = "hello"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"
declare void @solidity.main()
define void @main() {
    tail call void @solidity.main()
    ret void
}