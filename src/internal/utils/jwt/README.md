# JWT 工具包使用说明

## 环境变量验证

本 JWT 工具包在导入时会自动验证 `JWT_SECRET` 环境变量是否已设置。

### 验证机制

1. **自动验证**: 当包被导入时，`init()` 函数会自动执行并检查 `JWT_SECRET` 环境变量
2. **快速失败**: 如果环境变量未设置，程序会立即 panic，避免在运行时才发现问题
3. **安全保证**: 确保 JWT 签名密钥始终可用

### 使用方法

#### 1. 设置环境变量

在 `.env` 文件中添加：

```env
JWT_SECRET=your-secret-key-here
```

或者在命令行中设置：

```bash
export JWT_SECRET=your-secret-key-here
```

#### 2. 在代码中使用

```go
package main

import (
    "log"
    "app/src/internal/utils"
)

func main() {
    // 导入 utils 包会触发 init() 函数
    // 如果 JWT_SECRET 未设置，程序会在这里 panic

    // 现在可以安全地使用 JWT 功能
    token, err := utils.IssueJWT("user-123")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Generated token: %s", token)
}
```

#### 3. 基本用法

```go
func loginHandler(c fiber.Ctx) error {
    // 生成 JWT token（使用默认配置）
    token, err := utils.IssueJWT(userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to generate token",
        })
    }

    return c.JSON(fiber.Map{
        "token": token,
    })
}

func protectedHandler(c fiber.Ctx) error {
    // 从请求头获取 token
    token := c.Get("Authorization")

    // 验证 JWT token
    userID, err := utils.VerifyJWT(token)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    // 使用验证后的 userID
    return c.JSON(fiber.Map{
        "user_id": userID,
        "message": "Access granted",
    })
}
```

#### 4. 高级用法 - 自定义选项

```go
func adminLoginHandler(c fiber.Ctx) error {
    // 生成带有自定义选项的 JWT
    token, err := utils.IssueJWT(userID,
        utils.WithDuration(24*time.Hour),        // 24小时有效期
        utils.WithRole("admin"),                 // 设置用户角色
        utils.WithIssuer("my-app"),              // 设置发行者
        utils.WithAudience("my-app-users"),      // 设置受众
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to generate token",
        })
    }

    return c.JSON(fiber.Map{
        "token": token,
        "expires_in": 86400, // 24小时
    })
}

func adminProtectedHandler(c fiber.Ctx) error {
    token := c.Get("Authorization")

    // 解析 token 获取完整信息
    claims, err := utils.ParseJWT(token)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    // 检查用户角色
    if claims.Role != "admin" {
        return c.Status(403).JSON(fiber.Map{
            "error": "Insufficient permissions",
        })
    }

    return c.JSON(fiber.Map{
        "user_id": claims.UserID,
        "role": claims.Role,
        "message": "Admin access granted",
    })
}
```

#### 5. Token 过期管理

```go
func checkTokenStatus(c fiber.Ctx) error {
    token := c.Get("Authorization")

    // 检查 token 是否过期
    expired, err := utils.IsTokenExpired(token)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    if expired {
        return c.Status(401).JSON(fiber.Map{
            "error": "Token expired",
        })
    }

    // 获取过期时间
    expiration, err := utils.GetTokenExpiration(token)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to get expiration time",
        })
    }

    return c.JSON(fiber.Map{
        "valid": true,
        "expires_at": expiration,
        "time_until_expiration": expiration.Sub(time.Now()),
    })
}
```

### 可用函数

#### 核心函数

- `IssueJWT(userID string, options ...JWTOption) (string, error)` - 生成 JWT token
- `VerifyJWT(tokenString string) (string, error)` - 验证 JWT token 并返回用户 ID
- `ParseJWT(tokenString string) (*JWTClaims, error)` - 解析 JWT token 并返回完整声明信息

#### 辅助函数

- `IsTokenExpired(tokenString string) (bool, error)` - 检查 token 是否已过期
- `GetTokenExpiration(tokenString string) (time.Time, error)` - 获取 token 的过期时间

#### 选项函数

- `WithDuration(duration time.Duration) JWTOption` - 设置 token 有效期
- `WithRole(role string) JWTOption` - 设置用户角色
- `WithIssuer(issuer string) JWTOption` - 设置发行者
- `WithAudience(audience string) JWTOption` - 设置受众

### 错误处理

#### 环境变量未设置

```
panic: JWT_SECRET environment variable is not set
```

**解决方案**: 确保在运行程序前设置了 `JWT_SECRET` 环境变量。

#### JWT 生成失败

```go
token, err := utils.IssueJWT(userID)
if err != nil {
    // 处理错误
    log.Printf("Failed to generate JWT: %v", err)
    return
}
```

#### JWT 验证失败

```go
userID, err := utils.VerifyJWT(tokenString)
if err != nil {
    // 处理错误
    log.Printf("Failed to verify JWT: %v", err)
    return
}
```

#### 输入验证错误

```go
// 空 userID
_, err := utils.IssueJWT("")
if err != nil {
    log.Printf("Error: %v", err) // "userID cannot be empty"
}

// 无效持续时间
_, err := utils.IssueJWT("user", utils.WithDuration(-1*time.Hour))
if err != nil {
    log.Printf("Error: %v", err) // "duration must be positive"
}

// 空 token
_, err := utils.VerifyJWT("")
if err != nil {
    log.Printf("Error: %v", err) // "token string cannot be empty"
}
```

### 测试

运行测试：

```bash
# 设置环境变量并运行测试
JWT_SECRET=test-secret-key go test ./src/utils -v
```

### 安全注意事项

1. **密钥强度**: 使用强随机字符串作为 `JWT_SECRET`
2. **环境隔离**: 在不同环境（开发、测试、生产）使用不同的密钥
3. **密钥轮换**: 定期更换 JWT 密钥
4. **HTTPS**: 在生产环境中始终使用 HTTPS 传输 JWT
5. **密钥保护**: JWT 验证函数使用密钥副本，避免直接暴露原始密钥
6. **内存安全**: 密钥在内存中以安全的方式处理，减少泄露风险

### 示例密钥生成

```bash
# 生成 32 字节的随机密钥
openssl rand -base64 32

# 或者使用 Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```
