# 变基操作
将多次零碎的代码（commit）合并提交

1.commit多次零碎的代码后，git log --oneline留意你想合并的哈希值
2.git rebase -i HEAD~3 是你想合并的次数，我这里是3次所以是3
3.提示框可以更改成squash 和 reword ，有什么区别呢？ 答案是reword只改备注信息，而squash是合并代码，并合并备注
4.接着是强制提交 git push --force


