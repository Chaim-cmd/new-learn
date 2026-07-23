# ctx.Abort() 和 ctx.Next() 区别
阻止通行和放行
# 防忘小流程
1. 定义好claims,推荐使用jwt.RegisteredClaims规范内容,
2. 定义好Exp(过期时间)等因素, 使用jwt.NewWithClaims 生成token 
3. 使用token.SignedString(密钥) 发送字符串和error
4. 鉴权中间件,ctx.GetHeader("Authorization") 获取请求头
5. 使用SplintN()对Header进行分割,通常header为"Bearer token" -> ["Bearer","token"]
6. 失败c.Abort(), 成功c.Next
7. 解析验证token 使用ParseWithClaims(),传入jwtSecret
8. 最后判断过期等 !token.Valid
